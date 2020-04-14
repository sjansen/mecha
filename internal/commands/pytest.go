package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sync"
	"time"

	"github.com/fatih/color"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/sjansen/mecha/internal/subprocess"
)

type pytestCmd struct {
	args    []string
	timeout int
}

func (cmd *pytestCmd) register(app *kingpin.Application) {
	pytest := app.Command("pytest", "Run pytest while capturing output and metrics").
		Action(cmd.run)
	pytest.Flag("timeout", "maximum run time in seconds").
		Short('t').Default("60").
		IntVar(&cmd.timeout)
	pytest.Arg("ARGS", "pytest arguments").Required().
		StringsVar(&cmd.args)
}

// testdata/example/a_test.py::test__a__1 PASSED                     [  2%]
var testcase = regexp.MustCompile(
	`^(?P<file>.+?)::` +
		`(?P<tc>.+)[\s]+` +
		`(?P<result>(?:PASS|FAIL)[^\s]+)[\s]+` +
		`(....%.)?$`,
)

func (cmd *pytestCmd) run(pc *kingpin.ParseContext) error {
	if _, err := exec.LookPath("pytest"); err != nil {
		fmt.Fprintln(os.Stderr, "Command not found:", "pytest")
		return err
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(cmd.timeout)*time.Second,
	)
	defer cancel()
	p, err := subprocess.Run(ctx, "pytest", cmd.args...)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		red := color.New(color.FgRed)
		for line := range p.Stderr {
			red.Print(line, "\n")
		}
		os.Stdout.Sync()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		green := color.New(color.FgGreen)
		var progress, testfile string
		for line := range p.Stdout {
			if match := testcase.FindStringSubmatch(line); match == nil {
				if testfile != "" {
					testfile = ""
					os.Stdout.WriteString("\n")
				}
				green.Print(line, "\n")
			} else {
				if match[1] != testfile {
					if testfile != "" {
						os.Stdout.WriteString("\n")
					}
					if progress != "" {
						os.Stdout.WriteString(progress)
						os.Stdout.WriteString("  ")
					}
					testfile = match[1]
					os.Stdout.WriteString(testfile)
					os.Stdout.WriteString("  ")
				}
				if match[3] == "PASSED" {
					os.Stdout.WriteString(".")
				} else {
					os.Stdout.WriteString("F")
				}
				progress = match[4]
				os.Stdout.Sync()
			}
		}
	}()

	wg.Wait()
	status := <-p.Status
	return status.Error
}
