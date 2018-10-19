package scriptset

import (
	"github.com/google/skylark"
)

type cmd struct {
	args []string
}

func (c *cmd) Freeze() {
	return
}

func (c *cmd) Hash() (uint32, error) {
	return 0, errUnhashable
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
