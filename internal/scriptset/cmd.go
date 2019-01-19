package scriptset

import (
	"fmt"
	"hash/fnv"

	"go.starlark.net/starlark"
)

type cmd struct {
	Args []string `json:"args"`
}

func (c *cmd) init(args starlark.Tuple) error {
	tmp := make([]string, 0, len(args))
	for _, val := range args {
		switch x := val.(type) {
		case starlark.Float:
			tmp = append(tmp, x.String())
		case starlark.Int:
			tmp = append(tmp, x.String())
		case starlark.String:
			tmp = append(tmp, x.GoString())
		default:
			return fmt.Errorf(
				"cmd: got %s, want string, int, or float", val.Type(),
			)
		}
	}
	c.Args = tmp
	return nil
}

func (c *cmd) Freeze() {
}

func (c *cmd) Hash() (uint32, error) {
	h := fnv.New32()
	for _, arg := range c.Args {
		h.Write([]byte(arg))
		h.Write([]byte{0})
	}
	return h.Sum32(), nil
}

func (c *cmd) String() string {
	return "cmd"
}

func (c *cmd) Truth() starlark.Bool {
	if len(c.Args) > 0 {
		return starlark.True
	}
	return starlark.False
}

func (c *cmd) Type() string {
	return "cmd"
}
