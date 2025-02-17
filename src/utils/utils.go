package utils

type SymTab struct {
	Parent *SymTab
	Table  map[any]any
}

type Type interface {
	isType()
}

type Int struct {
	name string
}

func (Int) isType() {}

type Bool struct {
	name string
}

func (Bool) isType() {}

type Unit struct {
	name string
}

func (Unit) isType() {}
