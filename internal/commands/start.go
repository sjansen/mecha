package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
	humanize "github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/sjansen/mecha/internal/fs"
	"github.com/sjansen/mecha/internal/tui"
)

type startCmd struct{}

func (cmd *startCmd) register(app *kingpin.Application) {
	app.Command("start", "Start the application defined by Procfile").
		Action(cmd.run)
}

func (cmd *startCmd) run(pc *kingpin.ParseContext) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	root, err := fs.FindProjectRoot(wd)
	if err != nil {
		return err
	}

	screen := tui.NewStackedTextViews()
	screen.AddStdView("TODO", nil, nil).
		AddStdView("TODO", nil, nil)

	screen.AddStatusItem("Clock:", startClockStatus())
	screen.AddStatusItem("Disk:", startDiskStatus(root))
	screen.AddStatusItem("RAM:", startMemoryStatus())

	return screen.Run()
}

func startClockStatus() chan *tui.Status {
	updates := make(chan *tui.Status)
	go func() {
		for {
			updates <- &tui.Status{
				Severity: tui.Refresh,
				Message:  "Checking...",
			}

			update := &tui.Status{}
			options := ntp.QueryOptions{Timeout: 30 * time.Second}
			if x, err := ntp.QueryWithOptions("0.beevik-ntp.pool.ntp.org", options); err != nil {
				update.Severity = tui.Unknown
				update.Message = "???"
			} else {
				offset := x.ClockOffset.Round(time.Second)
				if offset < time.Minute {
					update.Severity = tui.Healthy
					update.Message = fmt.Sprintf("PASS (%s)", offset)
				} else if offset < 3*time.Minute {
					update.Severity = tui.Warning
					update.Message = fmt.Sprintf("WARNING (%s)", offset)
				} else {
					update.Severity = tui.Alert
					update.Message = fmt.Sprintf("FAIL (%s)", offset)
				}
			}
			updates <- update
			time.Sleep(15 * time.Minute)
		}
	}()
	return updates
}

func startDiskStatus(root string) chan *tui.Status {
	updates := make(chan *tui.Status)
	go func() {
		for {
			update := &tui.Status{}
			if x, err := disk.Usage(root); err != nil {
				update.Severity = tui.Unknown
				update.Message = "???"
			} else {
				if x.UsedPercent < 85 {
					update.Severity = tui.Healthy
				} else if x.UsedPercent < 95 {
					update.Severity = tui.Warning
				} else {
					update.Severity = tui.Alert
				}
				update.Message = fmt.Sprintf("%2.0f%% (%s/%s)",
					x.UsedPercent,
					humanize.IBytes(x.Used),
					humanize.IBytes(x.Total),
				)
			}
			updates <- update
			time.Sleep(30 * time.Second)
		}
	}()
	return updates
}

func startMemoryStatus() chan *tui.Status {
	updates := make(chan *tui.Status)
	go func() {
		for {
			update := &tui.Status{}
			if x, err := mem.VirtualMemory(); err != nil {
				update.Severity = tui.Unknown
				update.Message = "???"
			} else {
				if x.UsedPercent < 85 {
					update.Severity = tui.Healthy
				} else if x.UsedPercent < 95 {
					update.Severity = tui.Warning
				} else {
					update.Severity = tui.Alert
				}
				update.Message = fmt.Sprintf("%2.0f%% (%s/%s)",
					x.UsedPercent,
					humanize.IBytes(x.Used),
					humanize.IBytes(x.Total),
				)
			}
			updates <- update
			time.Sleep(5 * time.Second)
		}
	}()
	return updates
}
