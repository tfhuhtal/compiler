package utils

type Type interface {
	isType()
}

type Int struct {
	Name string
}

func (Int) isType() {}

type Bool struct {
	Name string
}

func (Bool) isType() {}

type Fun struct {
	Params []Type
	Res    Type
}

func (Fun) isType() {}

type Unit struct {
	Name string
}

func (Unit) isType() {}

type SymTab[T any] struct {
	Parent *SymTab[T]
	Table  map[string]T
}

func NewSymTab[T any](parent *SymTab[T]) *SymTab[T] {
	return &SymTab[T]{
		Parent: parent,
		Table:  make(map[string]T),
	}
}
