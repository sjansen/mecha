package subprocess

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"syscall"
)

type ExitStatus struct {
	Code  int
	Error error
}

type Subprocess struct {
	PID    int
	Status <-chan *ExitStatus
	Stdout <-chan string
	Stderr <-chan string
}

func Run(ctx context.Context, name string, args ...string) (*Subprocess, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	es := make(chan *ExitStatus)
	p := &Subprocess{
		Stdout: makeLineChannel(stdout),
		Stderr: makeLineChannel(stderr),
		Status: es,
	}

	pid, err := waitForCommand(cmd, es)
	p.PID = pid

	return p, err
}

func waitForCommand(cmd *exec.Cmd, es chan<- *ExitStatus) (int, error) {
	if err := cmd.Start(); err != nil {
		return 0, err
	}
	go func() {
		defer close(es)
		err := cmd.Wait()
		if err == nil {
			es <- &ExitStatus{}
		} else {
			es <- &ExitStatus{
				Code:  cmd.ProcessState.ExitCode(),
				Error: err,
			}
		}
	}()
	return cmd.Process.Pid, nil
}

func makeLineChannel(r io.Reader) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		reader := bufio.NewReader(r)
		for {
			line, err := reader.ReadString('\n')
			switch {
			case err == nil:
				ch <- line
			case err == io.EOF:
				ch <- line
				return
			default:
				return
			}
		}
	}()
	return ch
}
