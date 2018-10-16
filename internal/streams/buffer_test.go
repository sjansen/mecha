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
					"foo",
				},
			},
			lines: []string{
				"foo",
			},
		},
		{
			write: []byte("bar\nbaz"),
			buffer: &LineBuffer{
				tmp: []byte("baz"),
				lines: []string{
					"foo", "bar",
				},
			},
			lines: []string{
				"foo", "bar",
			},
		},
		{
			write: []byte("\nqux\n"),
			buffer: &LineBuffer{
				tmp: []byte{},
				lines: []string{
					"foo", "bar", "baz",
					"qux",
				},
			},
			lines: []string{
				"foo", "bar", "baz",
				"qux",
			},
		},
		{
			write: []byte("qu"),
			buffer: &LineBuffer{
				tmp: []byte("qu"),
				lines: []string{
					"foo", "bar", "baz",
					"qux",
				},
			},
			lines: []string{
				"foo", "bar", "baz", "qux",
			},
		},
		{
			write: []byte("ux\ncorge"),
			buffer: &LineBuffer{
				tmp: []byte("corge"),
				lines: []string{
					"foo", "bar", "baz",
					"qux", "quux",
				},
			},
			lines: []string{
				"foo", "bar", "baz",
				"qux", "quux",
			},
		},
		{
			write: []byte("\ngrault\ngarply\nwaldo\nf"),
			buffer: &LineBuffer{
				tmp: []byte("f"),
				lines: []string{
					"foo", "bar", "baz",
					"qux", "quux", "corge",
					"grault", "garply", "waldo",
				},
			},
			lines: []string{
				"foo", "bar", "baz",
				"qux", "quux", "corge",
				"grault", "garply", "waldo",
			},
		},
		{
			write: []byte("re"),
			buffer: &LineBuffer{
				tmp: []byte("fre"),
				lines: []string{
					"foo", "bar", "baz",
					"qux", "quux", "corge",
					"grault", "garply", "waldo",
				},
			},
			lines: []string{
				"foo", "bar", "baz",
				"qux", "quux", "corge",
				"grault", "garply", "waldo",
			},
		},
		{
			write: []byte("d\n"),
			buffer: &LineBuffer{
				tmp: []byte{},
				lines: []string{
					"foo", "bar", "baz",
					"qux", "quux", "corge",
					"grault", "garply", "waldo",
					"fred",
				},
			},
			lines: []string{
				"foo", "bar", "baz",
				"qux", "quux", "corge",
				"grault", "garply", "waldo",
				"fred",
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
	for _, expected := range []string{"foo", "bar", ""} {
		select {
		case actual := <-ch:
			require.Equal(expected, actual)
		default:
			require.Equal(expected, "")
		}
	}
}
