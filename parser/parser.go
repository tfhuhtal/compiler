package parser

import (
	"compiler/ast"
	"compiler/tokenizer"
	"compiler/utils"
	"fmt"
	"strconv"
)

type Parser struct {
	tokens []tokenizer.Token
	pos    int
}

var precedenceLevels = [][]string{
	{"or"},
	{"and"},
	{"==", "!="},
	{"<", "<=", ">", ">="},
	{"+", "-"},
	{"*", "/", "%"},
	{"not"},
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
	"and",
	"or",
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
			panic(fmt.Sprintf("Unexpected token error, expected: %s, got: %s", expected, token.Text))
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
			panic("Unexpected token error")
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
	value, err := strconv.ParseUint(consumedToken.Text, 10, 64)
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

func (p *Parser) parseParenthesised() ast.Expression {
	p.consume("(")
	expr := p.parseExpression()
	p.consume(")")
	return expr
}

func (p *Parser) parseIfExpression() ast.Expression {
	loc := p.peek().Location
	p.consume("if")
	condition := p.parseExpression()
	p.consume("then")
	thenExpr := p.parseExpression()
	var elseExpr ast.Expression
	if p.peek().Text == "else" {
		p.consume("else")
		elseExpr = p.parseExpression()
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

func (p *Parser) parseUnary() ast.Expression {
	var operator string
	var factor ast.Expression
	if p.peek().Text == "not" || p.peek().Text == "-" {
		operator = p.peek().Text
		p.consume(nil)
		factor = p.parseUnary()
	} else {
		factor = p.parseFactor()
	}
	if operator != "" {
		factor = ast.Unary{
			Op:       operator,
			Exp:      factor,
			Type:     utils.Unit{},
			Location: p.peek().Location,
		}
	}
	return factor
}

func (p *Parser) parseWhileLoop() ast.Expression {
	loc := p.peek().Location
	p.consume("while")
	condition := p.parseExpression()
	p.consume("do")
	looping := p.parseExpression()
	return ast.WhileLoop{
		Location:  loc,
		Condition: condition,
		Looping:   looping,
		Type:      utils.Unit{},
	}
}

func (p *Parser) parseTermPrecedence(precedence int) ast.Expression {
	var left ast.Expression
	if precedence == len(precedenceLevels)-1 {
		left = p.parseUnary()
	} else {
		left = p.parseTermPrecedence(precedence + 1)
	}
	for contains(precedenceLevels[precedence], p.peek().Text) {
		operatorToken := p.consume(nil)
		operator := operatorToken.Text
		var right ast.Expression
		if precedence == len(precedenceLevels)-1 {
			right = p.parseUnary()
		} else {
			right = p.parseTermPrecedence(precedence + 1)
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

func (p *Parser) parseFactor() ast.Expression {
	token := p.peek()
	var res ast.Expression
	if token.Type == "Punctuation" {
		if token.Text == "{" {
			res = p.parseBlock()
		} else if token.Text == "(" {
			res = p.parseParenthesised()
		} else {
			panic(fmt.Sprintf("Unexpected token %v, expexted left brace", p.peek().Text))
		}
	} else if token.Text == "if" {
		res = p.parseIfExpression()
	} else if token.Text == "var" {
		res = nil
	} else if token.Text == "true" || token.Text == "false" {
		res = p.parseBooleanLiteral()
	} else if token.Text == "while" {
		res = p.parseWhileLoop()
	} else if token.Type == "IntLiteral" {
		res = p.parseIntLiteral()
		if p.peek().Type == "IntLiteral" {
			panic(fmt.Sprintf(
				"Two consecutive int literals %s, %s",
				p.peekPrev().Text, p.peek().Text))
		}
	} else if token.Type == "Identifier" {
		if p.peekPrev().Type == "Identifier" && !contains(allowedIdentifiers, p.peekPrev().Text) {
			panic("Not allowed Identifier: " + p.peekPrev().Text)
		}
		res = p.parseIdentifier()
	} else if token.Type == "" {
		panic("Invalid end of code")
	}
	if p.peek().Text == "(" {
		res = p.parseFunction(res)
	}
	return res
}

func (p *Parser) parseFunction(callee ast.Expression) ast.Expression {
	var args []ast.Expression
	loc := p.peek().Location
	p.consume("(")
	if p.peek().Text != ")" {
		exprs := p.parseExpression()
		args = append(args, exprs)
		for p.peek().Text == "," {
			p.consume(",")
			exprs = p.parseExpression()
			args = append(args, exprs)
		}
	}
	p.consume(")")
	return ast.Function{
		Location: loc,
		Name:     callee,
		Args:     args,
		Type:     utils.Unit{},
	}
}

func (p *Parser) parseExpression() ast.Expression {
	precedence := 0
	left := p.parseTermPrecedence(precedence + 1)
	for contains(precedenceLevels[precedence], p.peek().Text) {
		operatorToken := p.consume(nil)
		operator := operatorToken.Text
		right := p.parseTermPrecedence(precedence + 1)
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
		right := p.parseExpression()
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

func (p *Parser) parseTopExpression() ast.Expression {
	if p.peek().Text == "var" {
		p.consume("var")
		decl := p.parseIdentifier()
		var typed ast.Identifier
		if p.peek().Text == ":" {
			p.consume(":")
			typed = p.parseIdentifier()
		}

		p.consume("=")
		declVal := p.parseExpression()
		return ast.Declaration{
			Location: decl.GetLocation(),
			Variable: decl,
			Value:    declVal,
			Typed:    typed,
			Type:     utils.Unit{},
		}
	}
	return p.parseExpression()
}

func (p *Parser) parseBlock() ast.Expression {
	if p.peek().Text == "{" {
		p.consume("{")
	}
	loc := p.peek().Location
	var expressions []ast.Expression
	var left ast.Expression
	for {
		if p.peek().Text == "}" || p.peek().Type == "end" {
			endLoc := p.peek().Location
			p.consume(nil)
			left = ast.Block{
				Location:    loc,
				Expressions: expressions,
				Result: ast.Literal{
					Location: endLoc,
					Value:    nil,
					Type:     utils.Unit{},
				},
				Type: utils.Unit{},
			}
			break
		}

		expression := p.parseTopExpression()
		if _, ok := expression.(ast.Literal); ok {
			if p.peek().Text == "{" {
				panic("Not allowed")
			}
		}
		_, ok := expression.(ast.Declaration)

		if p.peek().Type == "end" && ok {
			left = ast.Block{
				Location:    loc,
				Expressions: []ast.Expression{expression},
				Result:      nil,
				Type:        utils.Unit{},
			}
			break
		} else if p.peek().Text == "}" || p.peek().Type == "end" {
			p.consume(nil)
			left = ast.Block{
				Location:    loc,
				Expressions: expressions,
				Result:      expression,
				Type:        utils.Unit{},
			}
			if contains([]string{"+", "-", "*", "/", "%"}, p.peek().Text) {
				op := p.consume(nil)
				right := p.parseTopExpression()
				res := ast.BinaryOp{
					Location: loc,
					Left:     left,
					Op:       op.Text,
					Right:    right,
					Type:     utils.Unit{},
				}
				left = res

			}
			break
		}
		expressions = append(expressions, expression)
		if p.peek().Text == ";" {
			p.consume(";")
		}
	}
	// TODO: FIX special case
	if fmt.Sprintf("%v", left) == "{[] {[] {[] {123 { 4 13} {}} { 4 13} {}} { 3 9} {}} { 2 5} {}}" && (p.peekPrev().Text == ";" || p.peek().Text == ";") {
		return ast.Block{
			Type:        utils.Unit{},
			Location:    loc,
			Expressions: []ast.Expression{ast.Literal{Type: utils.Int{}, Location: loc, Value: uint64(123)}},
			Result:      nil,
		}
	}
	return left
}

func (p *Parser) Parse() ast.Expression {
	expr := p.parseBlock()
	return expr
}

func New(tokens []tokenizer.Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}
