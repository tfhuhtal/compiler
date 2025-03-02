package parser

import (
	"compiler/ast"
	"compiler/tokenizer"
	"compiler/utils"
	"strconv"
)

type Parser struct {
	tokens []tokenizer.Token
	pos    int
}

var precedenceLevels = [][]string{
	{"or"},
	{"and"},
	{"==", "!=", "<", "<=", ">", ">="},
	{"+", "-"},
	{"*", "/", "%"},
}

var allowedIdentifiers = []string{
	"var",
	"if",
	"else",
	"then",
	"while",
	"do",
	"not",
	"true",
	"false",
}

func contains(slice []string, item string) bool {
	for _, str := range slice {
		if str == item {
			return true
		}
	}
	return false
}

func (p *Parser) peek() tokenizer.Token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	} else if len(p.tokens) == 0 {
		return tokenizer.Token{
			Location: tokenizer.SourceLocation{File: "", Line: 1, Column: 1},
			Type:     "end",
			Text:     "",
		}
	}
	// If we're at or past the end, return an "end" token.
	return tokenizer.Token{
		Location: p.tokens[len(p.tokens)-1].Location,
		Type:     "end",
		Text:     "",
	}
}

func (p *Parser) peekPrev() tokenizer.Token {
	if p.pos-1 >= 0 {
		return p.tokens[p.pos-1]
	} else {
		return tokenizer.Token{
			Location: p.tokens[len(p.tokens)-1].Location,
			Type:     "end",
			Text:     "",
		}
	}
}

func (p *Parser) consume(expected interface{}) tokenizer.Token {
	token := p.peek()

	if expected == nil {
		p.pos++
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

	p.pos++
	return token
}

func (p *Parser) parseIntLiteral() ast.Literal {
	token := p.peek()
	if token.Type != "IntLiteral" {
		panic("Not int literal")
	}
	consumedToken := p.consume(nil)
	value, err := strconv.Atoi(consumedToken.Text)
	if err != nil {
		panic(err)
	}
	return ast.Literal{
		Location: consumedToken.Location,
		Value:    value,
		Type:     utils.Unit{},
	}
}

func (p *Parser) parseIdentifier() ast.Identifier {
	token := p.peek()
	if token.Type != "Identifier" {
		panic("Not identifier")
	}
	consumedToken := p.consume(nil)
	return ast.Identifier{
		Location: consumedToken.Location,
		Name:     consumedToken.Text,
		Type:     utils.Unit{},
	}
}

func (p *Parser) parseParenthesised(list []string, allow bool) ast.Expression {
	p.consume("(")
	expr := p.parseExpression(append([]string{")"}, list...), allow)
	p.consume(")")
	return expr
}

func (p *Parser) parseIfExpression(list []string, allow bool) ast.Expression {
	loc := p.peek().Location
	p.consume("if")
	condition := p.parseExpression(list, allow)
	p.consume("then")
	thenExpr := p.parseExpression(list, allow)
	var elseExpr ast.Expression
	if p.peek().Text == "else" {
		p.consume("else")
		elseExpr = p.parseExpression(list, allow)
	}
	return ast.IfExpression{
		Location:  loc,
		Condition: condition,
		Then:      thenExpr,
		Else:      elseExpr,
		Type:      utils.Unit{},
	}
}

func (p *Parser) parseBooleanLiteral() ast.BooleanLiteral {
	token := p.consume(nil)
	return ast.BooleanLiteral{
		Location: token.Location,
		Boolean:  token.Text,
		Type:     utils.Unit{},
	}
}

func (p *Parser) parseUnary(list []string, allow bool) ast.Expression {
	token := p.peek()
	var operators []string
	for token.Text == "not" || token.Text == "-" {
		operators = append(operators, token.Text)
	}
	factor := p.parseFactor(list, allow)
	if len(operators) > 0 {
		factor = ast.Unary{
			Location: token.Location,
			Ops:      operators,
			Exp:      factor,
			Type:     utils.Unit{},
		}
	}
	return factor
}

func (p *Parser) parseWhileLoop() ast.Expression {
	loc := p.peek().Location
	p.consume("while")
	condition := p.parseExpression([]string{"do"}, false)
	p.consume("do")
	looping := p.parseExpression([]string{}, false)
	return ast.WhileLoop{
		Location:  loc,
		Condition: condition,
		Looping:   looping,
		Type:      utils.Unit{},
	}
}

func (p *Parser) parseTerm(list []string, allow bool) ast.Expression {
	left := p.parseUnary(list, allow)
	for p.peek().Text == "*" || p.peek().Text == "/" || p.peek().Text == "%" {
		operatorToken := p.consume(nil)
		operator := operatorToken.Text
		right := p.parseUnary(list, allow)
		left = ast.BinaryOp{
			Location: left.GetLocation(),
			Left:     left,
			Op:       operator,
			Right:    right,
			Type:     utils.Unit{},
		}
	}
	return left
}

func (p *Parser) parseTermPrecedence(precedence int, list []string, allow bool) ast.Expression {
	var left ast.Expression
	if precedence == len(precedenceLevels)-1 {
		left = p.parseTerm(list, allow)
	} else {
		left = p.parseTermPrecedence(precedence+1, list, allow)
	}
	for contains(precedenceLevels[precedence], p.peek().Text) {
		operatorToken := p.consume(nil)
		operator := operatorToken.Text
		var right ast.Expression
		if precedence == len(precedenceLevels)-1 {
			right = p.parseTerm(list, allow)
		} else {
			right = p.parseTermPrecedence(precedence+1, list, allow)
		}
		left = ast.BinaryOp{
			Location: left.GetLocation(),
			Left:     left,
			Op:       operator,
			Right:    right,
			Type:     utils.Unit{},
		}
	}
	return left
}

func (p *Parser) parseFactor(list []string, allow bool) ast.Expression {
	token := p.peek()
	var res ast.Expression
	if token.Text == "{" {
		res = p.parseBlock()
	} else if token.Text == "(" {
		res = p.parseParenthesised(list, allow)
	} else if token.Text == "if" {
		res = p.parseIfExpression(list, allow)
	} else if token.Text == "var" {
		res = nil
	} else if token.Text == "true" || token.Text == "false" {
		res = p.parseBooleanLiteral()
	} else if token.Text == "while" {
		res = p.parseWhileLoop()
	} else if token.Type == "IntLiteral" {
		res = p.parseIntLiteral()
	} else if token.Type == "Identifier" {
		if p.peekPrev().Type == "Identifier" && !contains(allowedIdentifiers, p.peekPrev().Text) {
			panic("Not allowed Identifier: " + p.peekPrev().Text)
		}
		res = p.parseIdentifier()
	}
	if p.peek().Text == "(" {
		res = p.parseFunction(list, allow, res)
	}
	return res
}

func (p *Parser) parseFunction(list []string, allow bool, callee ast.Expression) ast.Expression {
	var args []ast.Expression
	loc := p.peek().Location
	p.consume("(")
	exprs := p.parseExpression(append([]string{",", ")"}, list...), allow)
	args = append(args, exprs)
	for p.peek().Text == "," {
		p.consume(",")
		exprs = p.parseExpression(append([]string{",", ")"}, list...), allow)
		args = append(args, exprs)
	}
	p.consume(")")
	return ast.Function{
		Location: loc,
		Name:     callee,
		Args:     args,
		Type:     utils.Unit{},
	}
}

func (p *Parser) parseExpression(list []string, allow bool) ast.Expression {
	precedence := 0
	left := p.parseTermPrecedence(precedence+1, list, allow)
	for contains(precedenceLevels[precedence], p.peek().Text) {
		operatorToken := p.consume(nil)
		operator := operatorToken.Text
		right := p.parseTermPrecedence(precedence+1, list, allow)
		left = ast.BinaryOp{
			Location: left.GetLocation(),
			Right:    right,
			Op:       operator,
			Left:     left,
			Type:     utils.Unit{},
		}
	}
	if p.peek().Text == "=" {
		operatorToken := p.consume(nil)
		operator := operatorToken.Text
		right := p.parseExpression(list, allow)
		left = ast.BinaryOp{
			Location: left.GetLocation(),
			Right:    right,
			Op:       operator,
			Left:     left,
			Type:     utils.Unit{},
		}
	}
	if !contains(allowedIdentifiers, p.peekPrev().Text) && p.peekPrev().Type == "Identifier" && p.peek().Type == "Identifier" {
		panic("Not allowed expression: " + p.peekPrev().Text)
	}
	return left
}

func (p *Parser) parseTypeExpression() ast.Expression {
	if p.peek().Text == "(" {
		p.consume("(")
		var params []ast.Expression
		param := p.parseTypeExpression()
		params = append(params, param)
		for p.peek().Text == "," {
			param = p.parseTypeExpression()
			params = append(params, param)
		}
		p.consume(")")
		p.consume("=>")
		res := p.parseTypeExpression()
		return ast.FunctionTypeExpression{
			Location:      p.peek().Location,
			VariableTypes: params,
			ResultType:    res,
			Type:          utils.Unit{},
		}
	} else {
		return p.parseIdentifier()
	}
}

func (p *Parser) parseTopExpression(list []string, allow bool) ast.Expression {
	if p.peek().Text == "var" {
		p.consume("var")
		decl := p.parseExpression(append([]string{":"}, list...), allow)
		if p.peek().Text == ":" {
			p.consume(":")
			typed := p.parseTypeExpression()
			p.consume("=")
			declVal := p.parseExpression(list, allow)
			return ast.Declaration{
				Location: decl.GetLocation(),
				Variable: decl,
				Value:    declVal,
				Typed:    typed,
				Type:     utils.Unit{},
			}
		}
	}
	return p.parseExpression(list, allow)
}

func (p *Parser) parseBlock() ast.Expression {
	loc := p.peek().Location
	if p.peek().Text == "{" {
		p.consume("{")
	}
	var seq []ast.Expression
	var res ast.Expression = nil
	if p.peek().Text != "}" {
		line := p.parseTopExpression([]string{";", "}"}, true)
		if p.peek().Text != "}" {
			for p.peek().Text == ";" || line != nil || p.peekPrev().Text == "}" {
				p.consume(";")
				if p.peek().Text == "}" {
					if p.peekPrev().Text == ";" {
						seq = append(seq, line)
						res = ast.Literal{
							Location: p.peek().Location,
							Value:    nil,
							Type:     utils.Unit{},
						}
					} else {
						res = line
					}
					break
				}
				seq = append(seq, line)
				line = p.parseTopExpression([]string{";", "}"}, true)
			}
		} else {
			res = line
		}
	}
	p.consume("}")
	return ast.Block{
		Location:    loc,
		Expressions: seq,
		Result:      res,
		Type:        utils.Unit{},
	}
}

func (p *Parser) Parse() ast.Expression {
	expr := p.parseBlock()
	return expr
}

func New(tokens []tokenizer.Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}
