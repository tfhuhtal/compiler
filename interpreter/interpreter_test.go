package interpreter

import (
	"compiler/parser"
	"compiler/tokenizer"
	// "fmt"
	"testing"
)

func TestInterpreter(t *testing.T) {
	tokens := tokenizer.Tokenize("var x = 0; while x < 100 do { x = x + 1; print_int(x)}", "")
	parsed := parser.Parse(tokens)
	Interpret(parsed)
}
