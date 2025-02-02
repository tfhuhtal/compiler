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
	if token.Type != "IntLiteral" {
		return ast.Literal{}, fmt.Errorf("%v: expected an integer literal", token.Location)
	}
	consumedToken, err := consume(pos, tokens, nil)
	if err != nil {
		return ast.Literal{}, err
	}
	value, err := strconv.Atoi(consumedToken.Text)
	return ast.Literal{Value: value}, err
}

func parseIdentifier(pos *int, tokens []tokenizer.Token) (ast.Identifier, error) {
	token := peek(pos, tokens)
	if token.Type != "Identifier" {
		return ast.Identifier{}, fmt.Errorf("%v: expected an identifier", token.Location)
	}
	consumedToken, err := consume(pos, tokens, nil)
	if err != nil {
		return ast.Identifier{}, err
	}
	return ast.Identifier{Name: consumedToken.Text}, nil
}

func parseParenthesised(pos *int, tokens []tokenizer.Token) (ast.Expression, error) {
	consume(pos, tokens, "(")
	expr, err := parseExpression(pos, tokens)
	consume(pos, tokens, ")")
	return expr, err
}

func parseFactor(pos *int, tokens []tokenizer.Token) (ast.Expression, error) {
	token := peek(pos, tokens)
	if token.Text == "(" {
		return parseParenthesised(pos, tokens)
	} else if token.Type == "IntLiteral" {
		return parseIntLiteral(pos, tokens)
	} else if token.Type == "Identifier" {
		return parseIdentifier(pos, tokens)
	}
	return ast.Identifier{}, fmt.Errorf("%v: expected an integer literal or an identifier", token.Location)
}

func parseTerm(pos *int, tokens []tokenizer.Token) (ast.Expression, error) {
	left, err := parseFactor(pos, tokens)
	for peek(pos, tokens).Text == "*" || peek(pos, tokens).Text == "/" {
		operatorToken, err := consume(pos, tokens, nil)
		if err != nil {
			return ast.BinaryOp{}, err
		}
		operator := operatorToken.Text
		right, err := parseFactor(pos, tokens)
		if err != nil {
			return ast.BinaryOp{}, err
		}
		left = ast.BinaryOp{
			Left:  left,
			Op:    operator,
			Right: right,
		}
	}
	return left, err
}

func parseExpression(pos *int, tokens []tokenizer.Token) (ast.Expression, error) {
	left, err := parseTerm(pos, tokens)
	if err != nil {
		return ast.BinaryOp{}, err
	}

	// looping while there is no more operation tokens
	for peek(pos, tokens).Text == "+" || peek(pos, tokens).Text == "-" {
		operatorToken, err := consume(pos, tokens, nil)
		if err != nil {
			return ast.BinaryOp{}, err
		}
		operator := operatorToken.Text
		right, err := parseTerm(pos, tokens)
		if err != nil {
			return ast.BinaryOp{}, err
		}
		left = ast.BinaryOp{
			Left:  left,
			Op:    operator,
			Right: right,
		}
	}
	return left, err
}

func Parse(tokens []tokenizer.Token) (ast.Expression, error) {
	pos := 0
	expr, err := parseExpression(&pos, tokens)
	return expr, err
}
