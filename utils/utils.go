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
