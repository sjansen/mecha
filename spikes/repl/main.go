package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/chzyer/readline"
	"github.com/google/skylark"
	"github.com/google/skylark/resolve"
	"github.com/google/skylark/syntax"
)

const ps1 = "> "
const ps2 = ": "

var completer = readline.NewPrefixCompleter(
	readline.PcItem("print("),
)

func main() {
	resolve.AllowFloat = true
	resolve.AllowSet = true

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

	predeclared := skylark.StringDict{}
	thread := &skylark.Thread{}

	var lines []string
	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		} else if line == "exit" {
			break
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

		_, err = syntax.ParseExpr("<stdin>", line, 0)
		if err != nil {
			if globals, err := skylark.ExecFile(thread, "<stdin>", buffer, predeclared); err != nil {
				fmt.Println(err)
			} else {
				for k, v := range globals {
					predeclared[k] = v
				}
			}
		} else {
			if v, err := skylark.Eval(thread, "<stdin>", buffer, predeclared); err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(v.String())
			}
		}
	}
}
