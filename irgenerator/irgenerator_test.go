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
	t.Run("With block with res as the statements", func(t *testing.T) {
		tokens := tokenizer.Tokenize("{123}", "")
		parsed := parser.Parse(tokens)
		generated := Generate(parsed)
		if len(generated) != 2 {
			t.Errorf("there should be only two ir commands, %v", generated)
		}
	})
}
