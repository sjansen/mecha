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

	scripts map[string]*script
}

func New() *ScriptSet {
	s := &ScriptSet{
		thread:  &skylark.Thread{},
		scripts: make(map[string]*script),
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
	tmp := make([]string, 0, len(args))
	for _, val := range args {
		switch x := val.(type) {
		case skylark.Float:
			tmp = append(tmp, x.String())
		case skylark.Int:
			tmp = append(tmp, x.String())
		case skylark.String:
			tmp = append(tmp, x.GoString())
		default:
			err := fmt.Errorf(
				"%s: got %s, want string, int, or float", fn.Name(), val.Type(),
			)
			return nil, err
		}
	}
	return &cmd{args: tmp}, nil
}

func (set *ScriptSet) script(
	thread *skylark.Thread,
	fn *skylark.Builtin,
	args skylark.Tuple,
	kwargs []skylark.Tuple,
) (skylark.Value, error) {
	var name skylark.String
	var val skylark.Value
	if err := skylark.UnpackArgs(fn.Name(), args, kwargs, "name", &name, "commands", &val); err != nil {
		return nil, err
	}
	k := name.GoString()
	v := &script{}
	switch x := val.(type) {
	case *cmd:
		v.commands = []*cmd{x}
	case *skylark.List:
		v.commands = make([]*cmd, 0, x.Len())
		if err := v.addCommands(fn.Name(), x); err != nil {
			return nil, err
		}
	case *skylark.Set:
		v.commands = make([]*cmd, 0, x.Len())
		if err := v.addCommands(fn.Name(), x); err != nil {
			return nil, err
		}
	default:
		err := fmt.Errorf(
			"%s: got %s, want cmd, list of cmd, or set of cmd", fn.Name(), val.Type(),
		)
		return nil, err
	}
	set.scripts[k] = v
	return skylark.None, nil
}