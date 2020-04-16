package pytest

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/sjansen/mecha/internal/subprocess"
)

var testcase = regexp.MustCompile(
	`^(?P<file>.+?)::` +
		`(?P<test>.+)[\s]+` +
		`(?P<result>(?:PASS|FAIL)[^\s]+)` +
		`[\s]*(?P<progress>....%.)?$`,
)

func Run(ctx context.Context, args ...string) error {
	args = append([]string{"-v"}, args...)
	p, err := subprocess.New(ctx, "pytest", args...).
		CaptureStderrLines().
		CaptureStdoutLines().
		SetEnv("PYTHONUNBUFFERED", "1").
		Start()
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
		red.Print(line, "\n")
	}
	os.Stdout.Sync()
}

func mapStdoutLines(ch <-chan string) {
	h := newDefaultLineHandler()
	for line := range ch {
		h.onStdoutLine(line)
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

type defaultLineHandler struct {
	green    *color.Color
	start    time.Time
	file     string
	progress string
}

func newDefaultLineHandler() *defaultLineHandler {
	return &defaultLineHandler{
		green: color.New(color.FgGreen),
		start: time.Now(),
	}
}

func (h *defaultLineHandler) onStdoutLine(l string) {
	m := matchLine(l)
	if m != nil {
		h.writeMatchedLine(m)
	} else {
		h.writeUnmatchedLine(l)
	}
	os.Stdout.Sync()
}

func (h *defaultLineHandler) writeMatchedLine(m map[string]string) {
	if m["file"] != h.file {
		if h.file != "" {
			os.Stdout.WriteString("\n")
		}
		if h.progress != "" {
			os.Stdout.WriteString(h.progress)
			os.Stdout.WriteString("  ")
		}
		h.file = m["file"]
		os.Stdout.WriteString(h.file)
		os.Stdout.WriteString("  ")
	}
	if m["result"] == "PASSED" {
		os.Stdout.WriteString(".")
	} else {
		os.Stdout.WriteString("F")
	}
	h.progress = m["progress"]
}

func (h *defaultLineHandler) writeUnmatchedLine(line string) {
	if h.file != "" {
		h.file = ""
		os.Stdout.WriteString("\n")
	}
	elapsed := time.Since(h.start).Truncate(time.Second)
	fmt.Fprintf(os.Stdout, "% 9s  ", elapsed.String())
	h.green.Fprint(os.Stdout, line, "\n")
}
