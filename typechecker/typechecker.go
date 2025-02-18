package typechecker

import (
	"compiler/ast"
	"compiler/utils"
	"fmt"
)

type SymTab struct {
	Parent *SymTab
	Table  map[string]utils.Type
}

func NewSymTab(parent *SymTab) SymTab {
	return SymTab{
		Parent: parent,
		Table:  make(map[string]utils.Type),
	}
}

func typecheck(node ast.Expression, symTab *SymTab) utils.Type {
	switch n := node.(type) {
	case ast.Literal:
		_, ok := n.Value.(int)
		if ok {
			return utils.Int{
				Name: "Int",
			}
		} else if n.Value == nil {
			fmt.Println("Nil literal")
			return utils.Unit{
				Name: "Nil",
			}
		} else {
			panic(fmt.Sprintf("Unknown literal type %s at location", n.Value))
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
				panic(fmt.Sprintf("Both left %s and right %s must be integers", left, right))
			}
			return utils.Int{
				Name: "Int",
			}

		case "<", ">", ">=", "<=", "==", "!=":
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
		}

	case ast.IfExpression:
		condition := typecheck(n.Condition, symTab)
		_, ok := condition.(utils.Bool)
		if !ok {
			panic(fmt.Sprintf("%s condition is not boolean %v", condition, n))
		}
		then := typecheck(n.Then, symTab)
		typecheck(n.Else, symTab)
		return then

	case ast.Declaration:
		value := typecheck(n.Value, symTab)
		var str string
		if identifier, ok := n.Variable.(ast.Identifier); ok {
			str = identifier.Name
		}
		if _, exists := symTab.Table[str]; exists {
			panic(fmt.Sprintf("%s already declared", n.Variable))
		}
		symTab.Table[str] = value
		return value

	case ast.Identifier:
		fmt.Println("Identifier", symTab.Table)
		if value, exists := symTab.Table[n.Name]; exists {
			return value
		}
		cur_scp := symTab.Parent
		for cur_scp != nil {
			if value, exists := cur_scp.Table[n.Name]; exists {
				return value
			}
			cur_scp = cur_scp.Parent
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
		return utils.Fun{
			Params: exprs,
			Res:    res,
		}

	case ast.WhileLoop:
		cond := typecheck(n.Condition, symTab)
		if _, ok := cond.(utils.Type); !ok {
			panic(fmt.Sprintf("%s condition is not boolean", cond))
		}
		return typecheck(n.Looping, symTab)

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
