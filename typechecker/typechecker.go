package typechecker

import (
	"compiler/ast"
	"compiler/utils"
	"fmt"
)

type SymTab = utils.SymTab[utils.Type]

func resolveType(name string) utils.Type {
	switch name {
	case "Int":
		return utils.Int{Name: "Int"}
	case "Bool":
		return utils.Bool{Name: "Bool"}
	case "Unit":
		return utils.Unit{Name: "Unit"}
	default:
		panic(fmt.Sprintf("Unknown type: %s", name))
	}
}

func typecheck(node ast.Expression, symTab *SymTab) utils.Type {
	switch n := node.(type) {
	case ast.Module:
		// First pass: register all function types for mutual recursion
		for _, fn := range n.Functions {
			fd := fn.(ast.FunctionDefinition)
			name := fd.Name.(ast.Identifier).Name
			var paramTypes []utils.Type
			for _, p := range fd.Params {
				paramTypes = append(paramTypes, resolveType(p.(ast.Param).Type.(ast.Identifier).Name))
			}
			retType := resolveType(fd.ResultType.(ast.Identifier).Name)
			symTab.Table[name] = utils.Fun{
				Params: paramTypes,
				Res:    retType,
			}
		}
		// Second pass: type-check function bodies
		for _, fn := range n.Functions {
			fd := fn.(ast.FunctionDefinition)
			fnTab := utils.NewSymTab(symTab)
			for _, p := range fd.Params {
				param := p.(ast.Param)
				pName := param.Name.(ast.Identifier).Name
				pType := resolveType(param.Type.(ast.Identifier).Name)
				fnTab.Table[pName] = pType
			}
			typecheck(fd.Body, fnTab)
		}
		return typecheck(n.Block, symTab)

	case ast.Literal:
		var res utils.Type
		_, ok := n.Value.(uint64)
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
		return res

	case ast.BinaryOp:
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

		case "!=", "==", "and", "or":
			if left != right {
				panic(fmt.Sprintf("Both left %s and right %s must be same type", left, right))
			}
			return utils.Bool{
				Name: "Bool",
			}
		}

	case ast.IfExpression:
		condition := typecheck(n.Condition, symTab)
		_, ok := condition.(utils.Bool)
		if !ok {
			panic(fmt.Sprintf("%s condition is not boolean", condition))
		}
		if n.Then == nil {
			panic(fmt.Sprintf("Not allowed to declare here %v", n.Then))
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
		if n.Typed.(ast.Identifier).Name == "Bool" {
			if _, ok := value.(utils.Bool); !ok {
				panic("Must be boolean")
			}
		} else if n.Typed.(ast.Identifier).Name == "Int" {
			if _, ok := value.(utils.Int); !ok {
				panic("Must be integer")
			}
		}
		symTab.Table[str] = value
		return value

	case ast.Identifier:
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
		if _, ok := value.(utils.Bool); !ok && n.Op == "not" {
			panic(fmt.Sprintf("Not allowed Unary %v", value))
		}
		return value

	case ast.BooleanLiteral:
		var res utils.Type
		if n.Boolean == "true" || n.Boolean == "false" {
			res = utils.Bool{
				Name: "Bool",
			}
		}
		return res

	case ast.FunctionCall:
		var argTypes []utils.Type
		for _, par := range n.Args {
			argTypes = append(argTypes, typecheck(par, symTab))
		}
		name := n.Name.(ast.Identifier).Name
		if name == "print_int" {
			if _, ok := argTypes[len(argTypes)-1].(utils.Int); !ok {
				panic("params should be int")
			}
			return utils.Int{Name: "Int"}
		} else if name == "print_bool" {
			if _, ok := argTypes[len(argTypes)-1].(utils.Bool); !ok {
				panic("params should be bool")
			}
			return utils.Bool{Name: "Bool"}
		} else if name == "read_int" {
			return utils.Int{Name: "Int"}
		}
		fnType := typecheck(n.Name, symTab)
		if ft, ok := fnType.(utils.Fun); ok {
			if len(ft.Params) != len(argTypes) {
				panic(fmt.Sprintf("Function %s expects %d args, got %d", name, len(ft.Params), len(argTypes)))
			}
			for i, pt := range ft.Params {
				if pt != argTypes[i] {
					panic(fmt.Sprintf("Argument %d type mismatch: expected %v, got %v", i, pt, argTypes[i]))
				}
			}
			return ft.Res
		}
		return fnType

	case ast.Block:
		var exprs []utils.Type
		tab := utils.NewSymTab(symTab)
		for _, expr := range n.Expressions {
			exprs = append(exprs, typecheck(expr, tab))
		}
		res := typecheck(n.Result, tab)
		return res

	case ast.WhileLoop:
		cond := typecheck(n.Condition, symTab)
		if _, ok := cond.(utils.Bool); !ok {
			panic(fmt.Sprintf("%s condition is not boolean", cond))
		}
		typecheck(n.Looping, symTab)
		return utils.Unit{}

	case ast.BreakExpression:
		return utils.Unit{}

	case ast.ContinueExpression:
		return utils.Unit{}

	case ast.ReturnExpression:
		return typecheck(n.Result, symTab)

	}
	return utils.Unit{}
}

func Type(nodes ast.Expression) any {
	tab := utils.NewSymTab[utils.Type](nil)
	res := typecheck(nodes, tab)
	return res
}
