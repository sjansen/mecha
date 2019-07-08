package main

import "fmt"

func main() {
	base := NewFrobnicator().
		WithA(3).
		WithB(4)
	f1 := base.
		WithAdd().
		Build()
	f2 := base.
		WithMultiply().
		Build()
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

func NewFrobnicator() Builder {
	return func(f *Frobnicator) *Frobnicator {
		return f
	}
}

func (b Builder) Build() *Frobnicator {
	return b(&Frobnicator{})
}

type Builder func(f *Frobnicator) *Frobnicator

func (b Builder) WithA(x int) Builder {
	return func(f *Frobnicator) *Frobnicator {
		b(f).a = x
		return f
	}
}

func (b Builder) WithB(x int) Builder {
	return func(f *Frobnicator) *Frobnicator {
		b(f).b = x
		return f
	}
}

func (b Builder) WithAdd() Builder {
	return func(f *Frobnicator) *Frobnicator {
		b(f).frob = func(a, b int) int {
			return a + b
		}
		return f
	}
}

func (b Builder) WithMultiply() Builder {
	return func(f *Frobnicator) *Frobnicator {
		b(f).frob = func(a, b int) int {
			return a * b
		}
		return f
	}
}
