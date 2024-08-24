package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/chzyer/readline"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

const ps1 = "> "
const ps2 = ": "

var completer = readline.NewPrefixCompleter(
	readline.PcItem("print("),
)

func main() {
	opts := &syntax.FileOptions{
		Set: true,
	}

	l, err := readline.NewEx(&readline.Config{
		AutoComplete:           completer,
		DisableAutoSaveHistory: true,
		EOFPrompt:              "exit",
		HistoryFile:            "/tmp/readline.tmp",
		InterruptPrompt:        "^C",
		Prompt:                 ps1,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	predeclared := starlark.StringDict{}
	thread := &starlark.Thread{}

	var lines []string
LOOP:
	for {
		line, err := l.Readline()
		switch {
		case err == readline.ErrInterrupt:
			if len(line) == 0 {
				break LOOP
			} else {
				continue
			}
		case err == io.EOF:
			break LOOP
		case line == "exit":
			break LOOP
		}

		lines = append(lines, line)
		switch {
		case strings.HasPrefix(line, " "):
			continue
		case strings.HasSuffix(line, ":"):
			l.SetPrompt(ps2)
			continue
		}

		buffer := strings.Join(lines, "\n")
		lines = lines[:0]
		l.SetPrompt(ps1)
		l.SaveHistory(buffer)

		_, err = opts.ParseExpr("<stdin>", line, 0)
		if err != nil {
			globals, err := starlark.ExecFileOptions(opts, thread, "<stdin>", buffer, predeclared)
			if err != nil {
				fmt.Println(err)
			} else {
				for k, v := range globals {
					predeclared[k] = v
				}
			}
		} else {
			if v, err := starlark.EvalOptions(opts, thread, "<stdin>", buffer, predeclared); err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(v.String())
			}
		}
	}
}
