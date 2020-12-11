package main

import (
	"fmt"

	"github.com/fatih/color"
	jsonnet "github.com/google/go-jsonnet"
)

var input = `
# Simple Example
local name = 'World';
{
  greeting: "Hello, %s!" % std.asciiUpper(name),
}
`

func main() {
	vm := jsonnet.MakeVM()
	vm.ErrorFormatter.SetColorFormatter(color.New(color.FgRed).Fprintf)
	if output, err := vm.EvaluateAnonymousSnippet("<stdin>", input); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(output)
	}
	input = `std.manifestIni({main:` + input + `, sections:{}})`
	if output, err := vm.EvaluateAnonymousSnippet("<stdin>", input); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(output)
	}
}
