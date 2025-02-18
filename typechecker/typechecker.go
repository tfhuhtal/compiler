package typechecker

import (
	"compiler/ast"
	"compiler/utils"
	"fmt"
)

type SymTab struct {
	Parent *SymTab
	Table  map[any]utils.Type
}

func NewSymTab(parent *SymTab) SymTab {
	return SymTab{
		Parent: parent,
		Table:  make(map[any]utils.Type),
	}
}

func typecheck(node ast.Expression, symTab *SymTab) utils.Type {
	switch n := node.(type) {
	case ast.Literal:
		fmt.Println("Literal")
		_, ok := n.Value.(int)
		if ok {
			return utils.Int{
				Name: "Int",
			}
		} else {
			panic("Unknown literal type")
		}

	case ast.BinaryOp:
		fmt.Println("BinaryOp")
		left := typecheck(n.Left, symTab)
		right := typecheck(n.Right, symTab)

		switch n.Op {
		case "+", "-", "*", "/", "%":
			_, lok := left.(utils.Int)
			_, rok := right.(utils.Int)
			if !lok || !rok {
				panic("Both left and right must be integers")
			}
			return utils.Int{
				Name: "Int",
			}

		case "<", ">", ">=", "<=":
			_, lok := left.(utils.Int)
			_, rok := right.(utils.Int)
			if !lok || !rok {
				panic("Both left and right must be integers")
			}
			return utils.Bool{
				Name: "Int",
			}

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
		fmt.Println("IfExpression")
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
		fmt.Println("Declaration")
		value := typecheck(n.Value, symTab)
		if _, exists := symTab.Table[n.Variable]; exists {
			panic(fmt.Sprintf("%s already declared", n.Variable))
		}
		symTab.Table[n.Variable] = value
		return value

	case ast.Identifier:
		fmt.Println("Identifier")
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

	case ast.Unary:
		value := typecheck(n.Exp, symTab)
		return value

	case ast.BooleanLiteral:
		if n.Boolean == "true" || n.Boolean == "false" {
			return utils.Bool{
				Name: "Bool",
			}
		}

	case ast.Function:
		fmt.Println("Function")
		var params []utils.Type
		for _, par := range n.Args {
			params = append(params, typecheck(par, symTab))
		}
		res := typecheck(n.Name, symTab)
		return utils.Fun{
			Params: params,
			Res:    res,
		}

	case ast.Block:
		var exprs []utils.Type
		tab := NewSymTab(symTab)
		for _, expr := range n.Expressions {
			exprs = append(exprs, typecheck(expr, &tab))
		}
		res := typecheck(n.Result, &tab)
		fmt.Println("Block", res)
		return utils.Fun{
			Params: exprs,
			Res:    res,
		}

	case ast.WhileLoop:
		return nil

	case ast.FunctionTypeExpression:
		return nil
	}
	return nil
}

func Type(nodes []ast.Expression) any {
	tab := NewSymTab(nil)
	var res []any

	for _, node := range nodes {
		res = append(res, typecheck(node, &tab))
	}
	return res
}
