package ast

import (
	"compiler/tokenizer"
)

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

type BooleanLiteral struct {
	Boolean  string
	Location Location
}

func (BooleanLiteral) isExpression() {}
func (b BooleanLiteral) GetLocation() Location {
	return b.Location
}

type Unary struct {
	Op       string
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

type FunctionCall struct {
	Name     Expression
	Args     []Expression
	Location Location
}

func (FunctionCall) isExpression() {}
func (f FunctionCall) GetLocation() Location {
	return f.Location
}

type Module struct {
	Functions []Expression
	Block     Expression
	Location  Location
}

// Module is not actually expression but the sake of GO it has to be done like this
func (Module) isExpression() {}
func (m Module) GetLocation() Location {
	return m.Location
}

type Param struct {
	Name     Expression
	Type     Expression
	Location Location
}

func (Param) isExpression() {}
func (p Param) GetLocation() Location {
	return p.Location
}

type ReturnExpression struct {
	Result   Expression
	Location Location
}

func (ReturnExpression) isExpression() {}
func (r ReturnExpression) GetLocation() Location {
	return r.Location
}

// Funtion is not actually expression but the sake of GO it has to be done like this
type FunctionDefinition struct {
	Name       Expression
	Params     []Expression
	ResultType Expression
	Body       Expression
	Location   Location
}

func (FunctionDefinition) isExpression() {}
func (f FunctionDefinition) GetLocation() Location {
	return f.Location
}
