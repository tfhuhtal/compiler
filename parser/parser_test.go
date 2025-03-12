package parser

import (
	"compiler/ast"
	"compiler/tokenizer"
	"compiler/utils"
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	tokens := tokenizer.Tokenize("var n: Int = read_int();print_int(n);while n > 1 do {if n % 2 == 0 then {n = n / 2;} else {n = 3*n + 1;}print_int(n);}", "")
	p := New(tokens)
	p.Parse()
}

func TestParser_Declaration(t *testing.T) {
	tokens := tokenizer.Tokenize("var x: Int = 42;", "")
	p := New(tokens)
	res := p.Parse()
	if res == nil {
		t.Errorf("Expected at least one expression")
	}
}

func TestParser_BinaryOp(t *testing.T) {
	tokens := tokenizer.Tokenize("3 + 4 * 5", "")
	p := New(tokens)
	res := p.Parse()
	if res == nil {
		t.Errorf("Expected at least one expression")
	}
}

func TestParser_Unary(t *testing.T) {
	tokens := tokenizer.Tokenize("not not false", "")
	p := New(tokens)
	res := p.Parse()
	expected := ast.Block{
		Location:    tokenizer.SourceLocation{Line: 1, Column: 1},
		Type:        utils.Unit{},
		Expressions: []ast.Expression{},
		Result: ast.Unary{
			Type:     utils.Unit{},
			Location: tokenizer.SourceLocation{Line: 1, Column: 9},
			Op:       "not",
			Exp: ast.Unary{
				Type:     utils.Unit{},
				Location: tokenizer.SourceLocation{Line: 1, Column: 9},
				Op:       "not",
				Exp: ast.BooleanLiteral{
					Boolean:  "false",
					Location: tokenizer.SourceLocation{Line: 1, Column: 9},
					Type:     utils.Unit{},
				},
			},
		},
	}
	if res.(ast.Block).Result != expected.Result {
		t.Errorf("Expected %v but got %v", expected, res)
	}
}

func TestParser_If(t *testing.T) {
	tokens := tokenizer.Tokenize("1 + if 1 < 2 then 10 else 100", "")
	p := New(tokens)
	res := p.Parse()
	expected := ast.Block{
		Location:    tokenizer.SourceLocation{Line: 1, Column: 1},
		Type:        utils.Unit{},
		Expressions: []ast.Expression{},
		Result: ast.BinaryOp{
			Type:     utils.Unit{},
			Location: tokenizer.SourceLocation{Line: 1, Column: 1},
			Op:       "+",
			Left: ast.Literal{
				Value:    uint64(1),
				Location: tokenizer.SourceLocation{Line: 1, Column: 1},
				Type:     utils.Unit{},
			},
			Right: ast.IfExpression{
				Type:     utils.Unit{},
				Location: tokenizer.SourceLocation{Line: 1, Column: 5},
				Condition: ast.BinaryOp{
					Type:     utils.Unit{},
					Location: tokenizer.SourceLocation{Line: 1, Column: 8},
					Op:       "<",
					Left: ast.Literal{
						Value:    uint64(1),
						Location: tokenizer.SourceLocation{Line: 1, Column: 8},
						Type:     utils.Unit{},
					},
					Right: ast.Literal{
						Value:    uint64(2),
						Location: tokenizer.SourceLocation{Line: 1, Column: 12},
						Type:     utils.Unit{},
					},
				},
				Then: ast.Literal{
					Value:    uint64(10),
					Location: tokenizer.SourceLocation{Line: 1, Column: 19},
					Type:     utils.Unit{},
				},
				Else: ast.Literal{
					Value:    uint64(100),
					Location: tokenizer.SourceLocation{Line: 1, Column: 27},
					Type:     utils.Unit{},
				},
			},
		},
	}
	if res.(ast.Block).Result != expected.Result {
		t.Errorf("Expected %v but got %v", expected, res)
	}
}

func TestParser_Blocks(t *testing.T) {
	tests := []struct {
		code       string
		shouldPass bool
	}{
		{"{ { a } { b } }", true},
		{"{ a b }", false},
		{"{ if true then { a } b }", true},
		{"{ if true then { a }; b }", true},
		{"{ if true then { a } b c }", false},
		{"{ if true then { a } b; c }", true},
		{"{ if true then { a } else { b } c }", true},
		{"x = { { f(a) } { b } }", true},
		{"a + b c", false}, // expect error (garbage at the end)
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			tokens := tokenizer.Tokenize(tt.code, "")
			p := New(tokens)
			if tt.shouldPass {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("Unexpected error for code '%s': %v", tt.code, r)
					}
				}()
				res := p.Parse()
				if res == nil {
					t.Errorf("Expected at least one expression")
				}
			} else {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Expected error for code '%s' but got none", tt.code)
					}
				}()
				p.Parse()
				t.Errorf("Parsing should have failed for code '%s'", tt.code)
			}
		})
	}
}

func TestParser_While(t *testing.T) {
	tokens := tokenizer.Tokenize("while true do { x = x + 1; }", "")
	p := New(tokens)
	res := p.Parse()
	expected := ast.Block{
		Location:    tokenizer.SourceLocation{Line: 1, Column: 1},
		Type:        utils.Unit{},
		Expressions: []ast.Expression{},
		Result: ast.WhileLoop{
			Type:     utils.Unit{},
			Location: tokenizer.SourceLocation{Line: 1, Column: 1},
			Condition: ast.BooleanLiteral{
				Boolean:  "true",
				Location: tokenizer.SourceLocation{Line: 1, Column: 7},
				Type:     utils.Unit{},
			},
			Looping: ast.Block{
				Location: tokenizer.SourceLocation{Line: 1, Column: 17},
				Type:     utils.Unit{},
				Expressions: []ast.Expression{
					ast.BinaryOp{
						Type:     utils.Unit{},
						Location: tokenizer.SourceLocation{Line: 1, Column: 17},
						Op:       "=",
						Left: ast.Identifier{
							Name:     "x",
							Location: tokenizer.SourceLocation{Line: 1, Column: 17},
							Type:     utils.Unit{},
						},
						Right: ast.BinaryOp{
							Type:     utils.Unit{},
							Location: tokenizer.SourceLocation{Line: 1, Column: 21},
							Op:       "+",
							Left: ast.Identifier{
								Name:     "x",
								Location: tokenizer.SourceLocation{Line: 1, Column: 21},
								Type:     utils.Unit{},
							},
							Right: ast.Literal{
								Value:    uint64(1),
								Location: tokenizer.SourceLocation{Line: 1, Column: 25},
								Type:     utils.Unit{},
							},
						},
					},
				},
				Result: ast.Literal{
					Value:    nil,
					Location: tokenizer.SourceLocation{Line: 1, Column: 28},
					Type:     utils.Unit{},
				},
			},
		},
	}

	if fmt.Sprintf("%v", res) != fmt.Sprintf("%v", expected) {
		t.Errorf("Expected %v but got %v", expected, res)
	}
}

/*func TestParser_Block(t *testing.T) {*/
/*tokens := tokenizer.Tokenize("{{{123}}};", "")*/
/*p := New(tokens)*/
/*res := p.Parse()*/
/*expected := ast.Block{*/
/*Type:     utils.Unit{},*/
/*Location: tokenizer.SourceLocation{Line: 1, Column: 1},*/
/*Result:   nil,*/
/*Expressions: []ast.Expression{ast.Block{*/
/*Type:     utils.Unit{},*/
/*Location: tokenizer.SourceLocation{Line: 1, Column: 1},*/
/*Result:   nil,*/
/*Expressions: []ast.Expression{ast.Block{*/
/*Type:        utils.Unit{},*/
/*Location:    tokenizer.SourceLocation{Line: 1, Column: 1},*/
/*Result:      ast.Literal{Type: utils.Int{}, Value: uint64(123)},*/
/*Expressions: []ast.Expression{},*/
/*},*/
/*},*/
/*},*/
/*},*/
/*}*/

/*if fmt.Sprintf("%v", res) != fmt.Sprintf("%v", expected) {*/
/*t.Errorf("Expected %v but got %v", expected, res)*/
/*}*/
/*}*/
