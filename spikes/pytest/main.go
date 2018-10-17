package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/fatih/color"
	"github.com/sjansen/mecha/internal/subprocess"
)

const pytest = "testdata/venv/bin/pytest"

// testdata/example/a_test.py::test__a__1 PASSED                     [  2%]
var testcase = regexp.MustCompile(`^(.+?)::[^\s]+ (.[^\s]+).*(....%.)$`)

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

func main() {
	if exists, err := exists(pytest); err != nil {
		die(err)
	} else if !exists {
		fmt.Fprintln(os.Stderr, "No such file:", pytest)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		60*time.Second,
	)
	p, err := subprocess.Run(ctx, pytest, "-v")
	if err != nil {
		die(err)
	}
	defer cancel()

	go func() {
		green := color.New(color.FgGreen)
		var progress, testfile string
		for line := range p.Stdout {
			if match := testcase.FindStringSubmatch(line); match == nil {
				if testfile != "" {
					testfile = ""
					os.Stdout.WriteString("\n")
				}
				green.Print("STDOUT: ", line, "\n")
			} else {
				if match[1] != testfile {
					if testfile == "" {
						os.Stdout.WriteString("[  0%]  ")
					} else {
						os.Stdout.WriteString("\n")
						os.Stdout.WriteString(progress)
						os.Stdout.WriteString("  ")
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
				progress = match[3]
				os.Stdout.Sync()
			}
		}
	}()
	go func() {
		red := color.New(color.FgRed)
		for line := range p.Stderr {
			red.Print(line, "\n")
		}
	}()

	status := <-p.Status
	if status.Error != nil {
		die(status.Error)
	}
}
