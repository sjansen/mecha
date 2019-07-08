package main

import "fmt"

func main() {
	f1 := NewFrobnicator(
		WithA(3),
		WithB(4),
		WithAdd(),
	)
	f2 := NewFrobnicator(
		WithA(3),
		WithB(4),
		WithMultiply(),
	)
	fmt.Printf("f1(1)=%-3d f1(2)=%-3d\n", f1.Frob(1), f1.Frob(2))
	fmt.Printf("f2(1)=%-3d f2(2)=%-3d\n", f2.Frob(1), f2.Frob(2))
}

type Frobnicator struct {
	a, b int
	frob func(int, int) int
}

func (f *Frobnicator) Frob(c int) int {
	return f.frob(
		f.frob(c, f.a),
		f.frob(c, f.b),
	)
}

func NewFrobnicator(options ...Option) *Frobnicator {
	f := &Frobnicator{}
	for _, option := range options {
		option(f)
	}
	return f
}

type Option func(f *Frobnicator) *Frobnicator

func WithA(x int) Option {
	return func(f *Frobnicator) *Frobnicator {
		f.a = x
		return f
	}
}

func WithB(x int) Option {
	return func(f *Frobnicator) *Frobnicator {
		f.b = x
		return f
	}
}

func WithAdd() Option {
	return func(f *Frobnicator) *Frobnicator {
		f.frob = func(a, b int) int {
			return a + b
		}
		return f
	}
}

func WithMultiply() Option {
	return func(f *Frobnicator) *Frobnicator {
		f.frob = func(a, b int) int {
			return a * b
		}
		return f
	}
}
