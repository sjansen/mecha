package text

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type step struct {
	write  []byte
	buffer *Buffer
	lines  []string
}

func TestBufferCore(t *testing.T) {
	require := require.New(t)

	b := &Buffer{}
	for _, step := range []step{
		{
			write: []byte("foo\n"),
			buffer: &Buffer{
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
			buffer: &Buffer{
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
			buffer: &Buffer{
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
			buffer: &Buffer{
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
			buffer: &Buffer{
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
			buffer: &Buffer{
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
			buffer: &Buffer{
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
			buffer: &Buffer{
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

func TestBufferPubSub(t *testing.T) {
	require := require.New(t)

	b := &Buffer{}
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
