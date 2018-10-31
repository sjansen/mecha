package scriptset

import (
	"hash/fnv"

	"github.com/google/skylark"
)

type cmd struct {
	args []string
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
	return skylark.True
}

func (c *cmd) Type() string {
	return "cmd"
}
