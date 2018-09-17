package main

import (
	"fmt"
	"os"

	"github.com/chzyer/readline"
)

func die(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func main() {
	stdin := readline.GetStdin()
	if !readline.IsTerminal(stdin) {
		fmt.Println("not a terminal")
		os.Exit(0)
	}

	name, err := readline.Line("What is your name? ")
	if err != nil {
		die(err)
	}
	quest, err := readline.Line("What is your quest? ")
	if err != nil {
		die(err)
	}
	color, err := readline.Password("What is your favorite color? ")
	if err != nil {
		die(err)
	}
	fmt.Printf("name=%q\nquest=%q\ncolor=%q\n", name, quest, color)
}
