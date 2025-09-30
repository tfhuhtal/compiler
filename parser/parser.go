package parser

import (
	"compiler/ast"
	"compiler/tokenizer"
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
	"fun",
	"return",
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
	return p.peekOffset(0)
}

func (p *Parser) peekOffset(n int) tokenizer.Token {
	if p.pos+n < len(p.tokens) && (p.pos+n) >= 0 {
		return p.tokens[p.pos+n]
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

func (p *Parser) consume(expected any) tokenizer.Token {
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
	}
}

func (p *Parser) parseBooleanLiteral() ast.BooleanLiteral {
	token := p.consume(nil)
	return ast.BooleanLiteral{
		Location: token.Location,
		Boolean:  token.Text,
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
		}
	}
	return left
}

func (p *Parser) parseFactor() ast.Expression {
	token := p.peek()
	var res ast.Expression
	if token.Type == "Punctuation" {
		if token.Text == "{" {
			p.consume("{")
			res = p.parseBlock()
			p.consume("}")
		} else if token.Text == "(" {
			res = p.parseParenthesised()
		} else {
			panic(fmt.Sprintf(
				"Unexpected token %v, expexted left brace at location %v",
				p.peek().Text, p.peek().Location))
		}
	} else if token.Text == "if" {
		res = p.parseIfExpression()
	} else if token.Text == "true" || token.Text == "false" {
		res = p.parseBooleanLiteral()
	} else if token.Text == "while" {
		res = p.parseWhileLoop()
	} else if token.Type == "IntLiteral" {
		res = p.parseIntLiteral()
		if p.peek().Type == "IntLiteral" {
			panic(fmt.Sprintf(
				"Two consecutive int literals %s, %s",
				p.peekOffset(-1).Text, p.peek().Text))
		}
	} else if token.Type == "Identifier" {
		if p.peekOffset(-1).Type == "Identifier" &&
			!contains(allowedIdentifiers, p.peekOffset(-1).Text) {
			panic("Not allowed Identifier: " + p.peekOffset(-1).Text)
		}
		res = p.parseIdentifier()
	} else if token.Type == "" {
		panic("Invalid end of code")
	}
	if p.peek().Text == "(" {
		res = p.parseFunctionCall(res)
	}
	return res
}

func (p *Parser) parseFunctionCall(callee ast.Expression) ast.Expression {
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
	return ast.FunctionCall{
		Location: loc,
		Name:     callee,
		Args:     args,
	}
}

func (p *Parser) parseExpression() ast.Expression {
	fmt.Println("goes herer")
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
		}
	}
	if !contains(allowedIdentifiers, p.peekOffset(-1).Text) &&
		p.peekOffset(-1).Type == "Identifier" && p.peek().Type == "Identifier" {
		panic("Not allowed expression: " + p.peekOffset(-1).Text)
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
		}
	}
	return p.parseExpression()
}

func (p *Parser) parseBlock() ast.Expression {
	var expressions []ast.Expression

	for {
		expression := p.parseTopExpression()

		if p.peek().Text == ";" || p.peekOffset(-1).Text == "}" &&
			!contains([]string{"", "}"}, p.peek().Text) {
			if p.peek().Text == ";" {
				p.consume(";")
			}
			expressions = append(expressions, expression)
			expression = nil
		}

		if p.peek().Text == "}" || p.peek().Type == "end" {
			endLoc := p.peek().Location
			return ast.Block{
				Location:    endLoc,
				Expressions: expressions,
				Result:      expression,
			}
		}

		if expression != nil {
			panic(fmt.Sprintf("Result but no block end: %v", expression))
		}
	}
}

func (p *Parser) parseParams() []ast.Expression {
	var params []ast.Expression
	for p.peek().Text != ")" {
		loc := p.peek().Location
		name := p.parseIdentifier()
		p.consume(":")
		typed := p.parseIdentifier()
		param := ast.Param{
			Name:     name,
			Type:     typed,
			Location: loc,
		}
		params = append(params, param)
		if p.peek().Text != "," {
			break
		}
		p.consume(",")
	}
	return params
}

func (p *Parser) parseFunctionDefinition() ast.Expression {
	loc := p.peek().Location
	p.consume("fun")
	name := p.parseIdentifier()
	p.consume("(")
	params := p.parseParams()
	p.consume(")")
	p.consume(":")
	resultType := p.parseIdentifier()
	body := p.parseBlock()
	fmt.Println(name, params)
	return ast.FunctionDefinition{
		Name:       name,
		Params:     params,
		ResultType: resultType,
		Body:       body,
		Location:   loc,
	}
}

func (p *Parser) parseModule() ast.Expression {
	loc := p.peek().Location
	var functionDefinitions []ast.Expression
	for p.peek().Text == "fun" {
		functionDefinitions = append(functionDefinitions, p.parseFunctionDefinition())
	}

	block := p.parseBlock()
	if len(functionDefinitions) == 0 {
		return block
	}

	return ast.Module{
		Functions: functionDefinitions,
		Block:     block,
		Location:  loc,
	}
}

func Parse(tokens []tokenizer.Token) ast.Expression {
	p := new(tokens)
	expr := p.parseModule()
	return expr
}

func new(tokens []tokenizer.Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}
