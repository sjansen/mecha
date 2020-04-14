package subprocess

import (
	"context"
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

func Run(ctx context.Context, name string, args ...string) (p *Subprocess, err error) {
	b1 := &lineBuffer{}
	b2 := &lineBuffer{}
	es := make(chan *ExitStatus)
	p = &Subprocess{
		Stdout: b1.Subscribe(),
		Stderr: b2.Subscribe(),
		Status: es,
	}

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = b1
	cmd.Stderr = b2
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err = cmd.Start(); err != nil {
		return nil, err
	}

	p.PID = cmd.Process.Pid
	go func() {
		err := cmd.Wait()
		b1.Close()
		b2.Close()
		if err == nil {
			es <- &ExitStatus{}
		} else {
			es <- &ExitStatus{
				Code:  cmd.ProcessState.ExitCode(),
				Error: err,
			}
		}
	}()

	return
}
