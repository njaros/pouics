package a

type A struct {
	i int
}

func NewA(i int) *A {
	return &A{i}
}

func (a A) GetI() int {
	return a.i
}

func (a *A) IncrI() {
	a.i++
}
