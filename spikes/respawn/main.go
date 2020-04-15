package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/sjansen/mecha/internal/subprocess"
)

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

func makeSomeNoise() {
	pid := os.Getpid()
	fmt.Printf("%2d: Started (pid=%d)\n", id, pid)

	rand.Seed(int64(id))
	n := rand.Intn(12) + 3
	time.Sleep(time.Duration(n) * time.Second)

	n = rand.Intn(7)
	if n < 5 {
		fmt.Printf("%2d: Stopped (pid=%d)\n", id, pid)
		os.Exit(0)
	} else {
		fmt.Printf("%2d: Crashed (pid=%d)\n", id, pid)
		os.Exit(1)
	}

}

func spawn(ctx context.Context, i int) int {
	p, err := subprocess.New(ctx, os.Args[0], "--as-test-child", strconv.Itoa(i)).
		CaptureStdoutLines().
		CaptureStderrLines().
		Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return -1
	}

	for {
		select {
		case line := <-p.Stdout:
			os.Stdout.Write([]byte(line))
			os.Stdout.Write([]byte("\n"))
		case line := <-p.Stderr:
			os.Stderr.Write([]byte(line))
			os.Stderr.Write([]byte("\n"))
		case status := <-p.Status:
			return status.Code
		}
	}
}

func startChildren() {
	var wg sync.WaitGroup

	start := 50
	var crashed, stopped, restarted int
	ctx, cancel := context.WithTimeout(context.Background(), 16*time.Second)
	defer cancel()
	for i := 1; i <= start; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for {
				if rc := spawn(ctx, i); rc == 0 {
					stopped++
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

		}(i)
	}

	wg.Wait()
	fmt.Println(
		"started:", start,
		"| restarted:", restarted,
		"| crashed:", crashed,
		"| stopped:", stopped,
	)
}
