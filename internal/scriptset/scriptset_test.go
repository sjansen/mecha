package scriptset

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScriptParsing(t *testing.T) {
	require := require.New(t)

	for _, tc := range []struct {
		data     string
		expected map[string]*script
	}{
		{
			data: `script(
			    name="a",
			    commands=cmd("date"),
			)`,
			expected: map[string]*script{
				"a": {
					commands: &cmd{},
				},
			},
		},
	} {
		s := New()
		require.NotNil(s)

		r := strings.NewReader(tc.data)
		err := s.Add("testcase", r)
		require.NoError(err)

		require.Equal(tc.expected, s.scripts)
	}
}
