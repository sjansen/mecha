package main

import (
	"fmt"
	"math"
	"os"

	"github.com/google/skylark"
	"github.com/google/skylark/resolve"
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

func load(_ *skylark.Thread, module string) (skylark.StringDict, error) {
	thread := &skylark.Thread{Load: load}
	globals, err := skylark.ExecFile(thread, module, files[module], nil)
	return globals, err
}

func sqrt(
	thread *skylark.Thread,
	_ *skylark.Builtin,
	args skylark.Tuple,
	kwargs []skylark.Tuple,
) (skylark.Value, error) {
	var x skylark.Float
	if err := skylark.UnpackArgs("sqrt", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	result := math.Sqrt(float64(x))
	return skylark.Float(result), nil
}

func main() {
	resolve.AllowFloat = true

	globals := skylark.StringDict{
		"sqrt": skylark.NewBuiltin("sqrt", sqrt),
	}
	thread := &skylark.Thread{Load: load}
	if result, err := skylark.ExecFile(thread, "<stdin>", script, globals); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		fmt.Println(result["c"])
	}
}
