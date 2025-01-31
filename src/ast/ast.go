package ast

type Expression interface {
	isExpression()
}

type Literal struct {
	Value any
}

func (Literal) isExpression() {}

type Identifier struct {
	Name string
}

func (Identifier) isExpression() {}

type BinaryOp struct {
	Left  Expression
	Op    string
	Right Expression
}

func (BinaryOp) isExpression() {}

type IfExpression struct {
	Condition Expression
	Then      Expression
	Else      Expression
}

func (IfExpression) isExpression() {}

type FunctionCall struct {
	Name string
	Args []Expression
}

func (FunctionCall) isExpression() {}