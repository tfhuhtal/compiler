package asmgenerator

import (
	"compiler/ir"
	"fmt"
	"math"
	"strings"
)

type void struct{}

var member void

type Op int

const (
	Add Op = iota
	Sub
	Mul
	Div
	Mod
	Equals
	NotEquals
	GT
	GTE
	LT
	LTE
	Not
	UnarySub
	AddressOf
	Deref
	New
	Delete
)

type Symbol struct {
	op    Op
	value string
}

type Locals struct {
	varToLocation map[ir.IRVar]string
	stackUsed     int
}

func collectAllVars(instructions []ir.Instruction) []ir.IRVar {
	varList := []ir.IRVar{}
	seen := make(map[ir.IRVar]void)
	for _, ins := range instructions {
		for _, v := range ins.GetVars() {
			if v == "" {
				continue
			}
			if _, ok := seen[v]; !ok {
				seen[v] = member
				varList = append(varList, v)
			}
		}
	}
	return varList
}

func GenerateASM(instructions []ir.Instruction) string {
	var lines []string
	emit := func(s string) { lines = append(lines, s) }

	// Initialize local variable tracking
	locs := Locals{
		varToLocation: make(map[ir.IRVar]string),
		stackUsed:     0,
	}

	// Gather all variables and assign them stack locations
	allVars := collectAllVars(instructions)
	for _, v := range allVars {
		locs.stackUsed++
		offset := -8 * locs.stackUsed
		locs.varToLocation[v] = fmt.Sprintf("%d(%%rbp)", offset)
	}

	// Align to 16 bytes if desired:
	if locs.stackUsed%2 != 0 {
		locs.stackUsed++
	}
	stackFrameSize := locs.stackUsed * 8

	// Emit a minimal function prologue
	emit(".extern print_int\n.extern print_bool\n.extern read_int\n.section .text\n\n")
	emit(".global main")
	emit(".type main, @function")
	emit("main:")
	for k, v := range locs.varToLocation {
		emit(fmt.Sprintf("# %s in %s", k, v))
	}
	emit("    pushq %rbp")
	emit("    movq %rsp, %rbp")
	emit(fmt.Sprintf("    subq $%d, %%rsp\n", stackFrameSize))

	for _, ins := range instructions {
		switch i := ins.(type) {

		case ir.LoadBoolConst:
			emit(fmt.Sprintf("# %s", i.String()))
			loc := locs.varToLocation[i.Dest]
			val := 0
			if i.Value {
				val = 1
			}
			emit(fmt.Sprintf("movq $%d, %s", val, loc))

		case ir.LoadIntConst:
			emit(fmt.Sprintf("# %s", i.String()))
			loc := locs.varToLocation[i.Dest]
			if i.Value > math.MaxUint32 {
				emit(fmt.Sprintf("movabsq $%d, %%rax", i.Value))
				emit(fmt.Sprintf("movq %%rax, %s\n", loc))
			} else {
				emit(fmt.Sprintf("movq $%d, %s\n", i.Value, loc))
			}

		case ir.Label:
			if i.String() != "" {
				emit(fmt.Sprintf(".%s:\n", i.Label))
			}

		case ir.Call:
			emit(fmt.Sprintf("# %s", i.String()))
			if i.Fun == "print_int" || i.Fun == "read_int" {
				emit("subq $8, %rsp")
				lines = append(lines, generateCall(i.Fun, i.Args, &locs)...)
				emit(mov("%rax", locs.varToLocation[i.Dest]))
				emit("add $8, %rsp\n")

			} else {
				lines = append(lines, generateCall(i.Fun, i.Args, &locs)...)
				emit(mov("%rax", locs.varToLocation[i.Dest]))
				emit("\n")

			}

		case ir.Copy:
			emit(fmt.Sprintf("# %s", i.String()))
			emit(mov(locs.varToLocation[i.Source], "%rax"))
			emit(mov("%rax", locs.varToLocation[i.Dest]))
			emit("\n")

		case ir.CondJump:
			emit(fmt.Sprintf("# %s", i.String()))
			emit(fmt.Sprintf("cmpq $0, %s", locs.varToLocation[i.Cond]))
			emit(fmt.Sprintf("jne .%s", i.ThenLabel.Label))
			emit(fmt.Sprintf("jmp .%s\n", i.ElseLabel.Label))

		case ir.Jump:
			emit(fmt.Sprintf("# %s", i.String()))
			emit(fmt.Sprintf("jmp .%s\n", i.Label.Label))

		default:
			emit(fmt.Sprintf("# Unhandled instruction: %v\n", i))
		}
	}

	// Emit a minimal function epilogue
	emit("movq $0, %rax")
	emit("movq %rbp, %rsp")
	emit("popq %rbp")
	emit("ret")

	// Optionally append standard library stubs if needed
	return strings.Join(lines, "\n")
}

func generateCall(fun ir.IRVar, args []ir.IRVar, locs *Locals) []string {
	calleeSym, ok := operatorFromStr(fun, len(args))
	var callee Symbol
	if ok {
		callee = calleeSym
	} else {
		callee = Symbol{value: fun}
	}

	fmt.Println(callee.value, callee.op)

	switch {
	case callee.value == "":
		if len(args) == 2 {
			arg1Loc := locs.varToLocation[args[0]]
			arg2Loc := locs.varToLocation[args[1]]
			lines := []string{}

			switch callee.op {
			case Add:
				lines = append(lines, binOp(&arg1Loc, &arg2Loc, "addq")...)
			case Sub:
				lines = append(lines, binOp(&arg1Loc, &arg2Loc, "subq")...)
			case Mul:
				lines = append(lines, binOp(&arg1Loc, &arg2Loc, "imulq")...)
			case Div:
				lines = append(lines, mov(arg1Loc, "%rax"), "cqto",
					fmt.Sprintf("idivq %s", arg2Loc))
			case Mod:
				lines = append(lines,
					mov(arg1Loc, "%rax"),
					"cqto",
					fmt.Sprintf("idivq %s", arg2Loc),
					mov("%rdx", "%rax"),
				)
			case Equals:
				lines = append(lines, comparison(&arg1Loc, &arg2Loc, "sete")...)
			case NotEquals:
				lines = append(lines, comparison(&arg1Loc, &arg2Loc, "setne")...)
			case GT:
				lines = append(lines, comparison(&arg1Loc, &arg2Loc, "setg")...)
			case GTE:
				lines = append(lines, comparison(&arg1Loc, &arg2Loc, "setge")...)
			case LT:
				lines = append(lines, comparison(&arg1Loc, &arg2Loc, "setl")...)
			case LTE:
				lines = append(lines, comparison(&arg1Loc, &arg2Loc, "setle")...)
			default:
				panic(fmt.Sprintf("operator %d does not have an intrinsic definition", callee.op))
			}
			return lines
		}

		lines := []string{}
		arg1Loc := locs.varToLocation[args[0]]

		switch callee.op {
		case Not:
			lines = append(lines,
				mov(arg1Loc, "%rax"),
				"xorq $0x1, %rax",
			)
		case UnarySub:
			lines = append(lines,
				mov(arg1Loc, "%rax"),
				"negq %rax",
			)
		default:
			lines = append(lines, fmt.Sprintf("; todo operator %d", callee.op))
		}
		return lines

	default:
		switch callee.value {
		default:
			return generateFunctionCall(fun, args, locs)
		}
	}
}

func generateFunctionCall(fun ir.IRVar, args []ir.IRVar, locs *Locals) []string {
	lines := []string{}

	if args[0] != "" {
		lines = append(lines, mov(locs.varToLocation[args[0]], "%rdi"))
	}
	lines = append(lines, fmt.Sprintf("callq %s", fun))

	return lines
}

func operatorFromStr(op string, argCount int) (Symbol, bool) {
	if argCount == 2 {
		switch op {
		case "+":
			return Symbol{op: Add}, true
		case "-":
			return Symbol{op: Sub}, true
		case "*":
			return Symbol{op: Mul}, true
		case "/":
			return Symbol{op: Div}, true
		case "==":
			return Symbol{op: Equals}, true
		case "!=":
			return Symbol{op: NotEquals}, true
		case ">":
			return Symbol{op: GT}, true
		case ">=":
			return Symbol{op: GTE}, true
		case "<":
			return Symbol{op: LT}, true
		case "<=":
			return Symbol{op: LTE}, true
		}
	} else if argCount == 1 {
		switch op {
		case "unary_-":
			return Symbol{op: UnarySub}, true
		case "not":
			return Symbol{op: Not}, true
		}
	}
	return Symbol{}, false
}

func comparison(a *string, b *string, setInstr string) []string {
	return []string{
		fmt.Sprintf("xor %%rax, %%rax"),
		mov(*a, "%rdx"),
		fmt.Sprintf("cmpq %s, %%rdx", *b),
		fmt.Sprintf("%s %%al", setInstr),
	}
}

func binOp(a *string, b *string, op string) []string {
	return []string{
		mov(*a, "%rax"),
		fmt.Sprintf("%s %s, %s", op, *b, "%rax")}
}

func mov(src string, dst string) string {
	return fmt.Sprintf("movq %s, %s", src, dst)
}
