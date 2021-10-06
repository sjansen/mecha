package main

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
)

// https://cuelang.org/

func main() {
	const config = `
	msg:   "Hello \(place)!"
	place: "world"
	`

	ctx := cuecontext.New()
	v := ctx.CompileString(config)
	str := v.LookupPath(cue.ParsePath("msg"))

	fmt.Println(str)
}
