package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

const children = 3
const max_count = 3
const max_sleep = 500
const min_sleep = 100

var id int

func init() {
	flag.IntVar(&id, "as-test-child", 0, "")
}

func main() {
	flag.Parse()
	if id > 0 {
		makeSomeNoise()
	} else {
		startChildren()
	}
}

func die(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func makeSomeNoise() {
	for i := 1; i <= max_count; i += 1 {
		sleep := rand.Intn(max_sleep-min_sleep) + min_sleep
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		fmt.Fprintln(os.Stdout, id, i)
		fmt.Fprintln(os.Stderr, id, i)
	}
}

func readUntilClosed(r io.Reader) {
	io.Copy(os.Stdout, r)
}

func startChild(ctx context.Context, i int, stdout, stderr io.WriteCloser) {
	cmd := exec.CommandContext(
		ctx,
		os.Args[0],
		"--as-test-child",
		strconv.Itoa(i),
	)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Start(); err != nil {
		die(err)
	}
	stdout.Close()
	stderr.Close()

	fmt.Fprintln(os.Stderr, "started:", i, cmd.Process.Pid)
	if err := cmd.Wait(); err != nil {
		die(err)
	}
}

func startChildren() {
	fmt.Fprintln(os.Stderr, "parent:", os.Getpid())

	max_time := max_count*max_sleep + 1000
	fmt.Fprintln(os.Stderr, "max time:", max_time)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(max_time)*time.Millisecond,
	)
	defer cancel()

	wg := &sync.WaitGroup{}
	for i := 1; i <= children; i += 1 {
		r, stdout, err := os.Pipe()
		if err != nil {
			die(err)
		}

		wg.Add(1)
		go func(i int, r io.Reader) {
			defer wg.Done()
			readUntilClosed(r)
			fmt.Fprintln(os.Stderr, "stdout:", i)
		}(i, r)

		r, stderr, err := os.Pipe()
		if err != nil {
			die(err)
		}

		wg.Add(1)
		go func(i int, r io.Reader) {
			defer wg.Done()
			readUntilClosed(r)
			fmt.Fprintln(os.Stderr, "stderr:", i)
		}(i, r)

		wg.Add(1)
		go func(i int, stdout, stderr io.WriteCloser) {
			defer wg.Done()
			fmt.Fprintln(os.Stderr, "starting:", i)
			startChild(ctx, i, stdout, stderr)
			fmt.Fprintln(os.Stderr, "stopped:", i)
		}(i, stdout, stderr)
	}
	wg.Wait()
}
