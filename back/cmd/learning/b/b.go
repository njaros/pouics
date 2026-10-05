package b

type A interface {
	GetI() int
	IncrI()
}

type B struct {
	a A
}

func NewB(a A) *B {
	return &B{a}
}

func (b *B) GetA() A {
	return b.a
}

func (b *B) SetA(a A) {
	b.a = a
}