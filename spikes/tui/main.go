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
	screen := tui.NewDemoScreen()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	addStreamPair := func() {
		// TODO report status
		p, err := subprocess.Run(
			ctx,
			os.Args[0],
			"--as-test-child",
		)
		if err != nil {
			die(err)
		}

		screen.AddStreamPair("TODO", p.Stdout, p.Stderr)
	}

	screen.AddMenuItem("add-row", "Add Row", addStreamPair).
		AddMenuItem("quit", "Quit", screen.Stop).
		AddStatusItem("clock", "Clock:").
		AddStatusItem("disk", "Disk:").
		AddStatusItem("ram", "RAM:")

	for _, id := range []string{"clock", "disk", "ram"} {
		id := id
		go func() {
			for {
				if ok := rand.Intn(100) > 20; ok {
					screen.UpdateStatusItem(id, "PASS", ok)
				} else {
					screen.UpdateStatusItem(id, "FAIL", ok)
				}
				time.Sleep(1 * time.Second)
			}
		}()
	}
	for i := 0; i < 3; i++ {
		addStreamPair()
	}
	if err := screen.Run(); err != nil {
		die(err)
	}
}
