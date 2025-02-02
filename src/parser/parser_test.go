package parser

import (
	"testing"
	"compiler/src/ast"
	"compiler/src/tokenizer"
)

func TestParser(t *testing.T) {
	tests := []struct {
		input    string
		expected ast.Expression
		hasError bool
	}{
		{
			input: "(a + b) * c",
			expected: ast.BinaryOp{
				Left: ast.BinaryOp{
					Left:  ast.Identifier{Name: "a"},
					Op:    "+",
					Right: ast.Identifier{Name: "b"},
				},
				Op:    "*",
				Right: ast.Identifier{Name: "c"},
			},
			hasError: false,
		},
		{
			input: "f(x, y + z)",
			expected: ast.FunctionCall{
				Name: ast.Identifier{Name: "f"},
				Args: []ast.Expression{
					ast.Identifier{Name: "x"},
					ast.BinaryOp{
						Left:  ast.Identifier{Name: "y"},
						Op:    "+",
						Right: ast.Identifier{Name: "z"},
					},
				},
			},
			hasError: false,
		},
		{
			input:    "a + b c",
			expected: nil,
			hasError: true,
		},
		{
			input:    "",
			expected: nil,
			hasError: true,
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.input, "test")
		result, err := Parse(tokens)
		if test.hasError {
			if err == nil {
				t.Errorf("Expected error for input %q, but got none", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input %q: %v", test.input, err)
			}
			if !compareAST(result, test.expected) {
				t.Errorf("For input %q, expected %v, but got %v", test.input, test.expected, result)
			}
		}
	}
}

func compareAST(a, b ast.Expression) bool {
	// Implement a function to compare two ASTs for equality.
	// This is a simplified example and may need to be expanded for full coverage.
	switch a := a.(type) {
	case ast.Literal:
		b, ok := b.(ast.Literal)
		return ok && a.Value == b.Value
	case ast.Identifier:
		b, ok := b.(ast.Identifier)
		return ok && a.Name == b.Name
	case ast.BinaryOp:
		b, ok := b.(ast.BinaryOp)
		return ok && a.Op == b.Op && compareAST(a.Left, b.Left) && compareAST(a.Right, b.Right)
	case ast.FunctionCall:
		b, ok := b.(ast.FunctionCall)
		if !ok || !compareAST(a.Name, b.Name) || len(a.Args) != len(b.Args) {
			return false
		}
		for i := range a.Args {
			if !compareAST(a.Args[i], b.Args[i]) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func TestParseBlock(t *testing.T) {
    tokens := tokenizer.Tokenize(`{
        f(a);
        x = y;
        f(x)
    }`, "test")
    expr, err := Parse(tokens)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    block, ok := expr.(ast.Block)
    if !ok {
        t.Fatalf("expected block, got %T", expr)
    }
    if len(block.Expressions) != 3 {
        t.Fatalf("expected 3 expressions, got %d", len(block.Expressions))
    }
}

func TestParseBlockWithOptionalSemicolon(t *testing.T) {
    tokens := tokenizer.Tokenize(`{
        f(a);
        x = y;
        f(x);
    }`, "test")
    expr, err := Parse(tokens)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    block, ok := expr.(ast.Block)
    if !ok {
        t.Fatalf("expected block, got %T", expr)
    }
    if len(block.Expressions) != 4 {
        t.Fatalf("expected 4 expressions, got %d", len(block.Expressions))
    }
    if _, ok := block.Expressions[3].(ast.Literal); !ok {
        t.Fatalf("expected last expression to be a literal, got %T", block.Expressions[3])
    }
}

func TestParseBlockWithError(t *testing.T) {
    tokens := tokenizer.Tokenize(`{
        f(a)
        x = y;
    }`, "test")
    _, err := Parse(tokens)
    if err == nil {
        t.Fatalf("expected error, got nil")
    }
}