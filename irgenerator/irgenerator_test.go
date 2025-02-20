package irgenerator

import (
	"compiler/parser"
	"compiler/tokenizer"
	"compiler/utils"
	"fmt"
	"testing"
)

func TestIrGenerator(t *testing.T) {
	t.Run("Testing with simple input", func(t *testing.T) {
		tokens := tokenizer.Tokenize("1 + 2 * 3", "")
		p := parser.New(tokens)
		res := p.Parse()
		var rootTypes = make(map[IRVar]utils.Type)
		rootTypes["+"] = utils.Int{}
		rootTypes["*"] = utils.Int{}
		instructions := Generate(rootTypes, res[0])
		expected := "[Label(start) LoadIntConst(1, x0) LoadIntConst(2, x1) LoadIntConst(3, x2) Call(*, [x1, x2], x3) Call(+, [x0, x3], x4) Call(print_int, [x4], x5)]"
		generatedInstructions := fmt.Sprintf("%v", instructions)

		if generatedInstructions != expected {
			t.Errorf("Expected: %s\nGot: %s", expected, generatedInstructions)
		}
	})
}
