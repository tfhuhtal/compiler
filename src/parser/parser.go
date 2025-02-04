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

func contains(slice []string, item string) bool {
	for _, str := range slice {
		if str == item {
			return true
		}
	}
	return false
}

func peek(pos *int, tokens []tokenizer.Token) tokenizer.Token {
	if *pos < len(tokens) {
		return tokens[*pos]
	} else if len(tokens) == 0 {
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

func peekPrev(pos *int, tokens []tokenizer.Token) tokenizer.Token {
	if *pos-1 >= 0 {
		return tokens[*pos-1]
	} else {
		return tokenizer.Token{
			Location: tokens[len(tokens)-1].Location,
			Type:     "end",
			Text:     "",
		}
	}
}

func consume(pos *int, tokens []tokenizer.Token, expected interface{}) tokenizer.Token {
	token := peek(pos, tokens)

	if expected == nil {
		*pos++
		return token
	}

	if expectedStr, ok := expected.(string); ok {
		if token.Text != expectedStr {
			return tokenizer.Token{}
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
			return tokenizer.Token{}
		}
	}

	*pos++
	return token
}

func parseIntLiteral(pos *int, tokens []tokenizer.Token) ast.Literal {
	token := peek(pos, tokens)
	if token.Type != "IntLiteral" {
		return ast.Literal{}
	}
	consumedToken := consume(pos, tokens, nil)
	value, err := strconv.Atoi(consumedToken.Text)
	if err != nil {
		fmt.Println(err)
	}
	return ast.Literal{Value: value}
}

func parseIdentifier(pos *int, tokens []tokenizer.Token) ast.Identifier {
	token := peek(pos, tokens)
	if token.Type != "Identifier" {
		return ast.Identifier{}
	}
	consumedToken := consume(pos, tokens, nil)
	return ast.Identifier{Name: consumedToken.Text}
}

func parseParenthesised(pos *int, tokens []tokenizer.Token, list []string, allow bool) ast.Expression {
	consume(pos, tokens, "(")
	expr := parseExpression(pos, tokens, append([]string{")"}, list...), allow)
	consume(pos, tokens, ")")
	return expr
}

func parseIfExpression(pos *int, tokens []tokenizer.Token, list []string, allow bool) ast.Expression {
	consume(pos, tokens, "if")
	condition := parseExpression(pos, tokens, list, allow)
	consume(pos, tokens, "then")
	thenExpr := parseExpression(pos, tokens, list, allow)
	var elseExpr ast.Expression
	if peek(pos, tokens).Text == "else" {
		consume(pos, tokens, "else")
		elseExpr = parseExpression(pos, tokens, list, allow)
	}
	return ast.IfExpression{
		Condition: condition,
		Then:      thenExpr,
		Else:      elseExpr,
	}
}

func parseBooleanLiteral(pos *int, tokens []tokenizer.Token) ast.BooleanLiteral {
	token := consume(pos, tokens, nil)
	return ast.BooleanLiteral{
		Boolean: token.Text,
	}
}

func parseUnary(pos *int, tokens []tokenizer.Token, list []string, allow bool) ast.Expression {
	token := peek(pos, tokens)
	var operators []string
	for token.Text == "not" || token.Text == "-" {
		operators = append(operators, token.Text)
	}
	factor := parseFactor(pos, tokens, list, allow)
	if len(operators) > 0 {
		factor = ast.Unary{
			Ops: operators,
			Exp: factor,
		}
	}
	return factor
}

func parseWhileLoop(pos *int, tokens []tokenizer.Token) ast.Expression {
	consume(pos, tokens, "while")
	condition := parseExpression(pos, tokens, []string{"do"}, false)
	consume(pos, tokens, "do")
	looping := parseExpression(pos, tokens, []string{}, false)
	return ast.WhileLoop{
		Condition: condition,
		Looping:   looping,
	}
}

func parseTerm(pos *int, tokens []tokenizer.Token, list []string, allow bool) ast.Expression {
	left := parseUnary(pos, tokens, list, allow)
	if peek(pos, tokens).Text == "(" {
		args := parseFunction(pos, tokens, list, allow)
		left = ast.Function{
			Name: left,
			Args: args,
		}
	}
	for peek(pos, tokens).Text == "*" || peek(pos, tokens).Text == "/" {
		operatorToken := consume(pos, tokens, nil)
		operator := operatorToken.Text
		right := parseUnary(pos, tokens, list, allow)
		left = ast.BinaryOp{
			Left:  left,
			Op:    operator,
			Right: right,
		}
	}
	return left
}

func parseTermPrecedence(pos *int, tokens []tokenizer.Token, precedence int, list []string, allow bool) ast.Expression {
	var left ast.Expression
	if precedence == len(precedenceLevels)-1 {
		left = parseTerm(pos, tokens, list, allow)
	} else {
		left = parseTermPrecedence(pos, tokens, precedence+1, list, allow)
	}
	for contains(precedenceLevels[precedence], peek(pos, tokens).Text) {
		operatorToken := consume(pos, tokens, nil)
		operator := operatorToken.Text
		var right ast.Expression
		if precedence == len(precedenceLevels)-1 {
			right = parseTerm(pos, tokens, list, allow)
		} else {
			right = parseTermPrecedence(pos, tokens, precedence+1, list, allow)
		}
		left = ast.BinaryOp{
			Left:  left,
			Op:    operator,
			Right: right,
		}
	}
	return left
}

func parseFactor(pos *int, tokens []tokenizer.Token, list []string, allow bool) ast.Expression {
	token := peek(pos, tokens)
	if token.Text == "{" {
		return parseBlock(pos, tokens)
	} else if token.Text == "(" {
		return parseParenthesised(pos, tokens, list, allow)
	} else if token.Text == "if" {
		return parseIfExpression(pos, tokens, list, allow)
	} else if token.Text == "var" {
		return nil
	} else if token.Text == "true" || token.Text == "false" {
		return parseBooleanLiteral(pos, tokens)
	} else if token.Text == "while" {
		return parseWhileLoop(pos, tokens)
	} else if token.Type == "IntLiteral" {
		return parseIntLiteral(pos, tokens)
	} else if token.Type == "Identifier" {
		return parseIdentifier(pos, tokens)
	}
	return ast.Identifier{}
}

func parseFunction(pos *int, tokens []tokenizer.Token, list []string, allow bool) []ast.Expression {
	var args []ast.Expression
	consume(pos, tokens, "(")
	exprs := parseExpression(pos, tokens, append([]string{",", ")"}, list...), allow)
	args = append(args, exprs)
	for peek(pos, tokens).Text == "," {
		consume(pos, tokens, ",")
		exprs = parseExpression(pos, tokens, append([]string{",", ")"}, list...), allow)
		args = append(args, exprs)
	}
	consume(pos, tokens, ")")
	return args
}

func parseTopExpression(pos *int, tokens []tokenizer.Token, list []string, allow bool) ast.Expression {
	if peek(pos, tokens).Text == "var" {
		consume(pos, tokens, "var")
		decl := parseExpression(pos, tokens, append([]string{":"}, list...), allow)
		if peek(pos, tokens).Text == ":" {
			consume(pos, tokens, ":")
			typed := parseTypeExpression(pos, tokens)
			consume(pos, tokens, "=")
			declVal := parseExpression(pos, tokens, list, allow)
			return ast.Declaration{
				Variable: decl,
				Value:    declVal,
				Typed:    typed,
			}
		}
	}
	return parseExpression(pos, tokens, list, allow)
}

func parseExpression(pos *int, tokens []tokenizer.Token, list []string, allow bool) ast.Expression {
	precedence := 0
	left := parseTermPrecedence(pos, tokens, precedence+1, list, allow)
	for contains(precedenceLevels[precedence], peek(pos, tokens).Text) {
		operatorToken := consume(pos, tokens, nil)
		operator := operatorToken.Text
		right := parseTermPrecedence(pos, tokens, precedence+1, list, allow)
		left = ast.BinaryOp{
			Right: right,
			Op:    operator,
			Left:  left,
		}
	}
	if peek(pos, tokens).Text == "=" {
		operatorToken := consume(pos, tokens, nil)
		operator := operatorToken.Text
		right := parseExpression(pos, tokens, list, allow)
		left = ast.BinaryOp{
			Right: right,
			Op:    operator,
			Left:  left,
		}
	}
	if peek(pos, tokens).Type != "end" && !contains(list, peek(pos, tokens).Text) && peek(pos, tokens).Text != strconv.FormatBool(allow) {
		panic("Unexpected token: " + peek(pos, tokens).Text)
	}
	return left
}

func parseTypeExpression(pos *int, tokens []tokenizer.Token) ast.Expression {
	if peek(pos, tokens).Text == "(" {
		consume(pos, tokens, "(")
		var params []ast.Expression
		param := parseTypeExpression(pos, tokens)
		params = append(params, param)
		for peek(pos, tokens).Text == "," {
			param = parseTypeExpression(pos, tokens)
			params = append(params, param)
		}
		consume(pos, tokens, ")")
		consume(pos, tokens, "=>")
		res := parseTypeExpression(pos, tokens)
		return ast.FunctionTypeExpression{
			VariableTypes: params,
			ResultType:    res,
		}
	} else {
		return parseIdentifier(pos, tokens)
	}
}

func parseBlock(pos *int, tokens []tokenizer.Token) ast.Expression {
	consume(pos, tokens, "{")
	var seq []ast.Expression
	var res ast.Expression = nil
	if peek(pos, tokens).Text != "}" {
		line := parseTopExpression(pos, tokens, []string{";", "}"}, true)
		if peek(pos, tokens).Text != "}" && peek(pos, tokens).Text != ";" && line != nil && peekPrev(pos, tokens).Text != "}" {
			return nil
		}
		if peek(pos, tokens).Text != "}" {
			for peek(pos, tokens).Text == ";" || line != nil || peekPrev(pos, tokens).Text == "}" {
				consume(pos, tokens, ";")
				seq = append(seq, line)

				if peek(pos, tokens).Text == "}" {
					res = line
					break
				}
			}
		} else {
			res = line
		}
	}
	consume(pos, tokens, "}")
	if res == nil {
		res = ast.Literal{
			Value: nil,
		}
	}
	return ast.Block{
		Expressions: seq,
		Result:      res,
	}
}

func parseAll(pos *int, tokens []tokenizer.Token) []ast.Expression {
	var exprs []ast.Expression
	top := parseTopExpression(pos, tokens, []string{";"}, false)
	exprs = append(exprs, top)
	for peek(pos, tokens).Text == ";" {
		consume(pos, tokens, ";")
		if peek(pos, tokens).Type == "end" {
			break
		}
		top = parseTopExpression(pos, tokens, []string{";"}, false)
		exprs = append(exprs, top)
	}
	return exprs
}

func Parse(tokens []tokenizer.Token) []ast.Expression {
	pos := 0
	expr := parseAll(&pos, tokens)
	return expr
}
