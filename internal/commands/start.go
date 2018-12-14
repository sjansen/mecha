package commands

import (
	"fmt"
	"math/rand"
	"os"
	"time"

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

	for _, label := range []string{"Clock:"} {
		updates := make(chan *tui.Status)
		screen.AddStatusItem(label, updates)
		go func() {
			for {
				if ok := rand.Intn(100) > 20; ok {
					updates <- &tui.Status{
						Severity: tui.Healthy,
						Message:  "PASS",
					}
				} else {
					updates <- &tui.Status{
						Severity: tui.Alert,
						Message:  "FAIL",
					}
				}
				time.Sleep(1 * time.Second)
			}
		}()
	}

	screen.AddStatusItem("Disk:", startDiskStatus(root))
	screen.AddStatusItem("RAM:", startMemoryStatus())

	return screen.Run()
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
