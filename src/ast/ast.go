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

type Function struct {
	Name Expression
	Args []Expression
}

func (Function) isExpression() {}

type BooleanLiteral struct {
	Boolean string
}

func (BooleanLiteral) isExpression() {}

type Unary struct {
	Ops []string
	Exp Expression
}

func (Unary) isExpression() {}

type Block struct {
	Expressions []Expression
	Result      Expression
}

func (Block) isExpression() {}

type FunctionTypeExpression struct {
	VariableTypes []Expression
	ResultType    Expression
}

func (FunctionTypeExpression) isExpression() {}

type Declaration struct {
	Variable Expression
	Value    Expression
	Typed    Expression
}

func (Declaration) isExpression() {}

type WhileLoop struct {
	Condition Expression
	Looping   Expression
}

func (WhileLoop) isExpression() {}
