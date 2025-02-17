package typechecker

import (
	"compiler/src/ast"
	"compiler/src/utils"
	"fmt"
)

func typecheck(node ast.Expression, symTab *utils.SymTab) utils.Type {
	switch n := node.(type) {
	case ast.Literal:
		_, ok := n.Value.(int)
		if ok {
			return utils.Int{}
		} else {
			panic("Unknown literal type")
		}

	case ast.BinaryOp:
		left := typecheck(n.Left, symTab)
		right := typecheck(n.Right, symTab)

		switch n.Op {
		case "+", "-", "*", "/", "%":
			_, lok := left.(utils.Int)
			_, rok := right.(utils.Int)
			if !lok || !rok {
				panic("Both left and right must be integers")
			}
			return utils.Int{}

		case "<", ">", ">=", "<=":
			_, lok := left.(utils.Int)
			_, rok := right.(utils.Int)
			if !lok || !rok {
				panic("Both left and right must be integers")
			}
			return utils.Bool{}

		case "=":
			if left != right {
				panic("Variables are not same type")
			}
			return left

		case "==", "!=":
			if left != right {
				panic("variables are not same type")
			}
			return left
		}

	case ast.IfExpression:
		condition := typecheck(n.Condition, symTab)
		_, ok := condition.(utils.Bool)
		if !ok {
			panic("The condition is not boolean")
		}
		then := typecheck(n.Then, symTab)
		els := typecheck(n.Else, symTab)
		if then != els {
			panic("In if clause then and else are not same type")
		}
		return then

	case ast.Declaration:
		value := typecheck(n.Value, symTab)
		if _, exists := symTab.Table[n.Variable]; exists {
			panic(fmt.Sprintf("%s already declared", n.Variable))
		}
		symTab.Table[n.Variable] = value
		return value

	case ast.Identifier:
		if symTab != nil {
			if value, exists := symTab.Table[n.Name]; exists {
				return typecheck(value.(ast.Expression), symTab)
			}
			cur_scp := symTab.Parent
			for cur_scp != nil {
				if value, exists := cur_scp.Table[n.Name]; exists {
					return typecheck(value.(ast.Expression), symTab)
				}
				cur_scp = cur_scp.Parent
			}
		}
	}
	return nil
}
