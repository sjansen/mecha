package streams

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type step struct {
	write  []byte
	buffer *LineBuffer
	lines  []string
}

func TestLineBufferCore(t *testing.T) {
	require := require.New(t)

	b := &LineBuffer{}
	for _, step := range []step{
		{
			write: []byte("foo\n"),
			buffer: &LineBuffer{
				tmp: nil,
				lines: []string{
					"foo\n",
				},
			},
			lines: []string{
				"foo\n",
			},
		},
		{
			write: []byte("bar\nbaz"),
			buffer: &LineBuffer{
				tmp: []byte("baz"),
				lines: []string{
					"foo\n", "bar\n",
				},
			},
			lines: []string{
				"foo\n", "bar\n",
			},
		},
		{
			write: []byte("\nqux\n"),
			buffer: &LineBuffer{
				tmp: []byte{},
				lines: []string{
					"foo\n", "bar\n", "baz\n",
					"qux\n",
				},
			},
			lines: []string{
				"foo\n", "bar\n", "baz\n",
				"qux\n",
			},
		},
		{
			write: []byte("qu"),
			buffer: &LineBuffer{
				tmp: []byte("qu"),
				lines: []string{
					"foo\n", "bar\n", "baz\n",
					"qux\n",
				},
			},
			lines: []string{
				"foo\n", "bar\n", "baz\n", "qux\n",
			},
		},
		{
			write: []byte("ux\ncorge"),
			buffer: &LineBuffer{
				tmp: []byte("corge"),
				lines: []string{
					"foo\n", "bar\n", "baz\n",
					"qux\n", "quux\n",
				},
			},
			lines: []string{
				"foo\n", "bar\n", "baz\n",
				"qux\n", "quux\n",
			},
		},
		{
			write: []byte("\ngrault\ngarply\nwaldo\nf"),
			buffer: &LineBuffer{
				tmp: []byte("f"),
				lines: []string{
					"foo\n", "bar\n", "baz\n",
					"qux\n", "quux\n", "corge\n",
					"grault\n", "garply\n", "waldo\n",
				},
			},
			lines: []string{
				"foo\n", "bar\n", "baz\n",
				"qux\n", "quux\n", "corge\n",
				"grault\n", "garply\n", "waldo\n",
			},
		},
		{
			write: []byte("re"),
			buffer: &LineBuffer{
				tmp: []byte("fre"),
				lines: []string{
					"foo\n", "bar\n", "baz\n",
					"qux\n", "quux\n", "corge\n",
					"grault\n", "garply\n", "waldo\n",
				},
			},
			lines: []string{
				"foo\n", "bar\n", "baz\n",
				"qux\n", "quux\n", "corge\n",
				"grault\n", "garply\n", "waldo\n",
			},
		},
		{
			write: []byte("d\n"),
			buffer: &LineBuffer{
				tmp: []byte{},
				lines: []string{
					"foo\n", "bar\n", "baz\n",
					"qux\n", "quux\n", "corge\n",
					"grault\n", "garply\n", "waldo\n",
					"fred\n",
				},
			},
			lines: []string{
				"foo\n", "bar\n", "baz\n",
				"qux\n", "quux\n", "corge\n",
				"grault\n", "garply\n", "waldo\n",
				"fred\n",
			},
		},
	} {
		_, err := b.Write(step.write)
		require.NoError(err)
		require.Equal(step.buffer, b)
		require.Equal(step.lines, b.Lines())
	}
}

func TestLineBufferPubSub(t *testing.T) {
	require := require.New(t)

	b := &LineBuffer{}
	defer b.Close()

	ch := b.Subscribe()
	require.NotNil(ch)

	_, err := b.Write([]byte("foo\nbar\nbaz"))
	require.NoError(err)
	for _, expected := range []string{"foo\n", "bar\n", ""} {
		select {
		case actual := <-ch:
			require.Equal(expected, actual)
		default:
			require.Equal(expected, "")
		}
	}
}
