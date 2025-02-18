package typechecker

import (
	"compiler/ast"
	"compiler/utils"
	"compiler/tokenizer"
	"compiler/parser"
	"testing"
)

func TestTypecheck(t *testing.T) {
	t.Run("IntegerLiteral returns Int", func(t *testing.T) {
			literal := ast.Literal{Value: 42}
			symTab := NewSymTab(nil)
			got := typecheck(literal, &symTab)
			if _, ok := got.(utils.Int); !ok {
					t.Errorf("Expected Int type, got %T", got)
			}
	})

	t.Run("BooleanLiteral returns Bool", func(t *testing.T) {
			literal := ast.BooleanLiteral{Boolean: "true"}
			symTab := NewSymTab(nil)
			got := typecheck(literal, &symTab)
			if _, ok := got.(utils.Bool); !ok {
					t.Errorf("Expected Bool type, got %T", got)
			}
	})
	t.Run("BinaryOp with Int operands returns Int", func(t *testing.T) {
			left := ast.Literal{Value: 42}
			right := ast.Literal{Value: 42}
			binaryOp := ast.BinaryOp{Left: left, Op: "+", Right: right}
			symTab := NewSymTab(nil)
			got := typecheck(binaryOp, &symTab)
			if _, ok := got.(utils.Int); !ok {
					t.Errorf("Expected Int type, got %T", got)
			}
	})
	t.Run("BinaryOp with Bool operands returns Bool", func(t *testing.T) {
			left := ast.BooleanLiteral{Boolean: "true"}
			right := ast.BooleanLiteral{Boolean: "true"}
			binaryOp := ast.BinaryOp{Left: left, Op: "==", Right: right}
			symTab := NewSymTab(nil)
			got := typecheck(binaryOp, &symTab)
			if _, ok := got.(utils.Bool); !ok {
					t.Errorf("Expected Bool type, got %T", got)
			}
	})
	t.Run("BinaryOp with Int and Bool operands returns error", func(t *testing.T) {
			left := ast.Literal{Value: 42}
			right := ast.BooleanLiteral{Boolean: "true"}
			binaryOp := ast.BinaryOp{Left: left, Op: "==", Right: right}
			symTab := NewSymTab(nil)
			defer func() {
					if r := recover(); r == nil {
							t.Errorf("Expected panic, got nil")
					}
			}()
			typecheck(binaryOp, &symTab)
	})
	t.Run("Testing more complex input", func(t *testing.T) {
			tokens := tokenizer.Tokenize("var n: Int = 2;print_int(n);while n > 1 do {if n % 2 == 0 then {n = n / 2;} else {n = 3*n + 1;}print_int(n);}", "")
			p := parser.New(tokens)
			res := p.Parse()
			defer func() {
					if r := recover(); r == nil {
							t.Errorf("Expected panic, got nil")
					}
			}()
			Type(res)
	})
}

