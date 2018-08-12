package text

import (
	"bytes"
	"sync"
)

type Buffer struct {
	mutex sync.Mutex
	tmp   []byte
	lines []string
}

func (b *Buffer) Lines() []string {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return b.lines
}

func (b *Buffer) Write(x []byte) (n int, err error) {
	n = len(x)
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for {
		if idx := bytes.IndexByte(x, byte('\n')); idx == -1 {
			break
		} else if len(b.tmp) < 1 {
			b.lines = append(b.lines, string(x[:idx+1]))
			x = x[idx+1:]
		} else {
			b.tmp = append(b.tmp, x[:idx+1]...)
			b.lines = append(b.lines, string(b.tmp))
			b.tmp = b.tmp[:0]
			x = x[idx+1:]
		}
	}
	b.tmp = append(b.tmp, x...)

	return
}
