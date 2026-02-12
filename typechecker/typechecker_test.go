package typechecker

import (
	"compiler/ast"
	"compiler/parser"
	"compiler/tokenizer"
	"compiler/utils"
	"testing"
)

func TestTypecheck(t *testing.T) {
	t.Run("IntegerLiteral returns Int", func(t *testing.T) {
		literal := ast.Literal{Value: uint64(42)}
		symTab := utils.NewSymTab[utils.Type](nil)
		got := typecheck(literal, symTab)
		if _, ok := got.(utils.Int); !ok {
			t.Errorf("Expected Int type, got %T", got)
		}
	})

	t.Run("BooleanLiteral returns Bool", func(t *testing.T) {
		literal := ast.BooleanLiteral{Boolean: "true"}
		symTab := utils.NewSymTab[utils.Type](nil)
		got := typecheck(literal, symTab)
		if _, ok := got.(utils.Bool); !ok {
			t.Errorf("Expected Bool type, got %T", got)
		}
	})
	t.Run("BinaryOp with Int operands returns Int", func(t *testing.T) {
		left := ast.Literal{Value: uint64(42)}
		right := ast.Literal{Value: uint64(42)}
		binaryOp := ast.BinaryOp{Left: left, Op: "+", Right: right}
		symTab := utils.NewSymTab[utils.Type](nil)
		got := typecheck(binaryOp, symTab)
		if _, ok := got.(utils.Int); !ok {
			t.Errorf("Expected Int type, got %T", got)
		}
	})
	t.Run("BinaryOp with Bool operands returns Bool", func(t *testing.T) {
		left := ast.BooleanLiteral{Boolean: "true"}
		right := ast.BooleanLiteral{Boolean: "true"}
		binaryOp := ast.BinaryOp{Left: left, Op: "==", Right: right}
		symTab := utils.NewSymTab[utils.Type](nil)
		got := typecheck(binaryOp, symTab)
		if _, ok := got.(utils.Bool); !ok {
			t.Errorf("Expected Bool type, got %T", got)
		}
	})
	t.Run("BinaryOp with Int and Bool operands returns error", func(t *testing.T) {
		left := ast.Literal{Value: uint64(42)}
		right := ast.BooleanLiteral{Boolean: "true"}
		binaryOp := ast.BinaryOp{Left: left, Op: "==", Right: right}
		symTab := utils.NewSymTab[utils.Type](nil)
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic, got nil")
			}
		}()
		typecheck(binaryOp, symTab)
	})
	t.Run("Testing more complex input", func(t *testing.T) {
		tokens := tokenizer.Tokenize("var n: Int = 2;print_int(n);while n > 1 do {if n % 2 == 0 then {n = n / 2;} else {n = 3*n + 1;}print_int(n);}", "")
		res := parser.Parse(tokens)
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Expected panic, got nil")
			}
		}()
		Type(res)
	})
	t.Run("Allowed unary", func(t *testing.T) {
		tokens := tokenizer.Tokenize("not (1*2)", "")
		res := parser.Parse(tokens)
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic, got nil")
			}
		}()
		Type(res)
	})
	t.Run("Not allowed not", func(t *testing.T) {
		tokens := tokenizer.Tokenize("-(1*2)", "")
		res := parser.Parse(tokens)
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Expected panic, got nil")
			}
		}()
		Type(res)
	})

	t.Run("Break in while returns Unit", func(t *testing.T) {
		tokens := tokenizer.Tokenize("while true do { break }", "")
		res := parser.Parse(tokens)
		got := Type(res)
		if _, ok := got.(utils.Unit); !ok {
			t.Errorf("Expected Unit type, got %T", got)
		}
	})
	t.Run("Continue in while returns Unit", func(t *testing.T) {
		tokens := tokenizer.Tokenize("var x: Int = 0; while x < 10 do { x = x + 1; continue }", "")
		res := parser.Parse(tokens)
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Unexpected panic: %v", r)
			}
		}()
		Type(res)
	})

	t.Run("Function definition type checks", func(t *testing.T) {
		tokens := tokenizer.Tokenize(`
			fun square(x: Int): Int {
				return x * x;
			}
			square(5)
		`, "")
		res := parser.Parse(tokens)
		got := Type(res)
		if _, ok := got.(utils.Int); !ok {
			t.Errorf("Expected Int type, got %T", got)
		}
	})

	t.Run("Function call argument type mismatch panics", func(t *testing.T) {
		tokens := tokenizer.Tokenize(`
			fun square(x: Int): Int {
				return x * x;
			}
			square(true)
		`, "")
		res := parser.Parse(tokens)
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for type mismatch")
			}
		}()
		Type(res)
	})

	t.Run("Mutual recursion type checks", func(t *testing.T) {
		tokens := tokenizer.Tokenize(`
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
		`, "")
		res := parser.Parse(tokens)
		got := Type(res)
		if _, ok := got.(utils.Bool); !ok {
			t.Errorf("Expected Bool type, got %T", got)
		}
	})

	t.Run("Multiple params function type checks", func(t *testing.T) {
		tokens := tokenizer.Tokenize(`
			fun add(a: Int, b: Int): Int {
				return a + b;
			}
			add(3, 4)
		`, "")
		res := parser.Parse(tokens)
		got := Type(res)
		if _, ok := got.(utils.Int); !ok {
			t.Errorf("Expected Int type, got %T", got)
		}
	})

	t.Run("Assign print_int to variable and call it", func(t *testing.T) {
		tokens := tokenizer.Tokenize("var x = print_int; x(4)", "")
		res := parser.Parse(tokens)
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Unexpected panic: %v", r)
			}
		}()
		Type(res)
	})

	t.Run("Assign print_int with function type annotation", func(t *testing.T) {
		tokens := tokenizer.Tokenize("var x: (Int) => Unit = print_int; x(4)", "")
		res := parser.Parse(tokens)
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Unexpected panic: %v", r)
			}
		}()
		Type(res)
	})

	t.Run("Assign print_bool with function type annotation", func(t *testing.T) {
		tokens := tokenizer.Tokenize("var x: (Bool) => Unit = print_bool; x(true)", "")
		res := parser.Parse(tokens)
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Unexpected panic: %v", r)
			}
		}()
		Type(res)
	})

	t.Run("Return type mismatch should fail", func(t *testing.T) {
		tokens := tokenizer.Tokenize(`
			fun f(x: Int): Bool { return x + x; }
			f(3)
		`, "")
		res := parser.Parse(tokens)
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for return type mismatch")
			}
		}()
		Type(res)
	})

	t.Run("Duplicate parameter names should fail", func(t *testing.T) {
		tokens := tokenizer.Tokenize(`
			fun f(x: Int, x: Int): Int { return x+1; }
			f(3, 4)
		`, "")
		res := parser.Parse(tokens)
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for duplicate parameter names")
			}
		}()
		Type(res)
	})

	t.Run("Simple function with return", func(t *testing.T) {
		tokens := tokenizer.Tokenize(`
			fun f(x: Int): Int { return x+1; }
			f(3)
		`, "")
		res := parser.Parse(tokens)
		got := Type(res)
		if _, ok := got.(utils.Int); !ok {
			t.Errorf("Expected Int type, got %T", got)
		}
	})

	t.Run("Function with early return in if", func(t *testing.T) {
		tokens := tokenizer.Tokenize(`
			fun f(n: Int): Int {
				if n > 5 then {
					return 5;
				}
				123
			}
			f(10)
		`, "")
		res := parser.Parse(tokens)
		got := Type(res)
		if _, ok := got.(utils.Int); !ok {
			t.Errorf("Expected Int type, got %T", got)
		}
	})

	t.Run("Function with early return and final return", func(t *testing.T) {
		tokens := tokenizer.Tokenize(`
			fun f(n: Int): Int {
				if n > 5 then {
					return 5;
				}
				return 123;
			}
			f(10)
		`, "")
		res := parser.Parse(tokens)
		got := Type(res)
		if _, ok := got.(utils.Int); !ok {
			t.Errorf("Expected Int type, got %T", got)
		}
	})

	t.Run("Multi param function with return", func(t *testing.T) {
		tokens := tokenizer.Tokenize(`
			fun f(p1: Int, p2: Int): Int {
				return p1 + p2;
			}
			f(1, 20)
		`, "")
		res := parser.Parse(tokens)
		got := Type(res)
		if _, ok := got.(utils.Int); !ok {
			t.Errorf("Expected Int type, got %T", got)
		}
	})

	t.Run("Mutual recursion with return", func(t *testing.T) {
		tokens := tokenizer.Tokenize(`
			fun f(x: Int): Int {
				return g(x + 1);
			}
			fun g(x: Int): Int {
				var result = 0;
				if x % 7 != 0 then {
					result = f(x);
				} else {
					result = x;
				}
				return result;
			}
			f(9)
		`, "")
		res := parser.Parse(tokens)
		got := Type(res)
		if _, ok := got.(utils.Int); !ok {
			t.Errorf("Expected Int type, got %T", got)
		}
	})

	t.Run("Factorial with var and return", func(t *testing.T) {
		tokens := tokenizer.Tokenize(`
			fun factorial(x: Int): Int {
				var result = 0;
				if x > 1 then {
					result = x * factorial(x - 1);
				} else {
					result = 1;
				}
				return result;
			}
			factorial(5)
		`, "")
		res := parser.Parse(tokens)
		got := Type(res)
		if _, ok := got.(utils.Int); !ok {
			t.Errorf("Expected Int type, got %T", got)
		}
	})

	t.Run("Factorial with return in branches", func(t *testing.T) {
		tokens := tokenizer.Tokenize(`
			fun factorial(x: Int): Int {
				if x > 1 then {
					return x * factorial(x - 1);
				} else {
					return 1;
				}
			}
			factorial(5)
		`, "")
		res := parser.Parse(tokens)
		got := Type(res)
		if _, ok := got.(utils.Int); !ok {
			t.Errorf("Expected Int type, got %T", got)
		}
	})

	t.Run("Return from inside while loop", func(t *testing.T) {
		tokens := tokenizer.Tokenize(`
			fun f(): Int {
				var i = 1;
				var sum = 0;
				while true do {
					sum = sum + i;
					i = i + 1;
					if i > 5 then {
						return sum;
					}
				}
			}
			f()
		`, "")
		res := parser.Parse(tokens)
		got := Type(res)
		if _, ok := got.(utils.Int); !ok {
			t.Errorf("Expected Int type, got %T", got)
		}
	})

	t.Run("Square function with return", func(t *testing.T) {
		tokens := tokenizer.Tokenize(`
			fun square(x: Int): Int {
				return x * x;
			}
			square(3)
		`, "")
		res := parser.Parse(tokens)
		got := Type(res)
		if _, ok := got.(utils.Int); !ok {
			t.Errorf("Expected Int type, got %T", got)
		}
	})

}
