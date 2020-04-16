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
		`(?P<result>(?:PASS|ERROR|FAIL|SKIP|XFAIL)[^\s]*)` +
		`[\s]*(?P<progress>\[.*?\])?$`,
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
		stdout := p.Stdout
		stderr := p.Stderr
		h := newDefaultLineHandler()
		for {
			select {
			case line, ok := <-stdout:
				h.onStdoutLine(line)
				if !ok {
					stdout = nil
				}
			case line, ok := <-stderr:
				h.onStderrLine(line)
				if !ok {
					stderr = nil
				}
			}
			if stdout == nil && stderr == nil {
				break
			}
		}
	}()

	wg.Wait()
	status := <-p.Status
	return status.Error
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
	red      *color.Color
	start    time.Time
	file     string
	progress string
}

func newDefaultLineHandler() *defaultLineHandler {
	return &defaultLineHandler{
		green: color.New(color.FgGreen),
		red:   color.New(color.FgRed),
		start: time.Now(),
	}
}

func (h *defaultLineHandler) onStderrLine(l string) {
	h.red.Print(l, "\n")
	os.Stdout.Sync()
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
	switch m["result"] {
	case "PASSED":
		os.Stdout.WriteString(".")
	case "ERROR":
		os.Stdout.WriteString("E")
	case "FAILED":
		os.Stdout.WriteString("F")
	case "SKIPPED":
		os.Stdout.WriteString("s")
	case "XFAIL":
		os.Stdout.WriteString("x")
	default:
		os.Stdout.WriteString("?")
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
