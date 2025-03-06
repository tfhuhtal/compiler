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
	varTypes     map[IRVar]Type
	rootTypes    map[IRVar]Type
	instructions []ir.Instruction
}

func NewIRGenerator(rootTypes map[IRVar]Type) *IRGenerator {
	gen := &IRGenerator{
		varTypes:     make(map[IRVar]Type),
		rootTypes:    rootTypes,
		instructions: []ir.Instruction{},
	}
	for k, v := range rootTypes {
		gen.varTypes[k] = v
	}
	gen.varTypes["unit"] = utils.Unit{}
	gen.instructions = append(gen.instructions, gen.newLabel())
	return gen
}

func (g *IRGenerator) Generate(rootExpr ast.Expression) []ir.Instruction {
	rootSymTab := utils.NewSymTab[IRVar](nil)
	for v := range g.varTypes {
		rootSymTab.Table[v] = v
	}
	result := g.visit(rootSymTab, rootExpr)

	fmt.Println(g.varTypes)

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
	return g.instructions
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
		if value, ok := e.Value.(int); ok {
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
		fmt.Println(st.Table)
		if _, exists := st.Table[e.Name]; !exists {
			panic(fmt.Sprintf("Undefined variable: %s, in location %v", e.Name, e.GetLocation()))
		}
		return st.Table[e.Name]

	case ast.BinaryOp:
		left := g.visit(st, e.Left)
		if e.Op == "=" {
			res := g.newVar(g.varTypes[left])
			g.instructions = append(g.instructions, ir.Copy{
				BaseInstruction: ir.BaseInstruction{Location: e.GetLocation()},
				Source:          res,
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
			panic(fmt.Sprintf("%s already declared", e.Variable))
		}
		newVar := g.newVar(e.Value.GetType())
		st.Table[name] = newVar
		g.instructions = append(g.instructions, ir.Copy{
			BaseInstruction: ir.BaseInstruction{Location: e.GetLocation()},
			Source:          value,
			Dest:            newVar,
		})
		return newVar

	case ast.IfExpression:
		thenLabel := g.newLabel()
		endLabel := g.newLabel()
		elseLabel := endLabel
		if e.Else != nil {
			elseLabel = g.newLabel()
		}
		condVar := g.visit(st, e.Condition)
		copyVar := g.newVar(utils.Int{Name: "copy"})
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
		whileBodyLabel := g.newLabel()
		whileEndLabel := g.newLabel()
		g.instructions = append(g.instructions, ir.CondJump{
			BaseInstruction: ir.BaseInstruction{Location: e.GetLocation()},
			Cond:            condVar,
			ThenLabel:       whileBodyLabel,
			ElseLabel:       whileEndLabel,
		})
		g.instructions = append(g.instructions, whileBodyLabel)
		g.visit(st, e.Looping)
		g.instructions = append(g.instructions, ir.Jump{
			BaseInstruction: ir.BaseInstruction{Location: e.GetLocation()},
			Label:           whileStartLabel,
		})
		g.instructions = append(g.instructions, whileEndLabel)
		return "unit"

	case ast.Block:
		innerTable := utils.NewSymTab(st)

		for v := range g.rootTypes {
			innerTable.Table[v] = v
		}

		for _, expr := range e.Expressions {
			g.visit(st, expr)
		}
		res := "unit"
		if e.Result != nil {
			res = g.visit(st, e.Result)
		}
		/*g.instructions = append(g.instructions, ir.Copy{*/
		/*BaseInstruction: ir.BaseInstruction{Location: loc},*/
		/*Source:          exprs,*/
		/*Dest:            res,*/
		/*})*/
		return res

	case ast.Function:
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

	case ast.Unary:
		var ops []IRVar
		for _, op := range e.Ops {
			ops = append(ops, op)
		}

		opFun, found := st.Table[e.Exp.(ast.Identifier).Name]
		if !found {
			panic(fmt.Sprintf("Unknown unary operator: %s", e.Exp))
		}
		resVar := g.newVar(e.GetType())
		g.instructions = append(g.instructions, ir.Call{
			BaseInstruction: ir.BaseInstruction{Location: e.GetLocation()},
			Fun:             opFun,
			Args:            ops,
			Dest:            resVar,
		})
		return resVar

	default:
		return ""
	}
}
