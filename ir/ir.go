package ir

import (
	"compiler/tokenizer"
	"compiler/utils"
	"fmt"
	"strings"
)

type Location = tokenizer.SourceLocation
type Type = utils.Type

type Instruction interface {
	isInstruction()
	String() string
	GetVars() []IRVar
}

type BaseInstruction struct {
	Location
}

func (BaseInstruction) isInstruction() {}

type IRVar = string

type LoadBoolConst struct {
	BaseInstruction
	Value bool
	Dest  IRVar
}

func (l LoadBoolConst) String() string {
	return fmt.Sprintf("LoadBoolConst(%v, %v)", l.Value, l.Dest)
}

func (l LoadBoolConst) GetVars() []IRVar {
	return []IRVar{l.Dest}
}

type LoadIntConst struct {
	BaseInstruction
	Value uint64
	Dest  IRVar
}

func (l LoadIntConst) String() string {
	return fmt.Sprintf("LoadIntConst(%d, %v)", l.Value, l.Dest)
}

func (l LoadIntConst) GetVars() []IRVar {
	return []IRVar{l.Dest}
}

type Copy struct {
	BaseInstruction
	Source IRVar
	Dest   IRVar
}

func (c Copy) String() string {
	return fmt.Sprintf("Copy(%v, %v)", c.Source, c.Dest)
}

func (c Copy) GetVars() []IRVar {
	return []IRVar{c.Source, c.Dest}
}

type Call struct {
	BaseInstruction
	Fun  IRVar
	Args []IRVar
	Dest IRVar
}

func (c Call) String() string {
	args := make([]string, len(c.Args))
	copy(args, c.Args)
	return fmt.Sprintf("Call(%v, [%s], %v)", c.Fun, strings.Join(args, ", "), c.Dest)
}

func (c Call) GetVars() []IRVar {
	return append([]IRVar{c.Dest}, c.Args...)
}

type Jump struct {
	BaseInstruction
	Label Label
}

func (j Jump) String() string {
	return fmt.Sprintf("Jump(%v)", j.Label)
}

func (j Jump) GetVars() []IRVar {
	return nil
}

type CondJump struct {
	BaseInstruction
	Cond      IRVar
	ThenLabel Label
	ElseLabel Label
}

func (c CondJump) String() string {
	return fmt.Sprintf("CondJump(%v, %v, %v)", c.Cond, c.ThenLabel, c.ElseLabel)
}

func (c CondJump) GetVars() []IRVar {
	return []IRVar{c.Cond}
}

type Label struct {
	BaseInstruction
	Label string
}

func (l Label) String() string {
	return fmt.Sprintf("Label(%s)", l.Label)
}

func (l Label) GetVars() []IRVar {
	return nil
}
