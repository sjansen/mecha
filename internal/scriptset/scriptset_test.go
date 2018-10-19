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
	}{{
		data: `
		script(
		    name="a",
		    commands=cmd("date"),
		)
		`,
		expected: map[string]*script{
			"a": {
				commands: &cmd{
					args: []string{"date"},
				},
			},
		},
	}, {
		data: `
		script(
		    name="b",
		    commands=cmd("make", "-j", 42),
		)

		script(
		    name="c",
		    commands=cmd("sleep", 0.5),
		)
		`,
		expected: map[string]*script{
			"b": {
				commands: &cmd{
					args: []string{"make", "-j", "42"},
				},
			},
			"c": {
				commands: &cmd{
					args: []string{"sleep", "0.5"},
				},
			},
		},
	}} {
		s := New()
		require.NotNil(s)

		data := strings.Replace(tc.data, "\t", "", -1)
		r := strings.NewReader(data)
		err := s.Add("testcase", r)
		require.NoError(err)

		require.Equal(tc.expected, s.scripts)
	}
}
