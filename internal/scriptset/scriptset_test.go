package scriptset

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScriptParsing(t *testing.T) {
	t.Parallel()

	testcases, err := filepath.Glob("testdata/*.sky")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, path := range testcases {
		path := path
		basename := filepath.Base(path)
		t.Run(basename, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			f, err := os.Open(path)
			defer func() {
				_ = f.Close()
			}()
			require.NoError(err)

			scriptset := New()
			require.NotNil(scriptset)
			err = scriptset.Add(path, f)
			require.NoError(err)

			ext := filepath.Ext(path)
			path = path[0:len(path)-len(ext)] + ".json"
			tmp, err := os.ReadFile(path)
			require.NoError(err)

			expected := &ScriptSet{}
			err = json.Unmarshal(tmp, expected)
			require.NoError(err)

			scriptset.globals = nil
			scriptset.thread = nil
			if !assert.Equal(expected, scriptset) {
				actual, err := json.MarshalIndent(scriptset, "", "  ")
				require.NoError(err)

				f, err := os.CreateTemp("", "actual.*.json")
				require.NoError(err)
				defer f.Close()

				if _, err := f.Write(actual); err == nil {
					t.Log(
						"Temp JSON file created to facilitate debugging.",
						"\nexpected:", path,
						"\nactual:", f.Name(),
					)
				}
			}
		})
	}
}
