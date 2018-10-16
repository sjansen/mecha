package main

import (
	"fmt"
	"os"

	humanize "github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	x, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf(
		"Available: %10s\nTotal: %14s\nFree: %9s (%2.0f%%)\nUsed: %9s (%2.0f%%)\n",
		humanize.IBytes(x.Available),
		humanize.IBytes(x.Total),
		humanize.IBytes(x.Free),
		100-x.UsedPercent,
		humanize.IBytes(x.Used),
		x.UsedPercent,
	)
}
