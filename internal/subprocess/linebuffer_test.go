package subprocess

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLineBufferCore(t *testing.T) {
	for _, tc := range []struct {
		input []byte
		lines []string
	}{{
		input: []byte("\n"),
		lines: []string{""},
	}, {
		input: []byte("foo\n"),
		lines: []string{"foo"},
	}, {
		input: []byte("foo\nbar"),
		lines: []string{"foo", "bar"},
	}} {
		tc := tc
		t.Run(string(tc.input), func(t *testing.T) {
			b := &lineBuffer{}
			ch := b.Subscribe()
			require.NotNil(t, ch)

			go func() {
				_, err := b.Write(tc.input)
				require.NoError(t, err)
				b.Close()
			}()

			for _, expected := range tc.lines {
				actual := <-ch
				require.Equal(t, expected, actual)
			}
		})
	}
}
