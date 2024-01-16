package scriptset

import (
	"fmt"
	"io"

	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

type ScriptSet struct {
	globals starlark.StringDict
	thread  *starlark.Thread

	Scripts map[string]*script `json:"scripts"`
}

func New() *ScriptSet {
	s := &ScriptSet{
		thread:  &starlark.Thread{},
		Scripts: make(map[string]*script),
	}
	s.globals = starlark.StringDict{
		"cmd":    starlark.NewBuiltin("cmd", s.cmd),
		"script": starlark.NewBuiltin("script", s.script),
	}
	return s
}

func (set *ScriptSet) Add(filename string, r io.Reader) error {
	opts := &syntax.FileOptions{
		Set: true,
	}
	if _, err := starlark.ExecFileOptions(opts, set.thread, filename, r, set.globals); err != nil {
		return err
	}
	return nil
}

func (set *ScriptSet) cmd(
	thread *starlark.Thread,
	fn *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	if len(kwargs) > 0 {
		err := fmt.Errorf("%s: unexpected keyword arguments", fn.Name())
		return nil, err
	}

	cmd := &cmd{}
	if err := cmd.init(args); err != nil {
		return nil, err
	}

	return cmd, nil
}

// script(
//
//	name,
//	steps,
//	check=None,
//	recover=None,
//
// )
func (set *ScriptSet) script(
	thread *starlark.Thread,
	fn *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var name starlark.String
	var steps starlark.Value
	err := starlark.UnpackArgs(
		fn.Name(), args, kwargs,
		"name", &name,
		"steps", &steps,
	)
	if err != nil {
		return nil, err
	}

	script := &script{}
	if err := script.init(steps); err != nil {
		return nil, err
	}

	set.Scripts[name.GoString()] = script
	return starlark.None, nil
}
