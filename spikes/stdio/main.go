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

	"github.com/fatih/color"
	"github.com/sjansen/mecha/internal/subprocess"
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

func startReader(wg *sync.WaitGroup, c *color.Color, ch <-chan string) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for line := range ch {
			line := c.Sprint(line)
			fmt.Print(line, "\n")
		}
	}()
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
		fmt.Fprintln(os.Stderr, "starting:", i)
		p, err := subprocess.Run(
			ctx,
			os.Args[0],
			"--as-test-child",
			strconv.Itoa(i),
		)
		if err != nil {
			die(err)
		}

		fmt.Fprintln(os.Stderr, "started:", i)
		startReader(wg, green, p.Stdout)
		startReader(wg, red, p.Stderr)

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			status := <-p.Status
			if status.Error != nil {
				fmt.Fprintf(os.Stderr, "stopped: %d (err=%s)\n", i, status.Error)
			} else {
				fmt.Fprintf(os.Stderr, "stopped: %d (rc=%d)\n", i, status.Code)
			}
		}(i)
	}
	wg.Wait()
}
