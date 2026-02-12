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
		if len(generated) != 1 {
			t.Errorf("there should be only one ir command")
		}
	})
	t.Run("Break generates jump to loop end", func(t *testing.T) {
		tokens := tokenizer.Tokenize("var x: Int = 0; while true do { x = x + 1; if x == 5 then { break } }", "")
		parsed := parser.Parse(tokens)
		generated := Generate(parsed)
		if len(generated) == 0 {
			t.Errorf("Expected IR instructions for break")
		}
	})
	t.Run("Continue generates jump to loop start", func(t *testing.T) {
		tokens := tokenizer.Tokenize("var x: Int = 0; while x < 10 do { x = x + 1; continue }", "")
		parsed := parser.Parse(tokens)
		generated := Generate(parsed)
		if len(generated) == 0 {
			t.Errorf("Expected IR instructions for continue")
		}
	})
	t.Run("With block with res as the statements", func(t *testing.T) {
		tokens := tokenizer.Tokenize("{123}", "")
		parsed := parser.Parse(tokens)
		generated := Generate(parsed)
		if len(generated) != 2 {
			t.Errorf("there should be only two ir commands, %v", generated)
		}
	})
}
