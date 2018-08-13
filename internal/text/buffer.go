package text

import (
	"bytes"
	"sync"
)

type Buffer struct {
	mutex       sync.RWMutex
	tmp         []byte
	lines       []string
	subscribers []chan string
}

func (b *Buffer) Close() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	for _, ch := range b.subscribers {
		close(ch)
	}
	b.subscribers = nil
	return nil
}

func (b *Buffer) Lines() []string {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.lines
}

func (b *Buffer) Write(x []byte) (n int, err error) {
	var lines []string
	n = len(x)

	b.mutex.Lock()
	defer func() {
		subscribers := b.subscribers
		b.lines = append(b.lines, lines...)
		b.mutex.Unlock()
		for _, ch := range subscribers {
			for _, line := range lines {
				select {
				case ch <- line:
				default:
				}
			}
		}
	}()

	for {
		idx := bytes.IndexByte(x, byte('\n'))
		if idx == -1 {
			break
		}

		var line string
		if len(b.tmp) < 1 {
			line = string(x[:idx+1])
		} else {
			b.tmp = append(b.tmp, x[:idx+1]...)
			line = string(b.tmp)
			b.tmp = b.tmp[:0]
		}
		lines = append(lines, line)
		x = x[idx+1:]
	}
	b.tmp = append(b.tmp, x...)

	return
}

func (b *Buffer) Subscribe() <-chan string {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	ch := make(chan string, 10)
	b.subscribers = append(b.subscribers, ch)
	return ch
}
