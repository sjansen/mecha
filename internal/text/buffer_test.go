package text

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type step struct {
	write  []byte
	buffer *Buffer
	lines  [][]byte
}

func TestBuffer(t *testing.T) {
	require := require.New(t)

	b := &Buffer{}
	for _, step := range []step{
		step{
			write: []byte("foo\n"),
			buffer: &Buffer{
				lines: []int{4},
				text:  []byte("foo\n"),
			},
			lines: [][]byte{
				[]byte("foo\n"),
			},
		},
		step{
			write: []byte("bar\nbaz"),
			buffer: &Buffer{
				lines: []int{4, 8},
				text:  []byte("foo\nbar\nbaz"),
			},
			lines: [][]byte{
				[]byte("foo\n"),
				[]byte("bar\n"),
				[]byte("baz"),
			},
		},
		step{
			write: []byte("\nqux\n"),
			buffer: &Buffer{
				lines: []int{4, 8, 12, 16},
				text:  []byte("foo\nbar\nbaz\nqux\n"),
			},
			lines: [][]byte{
				[]byte("foo\n"),
				[]byte("bar\n"),
				[]byte("baz\n"),
				[]byte("qux\n"),
			},
		},
	} {
		_, err := b.Write(step.write)
		require.NoError(err)
		require.Equal(step.buffer, b)
		require.Equal(step.lines, b.Lines())
	}
}
