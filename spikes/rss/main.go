package main

import (
	"fmt"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/v3/process"
)

func main() {
	pid := os.Getpid()
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		die(err)
	}

	for s := "0123456789"; len(s) <= 2000000; s += s {
		info, err := proc.MemoryInfo()
		if err != nil {
			die(err)
		}

		fmt.Printf("RSS: %v (%v)\n",
			humanize.IBytes(info.RSS),
			humanize.IBytes(uint64(len(s))),
		)
	}
}

func die(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
