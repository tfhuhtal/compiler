package ir

import (
	"compiler/tokenizer"
	"compiler/utils"
	"fmt"
	"strings"
)

// Location is an alias for tokenizer.SourceLocation
type Location = tokenizer.SourceLocation
type Type = utils.Type

// Instruction defines an IR instruction.
type Instruction interface {
	isInstruction()
	String() string
}

// BaseInstruction is embedded by all concrete instruction types.
type BaseInstruction struct {
	Location
}

func (BaseInstruction) isInstruction() {}

// IRVar represents an IR variable.
type IRVar = string

// LoadBoolConst represents loading a boolean constant into Dest.
type LoadBoolConst struct {
	BaseInstruction
	Value bool
	Dest  IRVar
}

func (l LoadBoolConst) String() string {
	return fmt.Sprintf("LoadBoolConst(%v, %v)", l.Value, l.Dest)
}

// LoadIntConst represents loading an integer constant into Dest.
type LoadIntConst struct {
	BaseInstruction
	Value int
	Dest  IRVar
}

func (l LoadIntConst) String() string {
	return fmt.Sprintf("LoadIntConst(%d, %v)", l.Value, l.Dest)
}

// Copy represents copying a value from Source to Dest.
type Copy struct {
	BaseInstruction
	Source IRVar
	Dest   IRVar
}

func (c Copy) String() string {
	return fmt.Sprintf("Copy(%v, %v)", c.Source, c.Dest)
}

// Call represents calling a function or built-in.
type Call struct {
	BaseInstruction
	Fun  IRVar
	Args []IRVar
	Dest IRVar
}

func (c Call) String() string {
	args := make([]string, len(c.Args))
	for i, arg := range c.Args {
		args[i] = arg
	}
	return fmt.Sprintf("Call(%v, [%s], %v)", c.Fun, strings.Join(args, ", "), c.Dest)
}

// Jump represents an unconditional jump.
type Jump struct {
	BaseInstruction
	Label Label
}

func (j Jump) String() string {
	return fmt.Sprintf("Jump(Label: %v)", j.Label)
}

// CondJump represents a conditional jump.
type CondJump struct {
	BaseInstruction
	Cond      IRVar
	ThenLabel Label
	ElseLabel Label
}

func (c CondJump) String() string {
	return fmt.Sprintf("CondJump(%v, %v, %v)", c.Cond, c.ThenLabel, c.ElseLabel)
}

// Label is used for jump targets.
type Label struct {
	BaseInstruction
	Label string
}

func (l Label) String() string {
	return fmt.Sprintf("Label(%s)", l.Label)
}

type PrintInt struct {
	BaseInstruction
	Value IRVar
}

func (p PrintInt) String() string {
	return fmt.Sprintf("Call(print_int, %s)", p.Value)
}

type PrintBool struct {
	BaseInstruction
	Value IRVar
}

func (p PrintBool) String() string {
	return fmt.Sprintf("Call(print_bool, %s)", p.Value)
}
