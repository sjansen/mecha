package main

import (
	"fmt"
	"os/exec"
)

var COMMANDS = []string{
	"go", "mecha", "python2.7", "python3", "source", "virtualenv",
}

func main() {
	var width int
	for _, cmd := range COMMANDS {
		if width < len(cmd) {
			width = len(cmd)
		}
	}
	for _, cmd := range COMMANDS {
		path, err := exec.LookPath(cmd)
		if err != nil {
			path = "not found"
		}
		fmt.Printf("  %-*s  %s\n", width, cmd, path)
	}
}
