package main

import (
	"bufio"
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

	"github.com/fatih/color"
)

const children = 3
const maxCount = 3
const maxSleep = 500
const minSleep = 100

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
	for i := 1; i <= maxCount; i++ {
		sleep := rand.Intn(maxSleep-minSleep) + minSleep
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		fmt.Fprintln(os.Stdout, id, i, "stdout")
		fmt.Fprintln(os.Stderr, id, i, "stderr")
	}
}

func readUntilClosed(c *color.Color, r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		text := c.Sprint(scanner.Text())
		fmt.Println(text)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func startReader(c *color.Color, wg *sync.WaitGroup) (io.WriteCloser, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	wg.Add(1)
	go func(r io.Reader) {
		defer wg.Done()
		readUntilClosed(c, r)
	}(r)

	return w, nil
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

	maxTime := maxCount*maxSleep + 1000
	fmt.Fprintln(os.Stderr, "max time:", maxTime)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(maxTime)*time.Millisecond,
	)
	defer cancel()

	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	wg := &sync.WaitGroup{}
	for i := 1; i <= children; i++ {
		stdout, err := startReader(green, wg)
		if err != nil {
			die(err)
		}

		stderr, err := startReader(red, wg)
		if err != nil {
			die(err)
		}

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
