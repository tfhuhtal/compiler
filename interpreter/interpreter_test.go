package interpreter

import (
	"compiler/parser"
	"compiler/tokenizer"
	"fmt"
	"testing"
)

func helper(input string) any {
	tokens := tokenizer.Tokenize(input, "")
	parsed := parser.Parse(tokens)
	interpreted := Interpret(parsed)
	return interpreted
}

func TestInterpreter_While(t *testing.T) {
	res := helper("var x = 0; while x < 100 do { x = x + 1; print_int(x)}; x ")
	expected := "100"
	if fmt.Sprintf("%v", res) != expected {
		t.Errorf("Expected %v but got %v", expected, res)
	}
}

func TestInterpreter_LangSpec(t *testing.T) {
	res := helper(`
								var n: Int = 100;
								print_int(n);
								while n > 1 do {
								    if n % 2 == 0 then {
								        n = n / 2;
								    } else {
								        n = 3*n + 1;
								    }
								    print_int(n);
								}; print_bool(n == 1)`)
	expected := "true"
	if fmt.Sprintf("%v", res) != expected {
		t.Errorf("Expected %v but got %v", expected, res)
	}
}
