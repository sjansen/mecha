package subprocess

import (
	"context"
	"os/exec"
	"syscall"

	"github.com/sjansen/mecha/internal/streams"
)

type ExitStatus struct {
	Error  error
	Status int
}

func Run(ctx context.Context, name string, args ...string) (
	stdout, stderr <-chan string, status <-chan *ExitStatus, err error,
) {
	b1 := &streams.LineBuffer{}
	b2 := &streams.LineBuffer{}
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = b1
	cmd.Stderr = b2

	var es chan *ExitStatus
	stdout = b1.Subscribe()
	stderr = b2.Subscribe()
	if err = cmd.Start(); err == nil {
		es = make(chan *ExitStatus)
		status = es
	}

	go func() {
		err := cmd.Wait()
		if err == nil {
			es <- &ExitStatus{Status: 0}
		} else if _, ok := err.(*exec.ExitError); ok {
			ws, _ := cmd.ProcessState.Sys().(syscall.WaitStatus)
			es <- &ExitStatus{Status: ws.ExitStatus()}
		} else {
			es <- &ExitStatus{Error: err, Status: -1}
		}
		b1.Close()
		b2.Close()
	}()

	return
}
