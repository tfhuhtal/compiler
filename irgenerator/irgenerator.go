package irgenerator

import (
	"compiler/ast"
	"compiler/ir"
	"compiler/utils"
	"fmt"
	"strconv"
)

type IRVar = ir.IRVar
type Type = utils.Type
type SymTab = utils.SymTab[IRVar]

type IRGenerator struct {
	varTypes       map[IRVar]Type
	rootTypes      map[IRVar]Type
	instructions   []ir.Instruction
	loopStartLabel *ir.Label
	loopEndLabel   *ir.Label
}

func new(rootTypes map[IRVar]Type) *IRGenerator {
	gen := &IRGenerator{
		varTypes:     make(map[IRVar]Type),
		rootTypes:    rootTypes,
		instructions: []ir.Instruction{},
	}
	for k, v := range rootTypes {
		gen.varTypes[k] = v
	}
	gen.varTypes["unit"] = utils.Unit{}
	return gen
}

func Generate(rootExpr ast.Expression) map[string][]ir.Instruction {
	rootTypes := map[IRVar]utils.Type{
		"+":   utils.Int{},
		"*":   utils.Int{},
		"/":   utils.Int{},
		"%":   utils.Int{},
		"-":   utils.Int{},
		">":   utils.Bool{},
		"==":  utils.Bool{},
		"<=":  utils.Bool{},
		"<":   utils.Bool{},
		">=":  utils.Bool{},
		"!=":  utils.Bool{},
		"and": utils.Bool{},
		"or":  utils.Bool{},
	}

	funcs := make(map[string][]ir.Instruction)

	// Handle Module: generate IR for each function definition
	if mod, ok := rootExpr.(ast.Module); ok {
		rootSymTab := utils.NewSymTab[IRVar](nil)
		for v := range rootTypes {
			rootSymTab.Table[v] = v
		}

		for _, fn := range mod.Functions {
			fd := fn.(ast.FunctionDefinition)
			name := fd.Name.(ast.Identifier).Name
			rootSymTab.Table[name] = name
		}

		for _, fn := range mod.Functions {
			fd := fn.(ast.FunctionDefinition)
			name := fd.Name.(ast.Identifier).Name

			g := new(rootTypes)
			// Copy function names into this generator's varTypes
			for _, fn2 := range mod.Functions {
				fd2 := fn2.(ast.FunctionDefinition)
				n2 := fd2.Name.(ast.Identifier).Name
				g.varTypes[n2] = utils.Fun{}
			}

			fnSymTab := utils.NewSymTab(rootSymTab)
			for i, p := range fd.Params {
				param := p.(ast.Param)
				pName := param.Name.(ast.Identifier).Name
				pType := resolveIRType(param.Type.(ast.Identifier).Name)
				paramVar := g.newVar(pType)
				fnSymTab.Table[pName] = paramVar
				g.instructions = append(g.instructions, ir.LoadParam{
					BaseInstruction: ir.BaseInstruction{Location: param.GetLocation()},
					Index:           i,
					Dest:            paramVar,
				})
				fnSymTab.Table[pName] = paramVar
			}
			g.visit(fnSymTab, fd.Body)
			funcs[name] = g.instructions
		}

		// Generate main
		g := new(rootTypes)
		for _, fn := range mod.Functions {
			fd := fn.(ast.FunctionDefinition)
			name := fd.Name.(ast.Identifier).Name
			g.varTypes[name] = utils.Fun{}
		}
		mainSymTab := utils.NewSymTab(rootSymTab)
		result := g.visit(mainSymTab, mod.Block)
		emitTopLevelPrint(g, result, rootExpr)
		funcs["main"] = g.instructions
	} else {
		// No module, just a top-level expression
		g := new(rootTypes)
		rootSymTab := utils.NewSymTab[IRVar](nil)
		for v := range g.varTypes {
			rootSymTab.Table[v] = v
		}
		result := g.visit(rootSymTab, rootExpr)
		emitTopLevelPrint(g, result, rootExpr)
		funcs["main"] = g.instructions
	}

	return funcs
}

func resolveIRType(name string) utils.Type {
	switch name {
	case "Int":
		return utils.Int{Name: "Int"}
	case "Bool":
		return utils.Bool{Name: "Bool"}
	default:
		return utils.Unit{Name: name}
	}
}

func emitTopLevelPrint(g *IRGenerator, result IRVar, rootExpr ast.Expression) {
	if _, ok := g.varTypes[result].(utils.Int); ok {
		g.instructions = append(g.instructions, ir.Call{
			BaseInstruction: ir.BaseInstruction{Location: rootExpr.GetLocation()},
			Fun:             "print_int",
			Args:            []IRVar{result},
			Dest:            g.newVar(utils.Unit{}),
		})
	} else if _, ok := g.varTypes[result].(utils.Bool); ok {
		g.instructions = append(g.instructions, ir.Call{
			BaseInstruction: ir.BaseInstruction{Location: rootExpr.GetLocation()},
			Fun:             "print_bool",
			Args:            []IRVar{result},
			Dest:            g.newVar(utils.Unit{}),
		})
	}
}

func (g *IRGenerator) newVar(t Type) IRVar {
	idx := 0
	name := fmt.Sprintf("x%d", idx)
	for {
		if _, exists := g.varTypes[name]; !exists {
			break
		}
		idx++
		name = fmt.Sprintf("x%d", idx)
	}
	g.varTypes[name] = t
	return name
}

func (g *IRGenerator) newLabel() ir.Label {
	idx := 0
	name := fmt.Sprintf("L%d", idx)
	for {
		if _, exists := g.varTypes[name]; !exists {
			break
		}
		idx++
		name = fmt.Sprintf("L%d", idx)
	}
	g.varTypes[name] = utils.Unit{Name: name}
	return ir.Label{
		BaseInstruction: ir.BaseInstruction{},
		Label:           name,
	}
}

func (g *IRGenerator) visit(st *SymTab, expr ast.Expression) IRVar {
	switch e := expr.(type) {
	case ast.Literal:
		if value, ok := e.Value.(uint64); ok {
			variable := g.newVar(utils.Int{Name: "Int"})
			g.instructions = append(g.instructions, ir.LoadIntConst{
				BaseInstruction: ir.BaseInstruction{Location: e.GetLocation()},
				Value:           value,
				Dest:            variable,
			})
			return variable
		} else if e.Value == nil {
			return "unit"
		}
		panic(fmt.Sprintf("Unsupported literal: %v", e.Value))

	case ast.BooleanLiteral:
		if e.Boolean == "true" || e.Boolean == "false" {
			value, _ := strconv.ParseBool(e.Boolean)
			variable := g.newVar(utils.Bool{Name: "Bool"})
			g.instructions = append(g.instructions, ir.LoadBoolConst{
				BaseInstruction: ir.BaseInstruction{Location: e.GetLocation()},
				Value:           value,
				Dest:            variable,
			})
			return variable
		}
		panic("Unsupported boolean literal")

	case ast.Identifier:
		value, exists := st.Table[e.Name]

		if st.Parent != nil && !exists {
			return g.visit(st.Parent, e)
		} else if !exists {
			panic(fmt.Sprintf("Undefined variable: %s, in location %v", e.Name, e.GetLocation()))
		}
		return value

	case ast.BinaryOp:
		left := g.visit(st, e.Left)
		if e.Op == "=" {
			right := g.visit(st, e.Right)
			g.instructions = append(g.instructions, ir.Copy{
				BaseInstruction: ir.BaseInstruction{Location: e.GetLocation()},
				Source:          right,
				Dest:            left,
			})

			return left

		} else if e.Op == "and" || e.Op == "or" {
			rightLabel := g.newLabel()
			leftLabel := g.newLabel()

			if e.Op == "and" {
				g.instructions = append(g.instructions, ir.CondJump{
					BaseInstruction: ir.BaseInstruction{Location: e.GetLocation()},
					Cond:            left,
					ThenLabel:       rightLabel,
					ElseLabel:       leftLabel,
				})
			} else {
				g.instructions = append(g.instructions, ir.CondJump{
					BaseInstruction: ir.BaseInstruction{Location: e.GetLocation()},
					Cond:            left,
					ThenLabel:       leftLabel,
					ElseLabel:       rightLabel,
				})
			}

			g.instructions = append(g.instructions, rightLabel)
			right := g.visit(st, e.Right)
			newVar := g.newVar(utils.Bool{})
			g.instructions = append(g.instructions, ir.Copy{
				BaseInstruction: ir.BaseInstruction{Location: e.GetLocation()},
				Source:          right,
				Dest:            newVar,
			})
			endLabel := g.newLabel()
			g.instructions = append(g.instructions, ir.Jump{
				BaseInstruction: ir.BaseInstruction{Location: e.GetLocation()},
				Label:           endLabel,
			})

			g.instructions = append(g.instructions, leftLabel)
			var value bool
			if e.Op == "and" {
				value = false
			} else {
				value = true
			}

			g.instructions = append(g.instructions, ir.LoadBoolConst{
				BaseInstruction: ir.BaseInstruction{Location: e.Right.GetLocation()},
				Value:           value,
				Dest:            newVar,
			})

			g.instructions = append(g.instructions, ir.Jump{BaseInstruction: ir.BaseInstruction{}, Label: endLabel})
			g.instructions = append(g.instructions, endLabel)

			return newVar
		}
		right := g.visit(st, e.Right)
		varOp, exists := st.Table[e.Op]
		if !exists {
			panic(fmt.Sprintf("Unknown operator: %s", e.Op))
		}
		res := g.newVar(g.varTypes[varOp])
		g.instructions = append(g.instructions, ir.Call{
			BaseInstruction: ir.BaseInstruction{Location: e.GetLocation()},
			Fun:             varOp,
			Args:            []IRVar{left, right},
			Dest:            res,
		})
		return res

	case ast.Declaration:
		value := g.visit(st, e.Value)
		var name string
		if identifier, ok := e.Variable.(ast.Identifier); ok {
			name = identifier.Name
		}
		if _, exists := st.Table[name]; exists {
			panic(fmt.Sprintf("%v already declared", e.Variable))
		}
		newVar := g.newVar(g.varTypes[value])
		st.Table[name] = newVar
		g.instructions = append(g.instructions, ir.Copy{
			BaseInstruction: ir.BaseInstruction{Location: e.GetLocation()},
			Source:          value,
			Dest:            newVar,
		})
		return "unit"

	case ast.IfExpression:
		thenLabel := g.newLabel()
		endLabel := g.newLabel()
		elseLabel := endLabel
		if e.Else != nil {
			elseLabel = g.newLabel()
		}
		condVar := g.visit(st, e.Condition)
		copyVar := "unit"
		switch e.Then.(type) {
		case ast.Block:
			copyVar = g.newVar(utils.Unit{})
		case ast.Literal:
			copyVar = g.newVar(utils.Int{})
		case ast.BinaryOp:
			copyVar = g.newVar(utils.Int{})
		case ast.BooleanLiteral:
			copyVar = g.newVar(utils.Bool{})
		case ast.IfExpression:
			copyVar = g.newVar(utils.Int{})
		}
		g.instructions = append(g.instructions, ir.CondJump{
			BaseInstruction: ir.BaseInstruction{Location: e.GetLocation()},
			Cond:            condVar,
			ThenLabel:       thenLabel,
			ElseLabel:       elseLabel,
		})
		g.instructions = append(g.instructions, thenLabel)
		thenVar := g.visit(st, e.Then)
		res := "unit"
		if e.Else != nil {
			g.instructions = append(g.instructions, ir.Copy{
				BaseInstruction: ir.BaseInstruction{Location: e.GetLocation()},
				Source:          thenVar,
				Dest:            copyVar,
			})
			g.instructions = append(g.instructions, ir.Jump{
				BaseInstruction: ir.BaseInstruction{Location: e.GetLocation()},
				Label:           endLabel,
			})
			g.instructions = append(g.instructions, elseLabel)
			elseVar := g.visit(st, e.Else)
			g.instructions = append(g.instructions, ir.Copy{
				BaseInstruction: ir.BaseInstruction{Location: e.GetLocation()},
				Source:          elseVar,
				Dest:            copyVar,
			})
			res = copyVar
		}
		g.instructions = append(g.instructions, endLabel)
		return res

	case ast.WhileLoop:
		whileStartLabel := g.newLabel()
		g.instructions = append(g.instructions, whileStartLabel)
		condVar := g.visit(st, e.Condition)
		if _, ok := g.varTypes[condVar].(utils.Bool); !ok {
			panic(fmt.Sprintf("Conditional should be boolean %s", g.varTypes[condVar]))
		}
		whileBodyLabel := g.newLabel()
		whileEndLabel := g.newLabel()
		g.instructions = append(g.instructions, ir.CondJump{
			BaseInstruction: ir.BaseInstruction{Location: e.GetLocation()},
			Cond:            condVar,
			ThenLabel:       whileBodyLabel,
			ElseLabel:       whileEndLabel,
		})
		g.instructions = append(g.instructions, whileBodyLabel)
		prevStart := g.loopStartLabel
		prevEnd := g.loopEndLabel
		g.loopStartLabel = &whileStartLabel
		g.loopEndLabel = &whileEndLabel
		g.visit(st, e.Looping)
		g.loopStartLabel = prevStart
		g.loopEndLabel = prevEnd
		g.instructions = append(g.instructions, ir.Jump{
			BaseInstruction: ir.BaseInstruction{Location: e.GetLocation()},
			Label:           whileStartLabel,
		})
		g.instructions = append(g.instructions, whileEndLabel)
		return "unit"

	case ast.BreakExpression:
		if g.loopEndLabel == nil {
			panic("break outside of loop")
		}
		g.instructions = append(g.instructions, ir.Jump{
			BaseInstruction: ir.BaseInstruction{Location: e.GetLocation()},
			Label:           *g.loopEndLabel,
		})
		return "unit"

	case ast.ContinueExpression:
		if g.loopStartLabel == nil {
			panic("continue outside of loop")
		}
		g.instructions = append(g.instructions, ir.Jump{
			BaseInstruction: ir.BaseInstruction{Location: e.GetLocation()},
			Label:           *g.loopStartLabel,
		})
		return "unit"

	case ast.Block:
		innerTable := utils.NewSymTab(st)

		for v := range g.rootTypes {
			innerTable.Table[v] = v
		}

		for _, expr := range e.Expressions {
			g.visit(innerTable, expr)
		}
		res := "unit"
		if e.Result != nil {
			res = g.visit(innerTable, e.Result)
		}
		return res

	case ast.FunctionCall:
		var args []IRVar
		for _, arg := range e.Args {
			args = append(args, g.visit(st, arg))
		}
		dest := g.newVar(utils.Unit{})
		g.instructions = append(g.instructions, ir.Call{
			BaseInstruction: ir.BaseInstruction{Location: e.GetLocation()},
			Fun:             e.Name.(ast.Identifier).Name,
			Args:            args,
			Dest:            dest,
		})
		return dest

	case ast.ReturnExpression:
		val := g.visit(st, e.Result)
		g.instructions = append(g.instructions, ir.Return{
			BaseInstruction: ir.BaseInstruction{Location: e.GetLocation()},
			Value:           val,
		})
		return "unit"

	case ast.Unary:
		var args []IRVar
		args = append(args, g.visit(st, e.Exp))
		dest := g.newVar(g.varTypes[args[0]])
		g.instructions = append(g.instructions, ir.Call{
			BaseInstruction: ir.BaseInstruction{Location: e.GetLocation()},
			Fun:             fmt.Sprintf("unary_%s", e.Op),
			Args:            args,
			Dest:            dest,
		})

		return dest

	default:
		return ""
	}
}
