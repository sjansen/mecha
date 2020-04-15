package pytest

import (
	"context"
	"os"
	"regexp"
	"sync"

	"github.com/fatih/color"
	"github.com/sjansen/mecha/internal/subprocess"
)

var testcase = regexp.MustCompile(
	`^(?P<file>.+?)::` +
		`(?P<test>.+)[\s]+` +
		`(?P<result>(?:PASS|FAIL)[^\s]+)` +
		`[\s]*(?P<progress>....%.)?\n?$`,
)

func Run(ctx context.Context, args ...string) error {
	// PYTHONUNBUFFERED=1
	args = append([]string{"-v"}, args...)
	p, err := subprocess.Run(ctx, "pytest", args...)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		mapStderrLines(p.Stderr)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		mapStdoutLines(p.Stdout)
	}()

	wg.Wait()
	status := <-p.Status
	return status.Error
}

func mapStderrLines(ch <-chan string) {
	red := color.New(color.FgRed)
	for line := range ch {
		red.Print(line)
	}
	os.Stdout.Sync()
}

func mapStdoutLines(ch <-chan string) {
	green := color.New(color.FgGreen)

	var file, progress string
	for line := range ch {
		m := matchLine(line)
		if m == nil {
			if file != "" {
				file = ""
				os.Stdout.WriteString("\n")
			}
			green.Print(line)
		} else {
			if m["file"] != file {
				if file != "" {
					os.Stdout.WriteString("\n")
				}
				if progress != "" {
					os.Stdout.WriteString(progress)
					os.Stdout.WriteString("  ")
				}
				file = m["file"]
				os.Stdout.WriteString(file)
				os.Stdout.WriteString("  ")
			}
			if m["result"] == "PASSED" {
				os.Stdout.WriteString(".")
			} else {
				os.Stdout.WriteString("F")
			}
			progress = m["progress"]
			os.Stdout.Sync()
		}
	}
}

func matchLine(line string) map[string]string {
	match := testcase.FindStringSubmatch(line)
	if match == nil {
		return nil
	}
	names := testcase.SubexpNames()
	result := make(map[string]string, len(names))
	for i, name := range names {
		if name != "" {
			result[name] = match[i]
		}
	}
	return result
}
