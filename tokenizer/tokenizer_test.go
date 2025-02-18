package tokenizer

import (
	"testing"
)

// Define a special SourceLocation object
var L = SourceLocation{}

func (t Token) Equal(other Token) bool {
	if t.Text != other.Text || t.Type != other.Type {
		return false
	}
	if t.Location == L || other.Location == L {
		return true
	}
	return t.Location == other.Location
}

func TestTokenize(t *testing.T) {
	tokens := Tokenize("if  3\nwhile + - 4 5 6 kissa", "test")
	expected := []Token{
		{Text: "if", Type: Identifier, Location: L},
		{Text: "3", Type: IntLiteral, Location: L},
		{Text: "while", Type: Identifier, Location: L},
		{Text: "+", Type: Operator, Location: L},
		{Text: "-", Type: Operator, Location: L},
		{Text: "4", Type: IntLiteral, Location: L},
		{Text: "5", Type: IntLiteral, Location: L},
		{Text: "6", Type: IntLiteral, Location: L},
		{Text: "kissa", Type: Identifier, Location: L},
	}
	if len(tokens) != len(expected) {
		t.Errorf("Expected %d tokens, got %d", len(expected), len(tokens))
	}
	for i := range tokens {
		if !tokens[i].Equal(expected[i]) {
			t.Errorf("Expected token %v, got %v", expected[i], tokens[i])
		}
	}
}
