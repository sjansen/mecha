package subprocess

import (
	"context"
	"os/exec"
	"syscall"

	"github.com/sjansen/mecha/internal/streams"
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
	b1 := &streams.LineBuffer{}
	b2 := &streams.LineBuffer{}
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = b1
	cmd.Stderr = b2

	p = &Subprocess{
		Stdout: b1.Subscribe(),
		Stderr: b2.Subscribe(),
	}

	var es chan *ExitStatus
	if err = cmd.Start(); err == nil {
		es = make(chan *ExitStatus)
		p.PID = cmd.Process.Pid
		p.Status = es
	}

	go func() {
		err := cmd.Wait()
		b1.Close()
		b2.Close()
		if err == nil {
			es <- &ExitStatus{Code: 0}
		} else if _, ok := err.(*exec.ExitError); ok {
			ws, _ := cmd.ProcessState.Sys().(syscall.WaitStatus)
			es <- &ExitStatus{Code: ws.ExitStatus(), Error: err}
		} else {
			es <- &ExitStatus{Code: -1, Error: err}
		}
	}()

	return
}
