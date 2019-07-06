package main

import (
	"fmt"
	"os"

	"cuelang.org/go/cue"
)

func main() {
	const config = `
	msg:   "Hello \(place)!"
	place: "world"
	`

	var r cue.Runtime

	instance, err := r.Parse("test", config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	str, err := instance.Lookup("msg").String()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Println(str)
}
