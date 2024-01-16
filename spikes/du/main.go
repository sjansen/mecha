package main

import (
	"fmt"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/v3/disk"
)

func main() {
	partitions, err := disk.Partitions(false)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, p := range partitions {
		fmt.Printf(
			"Mountpoint: %v  Fstype: %v\n",
			p.Mountpoint, p.Fstype,
		)
		if u, err := disk.Usage(p.Mountpoint); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf(
				"Blocks: %2.0f%% (%s/%s)\n",
				u.UsedPercent,
				humanize.IBytes(u.Used),
				humanize.IBytes(u.Total),
			)
			fmt.Printf(
				"Inodes: %2.0f%% (%d/%d)\n",
				u.InodesUsedPercent,
				u.InodesUsed,
				u.InodesTotal,
			)
		}
		fmt.Println("")
	}
}
