package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

func spawn() int {
	cmd := exec.Command("testdata/script")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return -1
	}

	err = cmd.Wait()
	if err == nil {
		return 0
	} else if _, ok := err.(*exec.ExitError); ok {
		status, _ := cmd.ProcessState.Sys().(syscall.WaitStatus)
		return status.ExitStatus()
	} else {
		fmt.Fprintln(os.Stderr, err)
		return -1
	}
}

func main() {
	var wg sync.WaitGroup

	start := 250
	var crashed, exited, restarted int
	ctx, cancel := context.WithTimeout(context.Background(), 16*time.Second)
	defer cancel()
	for i := 0; i < start; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				if rc := spawn(); rc == 0 {
					exited++
				} else {
					crashed++
				}
				select {
				case <-ctx.Done():
					return
				default:
					restarted++
				}
			}

		}()
	}

	wg.Wait()
	fmt.Println(
		"started:", start,
		"crashed:", crashed,
		"exited:", exited,
		"restarted:", restarted,
	)
}
