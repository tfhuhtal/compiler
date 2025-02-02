package tokenizer

import (
	"regexp"
)

type TokenType string

const (
	IntLiteral  TokenType = "IntLiteral"
	Operator    TokenType = "Operator"
	Punctuation TokenType = "Punctuation"
	Identifier  TokenType = "Identifier"
)

type SourceLocation struct {
	File   string
	Line   int
	Column int
}

type Token struct {
	Text     string
	Type     TokenType
	Location SourceLocation
}

func Tokenize(sourceCode string, file string) []Token {
	var tokens []Token
	line, column := 1, 1

	tokenPatterns := map[TokenType]*regexp.Regexp{
		IntLiteral:  regexp.MustCompile(`^\d+`),
		Operator:    regexp.MustCompile(`^(==|!=|<=|>=|[+\-*/=<>])`),
		Punctuation: regexp.MustCompile(`^[(),{};]`),
		Identifier:  regexp.MustCompile(`^[a-zA-Z_]\w*`),
	}

	commentPattern := regexp.MustCompile(`^(//|#).*`)

	for len(sourceCode) > 0 {
		if sourceCode[0] == '\n' {
			line++
			column = 1
			sourceCode = sourceCode[1:]
			continue
		} else if sourceCode[0] == ' ' || sourceCode[0] == '\t' {
			column++
			sourceCode = sourceCode[1:]
			continue
		}

		if loc := commentPattern.FindStringIndex(sourceCode); loc != nil {
			endOfLine := regexp.MustCompile(`\n`).FindStringIndex(sourceCode)
			if endOfLine != nil {
				sourceCode = sourceCode[endOfLine[1]:]
				line++
				column = 1
			} else {
				break
			}
			continue
		}

		var matched bool
		for tokenType, pattern := range tokenPatterns {
			if loc := pattern.FindStringIndex(sourceCode); loc != nil {
				text := sourceCode[loc[0]:loc[1]]
				tokens = append(tokens, Token{
					Text: text,
					Type: tokenType,
					Location: SourceLocation{
						File:   file,
						Line:   line,
						Column: column,
					},
				})
				column += len(text)
				sourceCode = sourceCode[loc[1]:]
				matched = true
				break
			}
		}

		if !matched {
			column++
			sourceCode = sourceCode[1:]
		}
	}

	return tokens
}
