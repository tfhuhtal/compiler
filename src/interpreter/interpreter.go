package interpreter

import (
	"compiler/src/ast"
)

type Value any

func Interpret(node ast.Expression) Value {
	switch n := node.(type) {
	case ast.Literal:
		return n.Value
	case ast.BinaryOp:
		a := Interpret(n.Left)
		b := Interpret(n.Right)
		if n.Op == "+" {
			return a.(int) + b.(int)
		} else if n.Op == "<" {
			return a.(int) < b.(int)
		} else {
			panic("...")
		}
	case ast.IfExpression:
		condition := Interpret(n.Condition)
		if conditionBool, ok := condition.(bool); ok {
			if conditionBool {
				return Interpret(n.Then)
			} else {
				return Interpret(n.Else)
			}
		} else {
			panic("Condition is not a boolean")
		}
	default:
		return nil
	}
}
