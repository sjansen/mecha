package scriptset

import (
	"fmt"

	"github.com/google/skylark"
)

type script struct {
	Concurrent bool   `json:"concurrent"`
	Steps      []*cmd `json:"steps"`
}

func (s *script) init(
	steps skylark.Value,
) error {
	switch x := steps.(type) {
	case *cmd:
		s.Concurrent = false
		s.Steps = []*cmd{x}
	case *skylark.List:
		s.Concurrent = false
		if err := s.initSteps(x.Len(), x); err != nil {
			return err
		}
	case *skylark.Set:
		s.Concurrent = true
		if err := s.initSteps(x.Len(), x); err != nil {
			return err
		}
	default:
		return fmt.Errorf(
			"script: got %s, want cmd, list of cmd, or set of cmd", steps.Type(),
		)
	}

	return nil
}

func (s *script) initSteps(n int, x skylark.Iterable) error {
	s.Steps = make([]*cmd, 0, n)

	iter := x.Iterate()
	defer iter.Done()

	var item skylark.Value
	for iter.Next(&item) {
		if cmd, ok := item.(*cmd); ok {
			s.Steps = append(s.Steps, cmd)
		} else {
			return fmt.Errorf(
				"script: got %s, want cmd", x.Type(),
			)
		}
	}

	return nil
}
