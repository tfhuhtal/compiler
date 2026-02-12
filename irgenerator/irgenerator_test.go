package irgenerator

import (
	"compiler/parser"
	"compiler/tokenizer"
	"testing"
)

func TestIr(t *testing.T) {
	t.Run("With block with res unit", func(t *testing.T) {
		tokens := tokenizer.Tokenize("{123};", "")
		parsed := parser.Parse(tokens)
		generated := Generate(parsed)
		if len(generated["main"]) != 1 {
			t.Errorf("there should be only one ir command")
		}
	})
	t.Run("Break generates jump to loop end", func(t *testing.T) {
		tokens := tokenizer.Tokenize("var x: Int = 0; while true do { x = x + 1; if x == 5 then { break } }", "")
		parsed := parser.Parse(tokens)
		generated := Generate(parsed)
		if len(generated["main"]) == 0 {
			t.Errorf("Expected IR instructions for break")
		}
	})
	t.Run("Continue generates jump to loop start", func(t *testing.T) {
		tokens := tokenizer.Tokenize("var x: Int = 0; while x < 10 do { x = x + 1; continue }", "")
		parsed := parser.Parse(tokens)
		generated := Generate(parsed)
		if len(generated["main"]) == 0 {
			t.Errorf("Expected IR instructions for continue")
		}
	})
	t.Run("With block with res as the statements", func(t *testing.T) {
		tokens := tokenizer.Tokenize("{123}", "")
		parsed := parser.Parse(tokens)
		generated := Generate(parsed)
		if len(generated["main"]) != 2 {
			t.Errorf("there should be only two ir commands, %v", generated["main"])
		}
	})
	t.Run("Function definition generates separate instruction list", func(t *testing.T) {
		tokens := tokenizer.Tokenize(`
			fun square(x: Int): Int {
				return x * x;
			}
			square(5)
		`, "")
		parsed := parser.Parse(tokens)
		generated := Generate(parsed)
		if _, ok := generated["square"]; !ok {
			t.Errorf("Expected 'square' function in generated IR")
		}
		if _, ok := generated["main"]; !ok {
			t.Errorf("Expected 'main' in generated IR")
		}
	})
	t.Run("Assign print_int to variable and call", func(t *testing.T) {
		tokens := tokenizer.Tokenize("var x = print_int; x(4)", "")
		parsed := parser.Parse(tokens)
		generated := Generate(parsed)
		if len(generated["main"]) == 0 {
			t.Errorf("Expected IR instructions for function reference call")
		}
	})
	t.Run("Assign print_int to typed variable and call", func(t *testing.T) {
		tokens := tokenizer.Tokenize("var x: (Int) => Unit = print_int; x(4)", "")
		parsed := parser.Parse(tokens)
		generated := Generate(parsed)
		if len(generated["main"]) == 0 {
			t.Errorf("Expected IR instructions for typed function reference call")
		}
	})
	t.Run("Assign print_bool to typed variable and call", func(t *testing.T) {
		tokens := tokenizer.Tokenize("var x: (Bool) => Unit = print_bool; x(true)", "")
		parsed := parser.Parse(tokens)
		generated := Generate(parsed)
		if len(generated["main"]) == 0 {
			t.Errorf("Expected IR instructions for bool function reference call")
		}
	})
	t.Run("Multiple function definitions", func(t *testing.T) {
		tokens := tokenizer.Tokenize(`
			fun add(a: Int, b: Int): Int {
				return a + b;
			}
			fun double(x: Int): Int {
				return add(x, x);
			}
			double(21)
		`, "")
		parsed := parser.Parse(tokens)
		generated := Generate(parsed)
		if len(generated) != 3 {
			t.Errorf("Expected 3 function entries (add, double, main), got %d", len(generated))
		}
	})
}
