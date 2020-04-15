package pytest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegex(t *testing.T) {
	for _, tc := range []struct {
		line     string
		expected map[string]string
	}{{`tests/a_test.py::test__a__1 PASSED    [  2%]`,
		map[string]string{
			"file":     `tests/a_test.py`,
			"test":     `test__a__1`,
			"result":   `PASSED`,
			"progress": `[  2%]`,
		},
	}, {`tests/test_b.py::test_foo[bar baz] PASSED`,
		map[string]string{
			"file":     `tests/test_b.py`,
			"test":     `test_foo[bar baz]`,
			"result":   `PASSED`,
			"progress": ``,
		},
	}} {
		tc := tc
		t.Run(tc.line, func(t *testing.T) {
			actual := matchLine(tc.line)
			require.Equal(t, tc.expected, actual)
		})
	}
}
