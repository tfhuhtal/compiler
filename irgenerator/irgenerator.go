package irgenerator

import (
	"compiler/ast"
	"compiler/utils"
)

// represent the name of a memory location or built-in
type IRVar struct {
	name string
}

func Generate(rootTypes map[IRVar]utils.Type, rootExpr ast.Expression) {}
