package ast

import "compiler/tokenizer"

type Location = tokenizer.SourceLocation

type Expression interface {
	isExpression()
	GetLocation() Location
}

type Literal struct {
	Value    any
	Location Location
}

func (Literal) isExpression() {}
func (l Literal) GetLocation() Location {
	return l.Location
}

type Identifier struct {
	Name     string
	Location Location
}

func (Identifier) isExpression() {}
func (i Identifier) GetLocation() Location {
	return i.Location
}

type BinaryOp struct {
	Left     Expression
	Op       string
	Right    Expression
	Location Location
}

func (BinaryOp) isExpression() {}
func (b BinaryOp) GetLocation() Location {
	return b.Location
}

type IfExpression struct {
	Condition Expression
	Then      Expression
	Else      Expression
	Location  Location
}

func (IfExpression) isExpression() {}
func (i IfExpression) GetLocation() Location {
	return i.Location
}

type Function struct {
	Name     Expression
	Args     []Expression
	Location Location
}

func (Function) isExpression() {}
func (f Function) GetLocation() Location {
	return f.Location
}

type BooleanLiteral struct {
	Boolean  string
	Location Location
}

func (BooleanLiteral) isExpression() {}
func (b BooleanLiteral) GetLocation() Location {
	return b.Location
}

type Unary struct {
	Ops      []string
	Exp      Expression
	Location Location
}

func (Unary) isExpression() {}
func (u Unary) GetLocation() Location {
	return u.Location
}

type Block struct {
	Expressions []Expression
	Result      Expression
	Location    Location
}

func (Block) isExpression() {}
func (b Block) GetLocation() Location {
	return b.Location
}

type FunctionTypeExpression struct {
	VariableTypes []Expression
	ResultType    Expression
	Location      Location
}

func (FunctionTypeExpression) isExpression() {}
func (f FunctionTypeExpression) GetLocation() Location {
	return f.Location
}

type Declaration struct {
	Variable Expression
	Value    Expression
	Typed    Expression
	Location Location
}

func (Declaration) isExpression() {}
func (d Declaration) GetLocation() Location {
	return d.Location
}

type WhileLoop struct {
	Condition Expression
	Looping   Expression
	Location  Location
}

func (WhileLoop) isExpression() {}
func (w WhileLoop) GetLocation() Location {
	return w.Location
}
