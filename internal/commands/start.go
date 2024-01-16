package commands

import (
	"context"
	"fmt"
	"os"
	"sort"
	"syscall"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/beevik/ntp"
	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/sjansen/mecha/internal/config"
	"github.com/sjansen/mecha/internal/fs"
	"github.com/sjansen/mecha/internal/subprocess"
	"github.com/sjansen/mecha/internal/tui"
)

type startCmd struct {
	procfile string
}

func (cmd *startCmd) register(app *kingpin.Application) {
	x := app.Command("start", "Start the application defined by Procfile").
		Action(cmd.run)
	x.Flag("procfile", `proc file (default "Procfile")`).
		Short('f').
		Default("Procfile").
		ExistingFileVar(&cmd.procfile)
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
	screen.AddStatusItem("Clock:", startClockStatus())
	screen.AddStatusItem("Disk:", startDiskStatus(root))
	screen.AddStatusItem("RAM:", startMemoryStatus())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = startProcs(ctx, screen, cmd.procfile)
	if err != nil {
		return err
	}

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
			// https://www.ntppool.org/vendors.html
			server := "0.beevik-ntp.pool.ntp.org"
			options := ntp.QueryOptions{Timeout: 30 * time.Second}
			if x, err := ntp.QueryWithOptions(server, options); err != nil {
				update.Severity = tui.Unknown
				update.Message = err.Error()
			} else {
				offset := x.ClockOffset.Round(time.Second)
				switch {
				case offset < time.Minute:
					update.Severity = tui.Healthy
					update.Message = fmt.Sprintf("PASS (%s)", offset)
				case offset < 3*time.Minute:
					update.Severity = tui.Warning
					update.Message = fmt.Sprintf("WARNING (%s)", offset)
				default:
					update.Severity = tui.Alert
					update.Message = fmt.Sprintf("FAIL (%s)", offset)
				}
			}
			updates <- update
			time.Sleep(time.Hour)
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
				var status string
				switch {
				case x.UsedPercent < 85:
					update.Severity = tui.Healthy
					status = "PASS"
				case x.UsedPercent < 95:
					update.Severity = tui.Warning
					status = "WARNING"
				default:
					update.Severity = tui.Alert
					status = "FAIL"
				}
				update.Message = fmt.Sprintf("%s (%2.0f%% - %s/%s)",
					status,
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
				var status string
				switch {
				case x.UsedPercent < 85:
					update.Severity = tui.Healthy
					status = "PASS"
				case x.UsedPercent < 95:
					update.Severity = tui.Warning
					status = "WARNING"
				default:
					update.Severity = tui.Alert
					status = "FAIL"
				}
				update.Message = fmt.Sprintf("%s (%2.0f%% - %s/%s)",
					status,
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

func startProcs(ctx context.Context, screen *tui.StackedTextViews, filename string) error {
	procfile, err := os.Open(filename)
	if err != nil {
		return err
	}
	procs, err := config.ReadProcfile(procfile)
	if err != nil {
		return err
	}

	keys := make([]string, 0, len(procs))
	for k := range procs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		p, err := subprocess.New(ctx, "sh", "-c", procs[k]).
			CaptureStdoutLines().
			CaptureStderrLines().
			Start()
		if err != nil {
			return err
		}

		updates := make(chan *tui.Status)
		screen.AddStdView(" "+k+" ", p.Stdout, p.Stderr, updates)
		go func() {
			msg := fmt.Sprintf("PID: %d", p.PID)
			updates <- &tui.Status{Message: msg}

			select {
			case <-ctx.Done():
				// noop
			case status := <-p.Status:
				if status.Code == -1 {
					msg = fmt.Sprintf("exited; reason=%q", status.Error)
				} else {
					msg = fmt.Sprintf("exit=%d", status.Code)
				}
				updates <- &tui.Status{Message: msg}
			}
			syscall.Kill(-p.PID, syscall.SIGHUP)
		}()
	}

	return nil
}
