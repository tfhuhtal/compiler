package parser

import (
	"compiler/tokenizer"
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	tokens := tokenizer.Tokenize("var n: Int = read_int();print_int(n);while n > 1 do {if n % 2 == 0 then {n = n / 2;} else {n = 3*n + 1;}print_int(n);}", "")
	Parse(tokens)
}

func TestParser_Declaration(t *testing.T) {
	tokens := tokenizer.Tokenize("var x: Int = 42;", "")
	res := Parse(tokens)
	if res == nil {
		t.Errorf("Expected at least one expression")
	}
}

func TestParser_BinaryOp(t *testing.T) {
	tokens := tokenizer.Tokenize("3 + 4 * 5", "")
	res := Parse(tokens)
	if res == nil {
		t.Errorf("Expected at least one expression")
	}
}

func TestParser_Unary(t *testing.T) {
	tokens := tokenizer.Tokenize("not not false", "")
	res := Parse(tokens)
	expected := "{[] {not {not {false { 1 9}} { 1 9}} { 1 9}} { 1 9}}"
	if fmt.Sprintf("%v", res) != expected {
		t.Errorf("Expected %v but got %v", expected, res)
	}
}

func TestParser_If(t *testing.T) {
	tokens := tokenizer.Tokenize("1 + if 1 < 2 then 10 else 100", "")
	res := Parse(tokens)
	expected := "{[] {{1 { 1 1}} + {{{1 { 1 8}} < {2 { 1 12}} { 1 8}} {10 { 1 19}} {100 { 1 27}} { 1 5}} { 1 1}} { 1 27}}"
	if fmt.Sprintf("%v", res) != expected {
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
		{"if true then var x = 3;", false},
		{"{ { 1 }; 2 { 3 } }", false},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			tokens := tokenizer.Tokenize(tt.code, "")
			if tt.shouldPass {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("Unexpected error for code '%s': %v", tt.code, r)
					}
				}()
				res := Parse(tokens)
				if res == nil {
					t.Errorf("Expected at least one expression")
				}
			} else {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Expected error for code '%s' but got none", tt.code)
					}
				}()
				Parse(tokens)
				t.Errorf("Parsing should have failed for code '%s'", tt.code)
			}
		})
	}
}

func TestParser_While(t *testing.T) {
	tokens := tokenizer.Tokenize("while true do { x = x + 1; }", "")
	res := Parse(tokens)
	expected := "{[] {{true { 1 7}} {[{{x { 1 17}} = {{x { 1 21}} + {1 { 1 25}} { 1 21}} { 1 17}}] <nil> { 1 28}} { 1 1}} { 1 28}}"
	if fmt.Sprintf("%v", res) != fmt.Sprintf("%v", expected) {
		t.Errorf("Expected %v but got %v", expected, res)
	}
}

func TestParser_Block1(t *testing.T) {
	tokens := tokenizer.Tokenize("{123};", "")
	res := Parse(tokens)
	expected := "{[{[] {123 { 1 2}} { 1 5}}] <nil> { 1 6}}"
	if fmt.Sprintf("%v", res) != expected {
		t.Errorf("Expected %v but got %v", expected, res)
	}
}

func TestParser_Block2(t *testing.T) {
	tokens := tokenizer.Tokenize("{123}", "")
	res := Parse(tokens)
	expected := "{[] {[] {123 { 1 2}} { 1 5}} { 1 5}}"

	if fmt.Sprintf("%v", res) != fmt.Sprintf("%v", expected) {
		t.Errorf("Expected %v but got %v", expected, res)
	}
}

func TestParser_Fun1(t *testing.T) {
	tokens := tokenizer.Tokenize(`fun square(x: Int): Int {
								return x * x;
							  }`, "")
	expected := `{[{{square { 1 5}} [{{x { 1 12}} {Int { 1 15}} { 1 12}}] {Int { 1 21}} {[] {[{{{x { 2 16}} * {x { 2 20}} { 2 16}} { 2 9}}] <nil> { 3 10}} { 3 10}} { 1 1}}] {[] <nil> { 3 10}} { 1 1}}`
	result := Parse(tokens)
	if fmt.Sprintf("%v", result) != expected {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}

func TestParser_Fun2(t *testing.T) {
	tokens := tokenizer.Tokenize(`
									fun square(x: Int): Int {
									    return x * x;
									}

									fun vec_len_squared(x: Int, y: Int): Int {
									    return square(x) + square(y);
									}

									fun print_int_twice(x: Int): Unit {
									    print_int(x);
									    print_int(x);
									}

									print_int_twice(vec_len_squared(3, 4));
								`, "")
	expected := `{[{{square { 2 14}} [{{x { 2 21}} {Int { 2 24}} { 2 21}}] {Int { 2 30}} {[{[{{{x { 3 21}} * {x { 3 25}} { 3 21}} { 3 14}}] <nil> { 4 10}}] <nil> { 6 10}} { 2 10}} {{vec_len_squared { 6 14}} [{{x { 6 30}} {Int { 6 33}} { 6 30}} {{y { 6 38}} {Int { 6 41}} { 6 38}}] {Int { 6 47}} {[{[{{{{square { 7 21}} [{x { 7 28}}] { 7 27}} + {{square { 7 33}} [{y { 7 40}}] { 7 39}} { 7 27}} { 7 14}}] <nil> { 8 10}}] <nil> { 10 10}} { 6 10}} {{print_int_twice { 10 14}} [{{x { 10 30}} {Int { 10 33}} { 10 30}}] {Unit { 10 39}} {[{[{{print_int { 11 14}} [{x { 11 24}}] { 11 23}} {{print_int { 12 14}} [{x { 12 24}}] { 12 23}}] <nil> { 13 10}} {{print_int_twice { 15 10}} [{{vec_len_squared { 15 26}} [{3 { 15 42}} {4 { 15 45}}] { 15 41}}] { 15 25}}] <nil> { 15 48}} { 10 10}}] {[] <nil> { 15 48}} { 2 10}}`
	result := Parse(tokens)
	if fmt.Sprintf("%v", result) != expected {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}
