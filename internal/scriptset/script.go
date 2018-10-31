package scriptset

import (
	"fmt"

	"github.com/google/skylark"
)

type script struct {
	commands []*cmd
}

func (s *script) addCommands(fnName string, x skylark.Iterable) error {
	iter := x.Iterate()
	defer iter.Done()
	var item skylark.Value
	for iter.Next(&item) {
		if cmd, ok := item.(*cmd); ok {
			s.commands = append(s.commands, cmd)
		} else {
			return fmt.Errorf(
				"%s: got %s, want cmd", fnName, x.Type(),
			)
		}
	}
	return nil
}
