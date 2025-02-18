package irgenerator

import (
	"compiler/ast"
	"compiler/typechecker"
)

// represent the name of a memory location or built-in
type IRVar struct {
	name string
}

func Generate(rootExpr ast.Expression)
