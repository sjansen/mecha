package scriptset

import (
	"fmt"

	"github.com/google/skylark"
)

type script struct {
	Concurrent bool   `json:"concurrent"`
	Commands   []*cmd `json:"commands"`
}

func (s *script) init(
	commands skylark.Value,
) error {
	switch x := commands.(type) {
	case *cmd:
		s.Concurrent = false
		s.Commands = []*cmd{x}
	case *skylark.List:
		s.Concurrent = false
		if err := s.initCommands(x.Len(), x); err != nil {
			return err
		}
	case *skylark.Set:
		s.Concurrent = true
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
	s.Commands = make([]*cmd, 0, n)

	iter := x.Iterate()
	defer iter.Done()

	var item skylark.Value
	for iter.Next(&item) {
		if cmd, ok := item.(*cmd); ok {
			s.Commands = append(s.Commands, cmd)
		} else {
			return fmt.Errorf(
				"script: got %s, want cmd", x.Type(),
			)
		}
	}

	return nil
}
