package text

import "bytes"

type Buffer struct {
	lines []int
	text  []byte
}

func (b *Buffer) Lines() [][]byte {
	lines := make([][]byte, 0, len(b.lines)+1)
	begin := 0
	for _, end := range b.lines {
		lines = append(lines, b.text[begin:end])
		begin = end
	}
	if begin < len(b.text) {
		lines = append(lines, b.text[begin:])
	}
	return lines
}

func (b *Buffer) Write(x []byte) (n int, err error) {
	n = len(x)
	for {
		if idx := bytes.IndexByte(x, byte('\n')); idx == -1 {
			break
		} else {
			b.text = append(b.text, x[:idx+1]...)
			b.lines = append(b.lines, len(b.text))
			x = x[idx+1:]
		}
	}
	b.text = append(b.text, x...)
	return
}
