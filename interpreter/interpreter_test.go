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

func TestInterpreter_Break(t *testing.T) {
	res := helper("var x: Int = 0; while true do { x = x + 1; if x == 5 then { break } }; x")
	expected := "5"
	if fmt.Sprintf("%v", res) != expected {
		t.Errorf("Expected %v but got %v", expected, res)
	}
}

func TestInterpreter_Continue(t *testing.T) {
	res := helper("var x: Int = 0; var y: Int = 0; while x < 10 do { x = x + 1; if x % 2 == 0 then { continue }; y = y + 1 }; y")
	expected := "5"
	if fmt.Sprintf("%v", res) != expected {
		t.Errorf("Expected %v but got %v", expected, res)
	}
}

func TestInterpreter_BreakNested(t *testing.T) {
	res := helper("var x: Int = 0; var y: Int = 0; while x < 10 do { x = x + 1; while true do { y = y + 1; break } }; y")
	expected := "10"
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

func TestInterpreter_SimpleFunction(t *testing.T) {
	res := helper(`
		fun square(x: Int): Int {
			return x * x;
		}
		square(5)
	`)
	expected := "25"
	if fmt.Sprintf("%v", res) != expected {
		t.Errorf("Expected %v but got %v", expected, res)
	}
}

func TestInterpreter_MultipleFunctions_1(t *testing.T) {
	res := helper(`
		fun square(x: Int): Int {
			return x * x;
		}
		fun vec_len_squared(x: Int, y: Int): Int {
			return square(x) + square(y);
		}
		vec_len_squared(3, 4);
	`)
	expected := "<nil>"
	if fmt.Sprintf("%v", res) != expected {
		t.Errorf("Expected %v but got %v", expected, res)
	}
}

func TestInterpreter_MultipleFunctions_2(t *testing.T) {
	res := helper(`
		fun square(x: Int): Int {
			return x * x;
		}
		fun vec_len_squared(x: Int, y: Int): Int {
			return square(x) + square(y);
		}
		vec_len_squared(3, 4)
	`)
	expected := "25"
	if fmt.Sprintf("%v", res) != expected {
		t.Errorf("Expected %v but got %v", expected, res)
	}
}

func TestInterpreter_Recursion(t *testing.T) {
	res := helper(`
		fun factorial(n: Int): Int {
			if n <= 1 then {
				return 1;
			} else {
				return n * factorial(n - 1);
			}
		}
		factorial(5)
	`)
	expected := "120"
	if fmt.Sprintf("%v", res) != expected {
		t.Errorf("Expected %v but got %v", expected, res)
	}
}

func TestInterpreter_UnitFunction(t *testing.T) {
	res := helper(`
		fun do_print(x: Int): Unit {
			print_int(x);
		}
		do_print(42)
	`)
	// do_print returns nil (Unit)
	if res != nil {
		t.Errorf("Expected nil for Unit function, got %v", res)
	}
}

func TestInterpreter_MutualRecursion(t *testing.T) {
	res := helper(`
		fun is_even(n: Int): Bool {
			if n == 0 then {
				return true;
			} else {
				return is_odd(n - 1);
			}
		}
		fun is_odd(n: Int): Bool {
			if n == 0 then {
				return false;
			} else {
				return is_even(n - 1);
			}
		}
		is_even(10)
	`)
	expected := "true"
	if fmt.Sprintf("%v", res) != expected {
		t.Errorf("Expected %v but got %v", expected, res)
	}
}
