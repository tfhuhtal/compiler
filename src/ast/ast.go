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
	Name Expression
	Args []Expression
}

func (FunctionCall) isExpression() {}

type UnaryOp struct {
	Op    string
	Right Expression
}

func (UnaryOp) isExpression() {}

type Assignment struct {
	Left 	Expression
	Right	Expression
}

func (Assignment) isExpression() {}

type Block struct {
	Expressions []Expression
}

func (Block) isExpression() {}
