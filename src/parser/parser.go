package parser

import (
	"compiler/src/ast"
	"compiler/src/tokenizer"
	"fmt"
	"strconv"
)

var precedenceLevels = [][]string{
	{"or"},
	{"and"},
	{"==", "!=", "<", "<=", ">", ">="},
	{"+", "-"},
	{"*", "/", "%"},
}

func peek(pos *int, tokens []tokenizer.Token) tokenizer.Token {
	if *pos < len(tokens) {
		return tokens[*pos]
	}

	if len(tokens) == 0 {
		return tokenizer.Token{
			Location: tokenizer.SourceLocation{File: "", Line: 1, Column: 1},
			Type:     "end",
			Text:     "",
		}
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

func parseIfExpression(pos *int, tokens []tokenizer.Token) (ast.Expression, error) {
	consume(pos, tokens, "if")
	condition, err := parseExpression(pos, tokens)
	if err != nil {
		return nil, err
	}
	consume(pos, tokens, "then")
	thenExpr, err := parseExpression(pos, tokens)
	if err != nil {
		return nil, err
	}
	var elseExpr ast.Expression
	if peek(pos, tokens).Text == "else" {
		consume(pos, tokens, "else")
		elseExpr, err = parseExpression(pos, tokens)
		if err != nil {
			return nil, err
		}
	}
	return ast.IfExpression{
		Condition: condition,
		Then:      thenExpr,
		Else:      elseExpr,
	}, nil
}

func parseFunctionCall(pos *int, tokens []tokenizer.Token, functionName ast.Identifier) (ast.Expression, error) {
	consume(pos, tokens, "(")
	var args []ast.Expression
	for {
		if peek(pos, tokens).Text == ")" {
			break
		}
		arg, err := parseExpression(pos, tokens)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		if peek(pos, tokens).Text != "," {
			break
		}
		consume(pos, tokens, ",")
	}
	consume(pos, tokens, ")")
	return ast.FunctionCall{
		Name: functionName,
		Args: args,
	}, nil
}

func parseUnary(pos *int, tokens []tokenizer.Token) (ast.Expression, error) {
	token := peek(pos, tokens)
	if token.Text == "not" || token.Text == "-" {
		consume(pos, tokens, nil)
		expr, err := parseUnary(pos, tokens)
		if err != nil {
			return nil, err
		}
		return ast.UnaryOp{
			Op:    token.Text,
			Right: expr,
		}, nil
	}
	return parsePrimary(pos, tokens)
}

func parsePrimary(pos *int, tokens []tokenizer.Token) (ast.Expression, error) {
	token := peek(pos, tokens)
	if token.Text == "{" {
		return parseBlock(pos, tokens)
	} else if token.Text == "(" {
		return parseParenthesised(pos, tokens)
	} else if token.Text == "if" {
		return parseIfExpression(pos, tokens)
	} else if token.Type == "IntLiteral" {
		return parseIntLiteral(pos, tokens)
	} else if token.Type == "Identifier" {
		identifier, err := parseIdentifier(pos, tokens)
		if err != nil {
			return nil, err
		}
		if peek(pos, tokens).Text == "(" {
			return parseFunctionCall(pos, tokens, identifier)
		}
		return identifier, nil
	}
	return ast.Identifier{}, fmt.Errorf("%v: expected an integer literal or an identifier", token.Location)
}

func parseBinary(pos *int, tokens []tokenizer.Token, precedence int) (ast.Expression, error) {
	left, err := parseUnary(pos, tokens)
	if err != nil {
		return nil, err
	}

	for precedence < len(precedenceLevels) {
		for _, op := range precedenceLevels[precedence] {
			if peek(pos, tokens).Text == op {
				operatorToken, err := consume(pos, tokens, nil)
				if err != nil {
					return nil, err
				}
				right, err := parseBinary(pos, tokens, precedence+1)
				if err != nil {
					return nil, err
				}
				left = ast.BinaryOp{
					Left:  left,
					Op:    operatorToken.Text,
					Right: right,
				}
			}
		}
		precedence++
	}
	return left, nil
}

func parseAssignment(pos *int, tokens []tokenizer.Token) (ast.Expression, error) {
	left, err := parseBinary(pos, tokens, 0)
	if err != nil {
		return nil, err
	}

	if peek(pos, tokens).Text == "=" {
		consume(pos, tokens, "=")
		right, err := parseAssignment(pos, tokens)
		if err != nil {
			return nil, err
		}
		left = ast.Assignment{
			Left:  left,
			Right: right,
		}
	}
	return left, nil
}

func parseExpression(pos *int, tokens []tokenizer.Token) (ast.Expression, error) {
	return parseAssignment(pos, tokens)
}

func parseBlock(pos *int, tokens []tokenizer.Token) (ast.Expression, error) {
	consume(pos, tokens, "{")
	var expressions []ast.Expression
	for {
		if peek(pos, tokens).Text == "}" {
			break
		}
		expr, err := parseExpression(pos, tokens)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expr)
		if peek(pos, tokens).Text == ";" {
			consume(pos, tokens, ";")
		} else {
			break
		}
	}
	consume(pos, tokens, "}")
	if len(expressions) == 0 || peek(pos, tokens).Text == ";" {
		expressions = append(expressions, ast.Literal{Value: nil})
	}
	return ast.Block{Expressions: expressions}, nil
}

func Parse(tokens []tokenizer.Token) (ast.Expression, error) {
	pos := 0
	expr, err := parseExpression(&pos, tokens)
	return expr, err
}
