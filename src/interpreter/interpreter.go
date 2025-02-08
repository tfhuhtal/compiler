package interpreter

import (
	"compiler/src/ast"
)

type Value any

func interpret(node ast.Expression) Value {
	switch n := node.(type) {
	case ast.Literal:
		return n.Value
	case ast.BinaryOp:
		a := interpret(n.Left)
		b := interpret(n.Right)
		if n.Op == "+" {
			return a.(int) + b.(int)
		} else if n.Op == "<" {
			return a.(int) < b.(int)
		} else {
			panic("...")
		}
	case ast.IfExpression:
		if interpret(n.Condition) {
			return interpret(n.Then)
		} else {
			return interpret(n.Else)
		}
	default:
		return nil
	}
}
