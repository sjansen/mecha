package subprocess

import (
	"bytes"
	"sync"
)

type lineBuffer struct {
	sync.RWMutex
	tmp         []byte
	subscribers []chan string
}

func (b *lineBuffer) Close() error {
	b.Lock()
	defer b.Unlock()

	if len(b.tmp) > 0 {
		b.write('\n')
	}
	for _, ch := range b.subscribers {
		close(ch)
	}
	b.subscribers = nil

	return nil
}

func (b *lineBuffer) Write(x []byte) (int, error) {
	b.Lock()
	defer b.Unlock()

	return b.write(x...)
}

func (b *lineBuffer) write(x ...byte) (int, error) {
	var lines []string
	n := len(x)

	for {
		idx := bytes.IndexByte(x, byte('\n'))
		if idx == -1 {
			break
		}

		if len(b.tmp) < 1 {
			lines = append(lines, string(x[:idx]))
		} else {
			b.tmp = append(b.tmp, x[:idx]...)
			lines = append(lines, string(b.tmp))
			b.tmp = b.tmp[:0]
		}
		x = x[idx+1:]
	}
	b.tmp = append(b.tmp, x...)

	for _, ch := range b.subscribers {
		for _, line := range lines {
			ch <- line
		}
	}

	return n, nil
}

func (b *lineBuffer) Subscribe() <-chan string {
	b.Lock()
	defer b.Unlock()

	ch := make(chan string)
	b.subscribers = append(b.subscribers, ch)
	return ch
}
