package scriptset

import (
	"fmt"
	"io"

	"github.com/google/skylark"
	"github.com/google/skylark/resolve"
)

func init() {
	resolve.AllowFloat = true
	resolve.AllowSet = true
}

type ScriptSet struct {
	globals skylark.StringDict
	thread  *skylark.Thread

	Scripts map[string]*script `json:"scripts"`
}

func New() *ScriptSet {
	s := &ScriptSet{
		thread:  &skylark.Thread{},
		Scripts: make(map[string]*script),
	}
	s.globals = skylark.StringDict{
		"cmd":    skylark.NewBuiltin("cmd", s.cmd),
		"script": skylark.NewBuiltin("script", s.script),
	}
	return s
}

func (set *ScriptSet) Add(filename string, r io.Reader) error {
	if _, err := skylark.ExecFile(set.thread, filename, r, set.globals); err != nil {
		return err
	}
	return nil
}

func (set *ScriptSet) cmd(
	thread *skylark.Thread,
	fn *skylark.Builtin,
	args skylark.Tuple,
	kwargs []skylark.Tuple,
) (skylark.Value, error) {
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

func (set *ScriptSet) script(
	thread *skylark.Thread,
	fn *skylark.Builtin,
	args skylark.Tuple,
	kwargs []skylark.Tuple,
) (skylark.Value, error) {
	var name skylark.String
	var steps skylark.Value
	err := skylark.UnpackArgs(
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
	return skylark.None, nil
}
