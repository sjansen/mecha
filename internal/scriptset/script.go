package scriptset

import (
	"fmt"

	"github.com/google/skylark"
)

type script struct {
	commands []*cmd
}

func (s *script) init(
	commands skylark.Value,
) error {
	switch x := commands.(type) {
	case *cmd:
		s.commands = []*cmd{x}
	case *skylark.List:
		if err := s.initCommands(x.Len(), x); err != nil {
			return err
		}
	case *skylark.Set:
		if err := s.initCommands(x.Len(), x); err != nil {
			return err
		}
	default:
		return fmt.Errorf(
			"script: got %s, want cmd, list of cmd, or set of cmd", commands.Type(),
		)
	}

	return nil
}

func (s *script) initCommands(n int, x skylark.Iterable) error {
	s.commands = make([]*cmd, 0, n)

	iter := x.Iterate()
	defer iter.Done()

	var item skylark.Value
	for iter.Next(&item) {
		if cmd, ok := item.(*cmd); ok {
			s.commands = append(s.commands, cmd)
		} else {
			return fmt.Errorf(
				"script: got %s, want cmd", x.Type(),
			)
		}
	}

	return nil
}
