package scriptset

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScriptParsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		data     string
		expected map[string]*script
	}{{
		data: `
		script(
		    name="a",
		    commands=cmd("date"),
		)

		script(
		    name="b",
		    commands=cmd("make", "-j", 42),
		)

		script(
		    name="c",
		    commands=cmd("sleep", 0.5),
		)`,
		expected: map[string]*script{
			"a": {
				commands: []*cmd{{
					args: []string{"date"},
				}},
			},
			"b": {
				commands: []*cmd{{
					args: []string{"make", "-j", "42"},
				},
				}},
			"c": {
				commands: []*cmd{{
					args: []string{"sleep", "0.5"},
				},
				}},
		},
	}, {
		data: `
		script(
		    name="list",
		    commands=[
		        cmd("touch", "/tmp/foo"),
		        cmd("ls", "-l", "/tmp"),
		    ],
		)

		script(
		    name="set",
		    commands=set([
		        cmd("ddate"),
		        cmd("fortune"),
		    ]),
		)`,
		expected: map[string]*script{
			"list": {
				commands: []*cmd{{
					args: []string{"touch", "/tmp/foo"},
				}, {
					args: []string{"ls", "-l", "/tmp"},
				}},
			},
			"set": {
				commands: []*cmd{{
					args: []string{"ddate"},
				}, {
					args: []string{"fortune"},
				}},
			},
		},
	}}

	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("i=%d", i), func(t *testing.T) {
			require := require.New(t)
			data := strings.Replace(tc.data, "\t", "", -1)
			t.Log(data)

			s := New()
			require.NotNil(s)

			r := strings.NewReader(data)
			err := s.Add("testcase", r)
			require.NoError(err)
			require.Equal(tc.expected, s.scripts)
		})
	}
}
