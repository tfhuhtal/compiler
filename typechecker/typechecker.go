package typechecker

import (
	"compiler/ast"
	"compiler/utils"
	"fmt"
)

type SymTab = utils.SymTab[utils.Type]

func typecheck(node ast.Expression, symTab *SymTab) utils.Type {
	switch n := node.(type) {
	case *ast.Literal:
		var res utils.Type
		_, ok := n.Value.(int)
		if ok {
			res = utils.Int{
				Name: "Int",
			}
		} else if n.Value == nil {
			res = utils.Unit{
				Name: "Nil",
			}
		} else {
			panic(fmt.Sprintf("Unknown literal type %s at location", n.Value))
		}
		n.Type = res
		return res

	case *ast.BinaryOp:
		left := typecheck(n.Left, symTab)
		right := typecheck(n.Right, symTab)
		n.Type = left

		switch n.Op {
		case "+", "-", "*", "/", "%":
			_, lok := left.(utils.Int)
			_, rok := right.(utils.Int)
			if !lok || !rok {
				panic(fmt.Sprintf("Both left %s and right %s must be integers", left, right))
			}
			return utils.Int{
				Name: "Int",
			}

		case "<", ">", ">=", "<=":
			_, lok := left.(utils.Int)
			_, rok := right.(utils.Int)
			if !lok || !rok {
				panic(fmt.Sprintf("Both left %s and right %s must be integers", left, right))
			}
			return utils.Bool{
				Name: "Bool",
			}

		case "=":
			if left != right {
				panic(fmt.Sprintf("Both left %s and right %s must be same type", left, right))
			}
			return left

		case "!=", "==":
			if left != right {
				panic(fmt.Sprintf("Both left %s and right %s must be same type", left, right))
			}
			return utils.Bool{
				Name: "Bool",
			}
		}

	case *ast.IfExpression:
		condition := typecheck(n.Condition, symTab)
		_, ok := condition.(utils.Bool)
		if !ok {
			panic(fmt.Sprintf("%s condition is not boolean", condition))
		}
		then := typecheck(n.Then, symTab)
		typecheck(n.Else, symTab)
		n.Type = then
		return then

	case *ast.Declaration:
		value := typecheck(n.Value, symTab)
		var str string
		if identifier, ok := n.Variable.(ast.Identifier); ok {
			str = identifier.Name
		}
		if _, exists := symTab.Table[str]; exists {
			panic(fmt.Sprintf("%s already declared", n.Variable))
		}
		symTab.Table[str] = value
		n.Type = value
		return value

	case *ast.Identifier:
		if value, exists := symTab.Table[n.Name]; exists {
			n.Type = value
			return value
		}
		cur_scp := symTab.Parent
		for cur_scp != nil {
			if value, exists := cur_scp.Table[n.Name]; exists {
				n.Type = value
				return value
			}
			cur_scp = cur_scp.Parent
		}

	case *ast.Unary:
		value := typecheck(n.Exp, symTab)
		n.Type = value
		return value

	case *ast.BooleanLiteral:
		var res utils.Type
		if n.Boolean == "true" || n.Boolean == "false" {
			res = utils.Bool{
				Name: "Bool",
			}
		}
		n.Type = res
		return res

	case *ast.Function:
		var params []utils.Type
		for _, par := range n.Args {
			params = append(params, typecheck(par, symTab))
		}
		res := typecheck(n.Name, symTab)
		n.Type = res
		return utils.Fun{
			Params: params,
			Res:    res,
		}

	case *ast.Block:
		var exprs []utils.Type
		tab := utils.NewSymTab(symTab)
		for _, expr := range n.Expressions {
			exprs = append(exprs, typecheck(expr, tab))
		}
		res := typecheck(n.Result, tab)
		n.Type = res
		return utils.Fun{
			Params: exprs,
			Res:    res,
		}

	case *ast.WhileLoop:
		cond := typecheck(n.Condition, symTab)
		if _, ok := cond.(utils.Bool); !ok {
			panic(fmt.Sprintf("%s condition is not boolean", cond))
		}
		n.Type = cond
		return typecheck(n.Looping, symTab)

	case *ast.FunctionTypeExpression:
		return nil
	}
	return nil
}

func Type(nodes ast.Expression) any {
	tab := utils.NewSymTab[utils.Type](nil)
	res := typecheck(nodes, tab)
	return res
}
