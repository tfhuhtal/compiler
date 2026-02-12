package interpreter

import (
	"compiler/ast"
	"compiler/utils"
	"fmt"
)

type Value = any

type SymTab = utils.SymTab[Value]

type breakSignal struct{}
type continueSignal struct{}

type returnSignal struct {
	value Value
}

type userFunc struct {
	params []string
	body   ast.Expression
	symTab *SymTab
}

func interpret(node ast.Expression, symTab *SymTab) Value {
	switch n := node.(type) {

	case ast.Module:
		for _, fn := range n.Functions {
			fd := fn.(ast.FunctionDefinition)
			name := fd.Name.(ast.Identifier).Name
			var paramNames []string
			for _, p := range fd.Params {
				paramNames = append(paramNames, p.(ast.Param).Name.(ast.Identifier).Name)
			}
			symTab.Table[name] = userFunc{
				params: paramNames,
				body:   fd.Body,
				symTab: symTab,
			}
		}
		return interpret(n.Block, symTab)

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
		} else if n.Else != nil {
			e := n.Else.(ast.Block)
			for _, expr := range e.Expressions {
				interpret(expr, symTab)
			}
			return interpret(e.Result, symTab)
		}
		return nil

	case ast.Declaration:
		value := interpret(n.Value, symTab)
		var str string
		if identifier, ok := n.Variable.(ast.Identifier); ok {
			str = identifier.Name
		}
		if _, exists := symTab.Table[str]; exists {
			panic(fmt.Sprintf("%s already declared", n.Variable))
		}
		if n.Typed != nil {
			if n.Typed.(ast.Identifier).Name == "Bool" {
				if _, ok := value.(bool); !ok {
					panic("Must be boolean")
				}
			} else if n.Typed.(ast.Identifier).Name == "Int" {
				if _, ok := value.(uint64); !ok {
					panic("Must be integer")
				}
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
		name := n.Name.(ast.Identifier).Name
		if name == "print_int" {
			for _, a := range n.Args {
				v := interpret(a, symTab).(uint64)
				fmt.Println(v)
				return v
			}
		} else if name == "print_bool" {
			for _, a := range n.Args {
				v := interpret(a, symTab).(bool)
				fmt.Println(v)
				return v
			}
		} else if name == "read_int" {
			return uint64(0)
		}
		// Look up user-defined function
		fnVal := interpret(n.Name, symTab)
		uf, ok := fnVal.(userFunc)
		if !ok {
			panic(fmt.Sprintf("Not a function: %s", name))
		}
		fnTab := utils.NewSymTab(uf.symTab)
		for i, paramName := range uf.params {
			fnTab.Table[paramName] = interpret(n.Args[i], symTab)
		}
		var result Value
		func() {
			defer func() {
				if r := recover(); r != nil {
					if ret, ok := r.(returnSignal); ok {
						result = ret.value
					} else {
						panic(r)
					}
				}
			}()
			result = interpret(uf.body, fnTab)
		}()
		return result

	case ast.Block:
		tab := utils.NewSymTab(symTab)
		for _, expr := range n.Expressions {
			_ = interpret(expr, tab)
		}
		return interpret(n.Result, tab)

	case ast.WhileLoop:
		block := n.Looping.(ast.Block)
		for interpret(n.Condition, symTab).(bool) {
			brk := false
			cont := false
			allExprs := append(block.Expressions, block.Result)
			for _, expr := range allExprs {
				if expr == nil {
					continue
				}
				func() {
					defer func() {
						if r := recover(); r != nil {
							switch r.(type) {
							case breakSignal:
								brk = true
							case continueSignal:
								cont = true
							default:
								panic(r)
							}
						}
					}()
					_ = interpret(expr, symTab)
				}()
				if brk || cont {
					break
				}
			}
			if brk {
				break
			}
		}
		return nil

	case ast.BreakExpression:
		panic(breakSignal{})

	case ast.ContinueExpression:
		panic(continueSignal{})

	case ast.ReturnExpression:
		val := interpret(n.Result, symTab)
		panic(returnSignal{value: val})
	}
	return nil
}

func Interpret(nodes ast.Expression) Value {
	tab := utils.NewSymTab[Value](nil)
	res := interpret(nodes, tab)
	return res
}
