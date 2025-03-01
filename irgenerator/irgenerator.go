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

func Generate(rootTypes map[IRVar]Type, rootExpr ast.Expression) []ir.Instruction {
	varTypes := make(map[IRVar]Type)

	for k, v := range rootTypes {
		varTypes[k] = v
	}

	varUnit := "unit"
	varTypes[varUnit] = utils.Unit{}
	var ins = []ir.Instruction{
		ir.Label{
			BaseInstruction: ir.BaseInstruction{},
			Label:           "start",
		},
	}

	rootSymTab := utils.NewSymTab[IRVar](nil)
	for v := range rootTypes {
		rootSymTab.Table[v] = v
	}

	varFinalResult := visit(rootSymTab, rootExpr, varTypes, &ins)

	if _, ok := varTypes[varFinalResult].(utils.Int); ok {
		ins = append(ins, ir.Call{
			BaseInstruction: ir.BaseInstruction{Location: rootExpr.GetLocation()},
			Fun:             "print_int",
			Args:            []IRVar{varFinalResult},
			Dest:            newVar(utils.Unit{}, varTypes),
		})
	} else if _, ok := varTypes[varFinalResult].(utils.Bool); ok {
		ins = append(ins, ir.Call{
			BaseInstruction: ir.BaseInstruction{Location: rootExpr.GetLocation()},
			Fun:             "print_bool",
			Args:            []IRVar{varFinalResult},
			Dest:            newVar(utils.Unit{}, varTypes),
		})
	}

	return ins
}

func newVar(t Type, varTypes map[IRVar]Type) IRVar {
	idx := 0
	name := fmt.Sprintf("x%d", idx)
	for {
		if _, exists := varTypes[name]; !exists {
			break
		}
		idx++
		name = fmt.Sprintf("x%d", idx)
	}
	varTypes[name] = t
	return name
}

func newLabel(varTypes map[IRVar]Type) IRVar {
	idx := 0
	name := fmt.Sprintf("L%d", idx)
	for {
		if _, exists := varTypes[name]; !exists {
			break
		}
		idx++
		name = fmt.Sprintf("L%d", idx)
	}
	varTypes[name] = utils.Unit{Name: name}
	return name
}

func visit(st *SymTab, expr ast.Expression, varTypes map[IRVar]Type, ins *[]ir.Instruction) IRVar {
	loc := expr.GetLocation()

	switch e := expr.(type) {
	case ast.Literal:
		var variable IRVar
		if value, ok := e.Value.(int); ok {
			variable = newVar(utils.Int{Name: "Int"}, varTypes)
			*ins = append(*ins, ir.LoadIntConst{
				BaseInstruction: ir.BaseInstruction{Location: loc},
				Value:           value,
				Dest:            variable,
			})
		} else {
			panic("Unsupported literal")
		}
		return variable

	case ast.BooleanLiteral:
		var variable IRVar
		if e.Boolean == "true" || e.Boolean == "false" {
			value, _ := strconv.ParseBool(e.Boolean)
			variable = newVar(utils.Bool{Name: "Bool"}, varTypes)
			*ins = append(*ins, ir.LoadBoolConst{
				BaseInstruction: ir.BaseInstruction{Location: loc},
				Value:           value,
				Dest:            variable,
			})
		} else {
			panic("Unsupported boolean literal")
		}
		return variable

	case ast.Identifier:
		if _, exists := st.Table[e.Name]; !exists {
			panic("Perkele")
		}
		return st.Table[e.Name]

	case ast.BinaryOp:
		varOp, exists := st.Table[e.Op]
		if !exists {
			panic("jumankauti")
		}
		left := visit(st, e.Left, varTypes, ins)
		right := visit(st, e.Right, varTypes, ins)
		res := newVar(varTypes[left], varTypes)

		*ins = append(*ins, ir.Call{
			BaseInstruction: ir.BaseInstruction{Location: loc},
			Fun:             varOp,
			Args:            []IRVar{left, right},
			Dest:            res,
		})
		return res

		// TODO: maybe this new should also append to the st
	case ast.Declaration:
		value := visit(st, e.Value, varTypes, ins)
		var name string
		if identifier, ok := e.Variable.(ast.Identifier); ok {
			name = identifier.Name
		}
		if _, exists := st.Table[name]; exists {
			panic(fmt.Sprintf("%s already declared", e.Variable))
		}
		st.Table[name] = value
		new := newVar(e.Value.GetType(), varTypes)
		*ins = append(*ins, ir.Copy{
			BaseInstruction: ir.BaseInstruction{Location: loc},
			Source:          value,
			Dest:            new,
		})
		return st.Table[name]

	case ast.Unary:

	case ast.IfExpression:
		thenLabel := ir.Label{
			BaseInstruction: ir.BaseInstruction{Location: loc},
			Label:           newLabel(varTypes),
		}
		endLabel := ir.Label{
			BaseInstruction: ir.BaseInstruction{Location: loc},
			Label:           newLabel(varTypes),
		}
		var elseLabel ir.Label
		if e.Else != nil {
			elseLabel = ir.Label{
				BaseInstruction: ir.BaseInstruction{Location: loc},
				Label:           newLabel(varTypes),
			}
		} else {
			elseLabel = endLabel
		}
		cond := visit(st, e.Condition, varTypes, ins)
		*ins = append(*ins, ir.CondJump{
			BaseInstruction: ir.BaseInstruction{Location: loc},
			Cond:            cond,
			ThenLabel:       thenLabel,
			ElseLabel:       elseLabel,
		})
		*ins = append(*ins, thenLabel)
		visit(st, e.Then, varTypes, ins)

		/*if e.Else != nil {*/

		/*}*/
		*ins = append(*ins, endLabel)

	case ast.WhileLoop:

	case ast.Function:

	case ast.Block:

	}
	return ""
}
