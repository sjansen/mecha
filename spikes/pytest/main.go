package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/fatih/color"
	"github.com/sjansen/mecha/internal/text"
)

const pytest = "testdata/venv/bin/pytest"

// testdata/example/a_test.py::test__a__1 PASSED                     [  2%]
var testcase = regexp.MustCompile(`^(.+?)::[^\s]+ (.[^\s]+)`)

func die(err error) {
	fmt.Fprintln(os.Stderr, "FATAL:", err)
	os.Exit(1)
}

func exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

func newCommand() (cmd *exec.Cmd, stdout, stderr <-chan string) {
	b1 := &text.LineBuffer{}
	b2 := &text.LineBuffer{}

	cmd = exec.Command(pytest, "-v")
	cmd.Stdout = b1
	cmd.Stderr = b2

	stdout = b1.Subscribe()
	stderr = b2.Subscribe()

	return
}

func main() {
	if exists, err := exists(pytest); err != nil {
		die(err)
	} else if !exists {
		fmt.Fprintln(os.Stderr, "No such file:", pytest)
		os.Exit(1)
	}

	cmd, stdout, stderr := newCommand()
	go func() {
		green := color.New(color.FgGreen)
		var testfile string
		for line := range stdout {
			if match := testcase.FindStringSubmatch(line); match == nil {
				if testfile != "" {
					testfile = ""
					os.Stdout.WriteString("\n")
				}
				green.Print("STDOUT: ", line)
			} else {
				if match[1] != testfile {
					if testfile != "" {
						os.Stdout.WriteString("\n")
					}
					testfile = match[1]
					os.Stdout.WriteString(testfile)
					os.Stdout.WriteString("  ")
				}
				if match[2] == "PASSED" {
					os.Stdout.WriteString(".")
				} else {
					os.Stdout.WriteString("F")
				}
				os.Stdout.Sync()
			}
		}
	}()
	go func() {
		red := color.New(color.FgRed)
		for line := range stderr {
			red.Print(line)
		}
	}()

	err := cmd.Run()
	if err != nil {
		die(err)
	}
}
