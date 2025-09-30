package ast

import (
	"compiler/tokenizer"
	"compiler/utils"
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

type Module struct {
	Functions []Function
	Block     Block
}

type Param struct {
	Name string
	Type utils.Type
}

type Function struct {
	Name       Identifier
	Params     []Param
	ResultType utils.Type
	Body       Block
	Location   Location
}
