package utils

type SymTab struct {
	Parent *SymTab
	Table  map[any]any
}

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
