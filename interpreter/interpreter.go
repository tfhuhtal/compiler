package interpreter

import (
	"compiler/ast"
	"compiler/utils"
	"fmt"
)

type Value = any

type SymTab = utils.SymTab[Value]

func interpret(node ast.Expression, symTab *SymTab) Value {
	switch n := node.(type) {

	case ast.Literal:

		value, ok := n.Value.(uint64)
		if ok {
			return value
		} else {
			panic(fmt.Sprintf("Unknown literal type %s", n.Value))
		}

	case ast.BinaryOp:
		left := interpret(n.Left, symTab)
		right := interpret(n.Right, symTab)

		switch n.Op {
		case "+":
			return left.(uint64) + right.(uint64)
		case "-":
			return left.(uint64) - right.(uint64)
		case "*":
			return left.(uint64) * right.(uint64)
		case "/":
			return left.(uint64) / right.(uint64)
		case "%":
			return left.(uint64) % right.(uint64)
		case "<":
			return left.(uint64) < right.(uint64)
		case ">":
			return left.(uint64) > right.(uint64)
		case ">=":
			return left.(uint64) >= right.(uint64)
		case "<=":
			return left.(uint64) <= right.(uint64)
		case "!=":
			return left.(uint64) != right.(uint64)
		case "==":
			return left.(uint64) == right.(uint64)
		case "and":
			return left.(bool) && right.(bool)
		case "or":
			return left.(bool) || right.(bool)
		case "=": // asserttion handled differently
			symTab.Table[n.Left.(ast.Identifier).Name] = right.(uint64)
			return right
		}

	case ast.IfExpression:
		if interpret(n.Condition, symTab).(bool) {
			t := n.Then.(ast.Block)
			for _, expr := range t.Expressions {
				interpret(expr, symTab)
			}
			return interpret(t.Result, symTab)
		} else {
			e := n.Else.(ast.Block)
			for _, expr := range e.Expressions {
				interpret(expr, symTab)
			}
			return interpret(e.Result, symTab)
		}

	case ast.Declaration:
		value := interpret(n.Value, symTab)
		var str string
		if identifier, ok := n.Variable.(ast.Identifier); ok {
			str = identifier.Name
		}
		if _, exists := symTab.Table[str]; exists {
			panic(fmt.Sprintf("%s already declared", n.Variable))
		}
		if n.Typed.(ast.Identifier).Name == "Bool" {
			if _, ok := value.(bool); !ok {
				panic("Must be boolean")
			}
		} else if n.Typed.(ast.Identifier).Name == "Int" {
			if _, ok := value.(uint64); !ok {
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
		value := interpret(n.Exp, symTab)
		if _, ok := value.(utils.Bool); !ok && n.Op == "not" {
			panic(fmt.Sprintf("Not allowed Unary %v", value))
		}
		return value

	case ast.BooleanLiteral:
		if n.Boolean == "true" {
			return true
		}
		return false

	case ast.FunctionCall:
		if n.Name.(ast.Identifier).Name == "print_int" {
			for _, a := range n.Args {
				v := interpret(a, symTab).(uint64)
				fmt.Println(v)
				return v
			}
		} else if n.Name.(ast.Identifier).Name == "print_bool" {
			for _, a := range n.Args {
				v := interpret(a, symTab).(bool)
				fmt.Println(v)
				return v
			}
		} else if n.Name.(ast.Identifier).Name == "read_int" {
		}
		return interpret(n.Name, symTab)

	case ast.Block:
		tab := utils.NewSymTab(symTab)
		for _, expr := range n.Expressions {
			_ = interpret(expr, tab)
		}
		return interpret(n.Result, tab)

	case ast.WhileLoop:
		block := n.Looping.(ast.Block)
		for interpret(n.Condition, symTab).(bool) {
			loop := block.Expressions
			for _, expr := range loop {
				_ = interpret(expr, symTab)
			}
		}
		return interpret(block.Result, symTab)
	}
	return nil
}

func Interpret(nodes ast.Expression) Value {
	tab := utils.NewSymTab[Value](nil)
	res := interpret(nodes, tab)
	return res
}
