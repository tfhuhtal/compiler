package parser

import (
	"compiler/src/tokenizer"
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
	if len(res) == 0 {
		t.Errorf("Expected at least one expression")
	}
}

func TestParser_BinaryOp(t *testing.T) {
	tokens := tokenizer.Tokenize("3 + 4 * 5", "")
	p := New(tokens)
	res := p.Parse()
	if len(res) == 0 {
		t.Errorf("Expected at least one expression")
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
				if len(res) == 0 {
					t.Errorf("Expected non-empty AST for code '%s'", tt.code)
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
