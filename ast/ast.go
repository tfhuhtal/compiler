package ast

import (
	"compiler/tokenizer"
	"compiler/utils"
)

type Location = tokenizer.SourceLocation

type Expression interface {
	isExpression()
	GetLocation() Location
	GetType() utils.Type
}

type Literal struct {
	Value    any
	Location Location
	Type     utils.Type
}

func (Literal) isExpression() {}
func (l Literal) GetLocation() Location {
	return l.Location
}

type Identifier struct {
	Name     string
	Location Location
	Type     utils.Type
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
	Type     utils.Type
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
	Type      utils.Type
}

func (IfExpression) isExpression() {}
func (i IfExpression) GetLocation() Location {
	return i.Location
}

type Function struct {
	Name     Expression
	Args     []Expression
	Location Location
	Type     utils.Type
}

func (Function) isExpression() {}
func (f Function) GetLocation() Location {
	return f.Location
}

type BooleanLiteral struct {
	Boolean  string
	Location Location
	Type     utils.Type
}

func (BooleanLiteral) isExpression() {}
func (b BooleanLiteral) GetLocation() Location {
	return b.Location
}

type Unary struct {
	Ops      []string
	Exp      Expression
	Location Location
	Type     utils.Type
}

func (Unary) isExpression() {}
func (u Unary) GetLocation() Location {
	return u.Location
}

type Block struct {
	Expressions []Expression
	Result      Expression
	Location    Location
	Type        utils.Type
}

func (Block) isExpression() {}
func (b Block) GetLocation() Location {
	return b.Location
}

type FunctionTypeExpression struct {
	VariableTypes []Expression
	ResultType    Expression
	Location      Location
	Type          utils.Type
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
	Type     utils.Type
}

func (Declaration) isExpression() {}
func (d Declaration) GetLocation() Location {
	return d.Location
}

type WhileLoop struct {
	Condition Expression
	Looping   Expression
	Location  Location
	Type      utils.Type
}

func (WhileLoop) isExpression() {}
func (w WhileLoop) GetLocation() Location {
	return w.Location
}

func (l Literal) GetType() utils.Type {
	return l.Type
}

func (i Identifier) GetType() utils.Type {
	return i.Type
}

func (b BinaryOp) GetType() utils.Type {
	return b.Type
}

func (i IfExpression) GetType() utils.Type {
	return i.Type
}

func (f Function) GetType() utils.Type {
	return f.Type
}

func (b BooleanLiteral) GetType() utils.Type {
	return b.Type
}

func (u Unary) GetType() utils.Type {
	return u.Type
}

func (b Block) GetType() utils.Type {
	return b.Type
}

func (f FunctionTypeExpression) GetType() utils.Type {
	return f.Type
}

func (d Declaration) GetType() utils.Type {
	return d.Type
}

func (w WhileLoop) GetType() utils.Type {
	return w.Type
}
