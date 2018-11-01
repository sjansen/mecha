package scriptset

import (
	"fmt"
	"hash/fnv"

	"github.com/google/skylark"
)

type cmd struct {
	args []string
}

func (c *cmd) init(args skylark.Tuple) error {
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
			return fmt.Errorf(
				"cmd: got %s, want string, int, or float", val.Type(),
			)
		}
	}
	c.args = tmp
	return nil
}

func (c *cmd) Freeze() {
	return
}

func (c *cmd) Hash() (uint32, error) {
	h := fnv.New32()
	for _, arg := range c.args {
		h.Write([]byte(arg))
		h.Write([]byte{0})
	}
	return h.Sum32(), nil
}

func (c *cmd) String() string {
	return "cmd"
}

func (c *cmd) Truth() skylark.Bool {
	if len(c.args) > 0 {
		return skylark.True
	}
	return skylark.False
}

func (c *cmd) Type() string {
	return "cmd"
}
