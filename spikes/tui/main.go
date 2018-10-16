package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/sjansen/mecha/internal/tui"
)

func main() {
	screen := tui.NewScreen()
	addStreamPair := func() {
		stdout := make(chan string)
		stderr := make(chan string)
		screen.AddStreamPair("TODO", stdout, stderr)
		go func() {
			for i := 1; i <= 15; i++ {
				if i%10 == 0 {
					stderr <- fmt.Sprintf("line #%d", i)
				} else {
					stdout <- fmt.Sprintf("line #%d", i)
				}
				n := rand.Int()%750 + 250
				time.Sleep(time.Duration(n) * time.Millisecond)
			}
			close(stdout)
			close(stderr)
		}()
	}
	screen.AddMenuItem("Add Row", addStreamPair).
		AddMenuItem("Quit", screen.Stop)
	for i := 0; i < 3; i++ {
		addStreamPair()
	}
	if err := screen.Run(); err != nil {
		panic(err)
	}
}
