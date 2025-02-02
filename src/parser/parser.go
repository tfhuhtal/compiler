package parser

import (
	"compiler/src/ast"
	"compiler/src/tokenizer"
	"fmt"
	"strconv"
)

func peek(pos *int, tokens []tokenizer.Token) tokenizer.Token {
	if *pos < len(tokens) {
		return tokens[*pos]
	}
	// If we're at or past the end, return an "end" token.
	return tokenizer.Token{
		Location: tokens[len(tokens)-1].Location,
		Type:     "end",
		Text:     "",
	}
}

func consume(pos *int, tokens []tokenizer.Token, expected interface{}) (tokenizer.Token, error) {
	token := peek(pos, tokens)

	fmt.Println(expected)
	if expected == nil {
		*pos++
		return token, nil
	}

	if expectedStr, ok := expected.(string); ok {
		if token.Text != expectedStr {
			return tokenizer.Token{}, fmt.Errorf("%v: expected \"%s\"", token.Location, expectedStr)
		}
	}

	if expectedList, ok := expected.([]string); ok {
		matched := false
		for _, e := range expectedList {
			if token.Text == e {
				matched = true
				break
			}
		}
		if !matched {
			return tokenizer.Token{}, fmt.Errorf("%v: expected \"%v\"", token.Location, expectedList)
		}
	}

	*pos++
	return token, nil
}

func parseIntLiteral(pos *int, tokens []tokenizer.Token) (ast.Literal, error) {
	token := peek(pos, tokens)
	if token.Type != "Integer" {
		return ast.Literal{}, fmt.Errorf("%v: expected an integer literal", token.Location)
	}
	consumedToken, err := consume(pos, tokens, nil)
	if err != nil {
		return ast.Literal{}, err
	}
	fmt.Println(consumedToken)
	value, err := strconv.Atoi(consumedToken.Text)
	return ast.Literal{Value: value}, err
}

func parseExpression(pos *int, tokens []tokenizer.Token) (ast.BinaryOp, error) {
	fmt.Println(*pos)
	left, err := parseIntLiteral(pos, tokens)
	if err != nil {
		return ast.BinaryOp{}, err
	}

	operatorToken, err := consume(pos, tokens, []string{"+", "-"})
	if err != nil {
		return ast.BinaryOp{}, err
	}

	right, err := parseIntLiteral(pos, tokens)
	if err != nil {
		return ast.BinaryOp{}, err
	}

	return ast.BinaryOp{
		Left:  left,
		Op:    operatorToken.Text,
		Right: right,
	}, nil
}

func Parse(tokens []tokenizer.Token) (ast.Expression, error) {
	pos := 0
	expr, err := parseExpression(&pos, tokens)
	return expr, err
}
