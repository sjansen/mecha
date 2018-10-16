package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/sjansen/mecha/internal/subprocess"
	"github.com/sjansen/mecha/internal/tui"
)

const maxCount = 10
const maxSleep = 900
const minSleep = 200

var child bool

func init() {
	flag.BoolVar(&child, "as-test-child", false, "")
}

func main() {
	flag.Parse()
	if child {
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
		fmt.Fprintln(os.Stdout, i, "stdout")
		fmt.Fprintln(os.Stderr, i, "stderr")
	}
}

func startChildren() {
	screen := tui.NewScreen()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	addStreamPair := func() {
		// TODO report status
		stdout, stderr, _, err := subprocess.Run(
			ctx,
			os.Args[0],
			"--as-test-child",
		)
		if err != nil {
			die(err)
		}

		screen.AddStreamPair("TODO", stdout, stderr)
	}

	screen.AddMenuItem("Add Row", addStreamPair).
		AddMenuItem("Quit", screen.Stop)

	for i := 0; i < 3; i++ {
		addStreamPair()
	}
	if err := screen.Run(); err != nil {
		die(err)
	}
}
