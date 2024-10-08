package main

import (
	"fmt"
	"math"
	"os"

	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

var files = map[string]string{
	"a.sky": `a = 3`,
	"b.sky": `b = 4`,
}

const script = `
load("a.sky", "a")
load("b.sky", "b")

c = sqrt(float(a*a + b*b))
`

func load(_ *starlark.Thread, module string) (starlark.StringDict, error) {
	opts := &syntax.FileOptions{
		Set: true,
	}
	thread := &starlark.Thread{Load: load}
	globals, err := starlark.ExecFileOptions(opts, thread, module, files[module], nil)
	return globals, err
}

func sqrt(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var x starlark.Float
	if err := starlark.UnpackArgs("sqrt", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	result := math.Sqrt(float64(x))
	return starlark.Float(result), nil
}

func main() {
	opts := &syntax.FileOptions{
		Set: true,
	}

	globals := starlark.StringDict{
		"sqrt": starlark.NewBuiltin("sqrt", sqrt),
	}
	thread := &starlark.Thread{Load: load}
	if result, err := starlark.ExecFileOptions(opts, thread, "<stdin>", script, globals); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		fmt.Println(result["c"])
	}
}
