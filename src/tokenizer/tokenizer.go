package tokenizer

import (
	"regexp"
)

type TokenType string

const (
	Integer     TokenType = "Integer"
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

func Tokenize(source_code string, file string) []Token {
	var tokens []Token
	line, column := 1, 1

	// Define regex patterns for different token types
	tokenPatterns := map[TokenType]*regexp.Regexp{
		Integer:     regexp.MustCompile(`^\d+`),
		Operator:    regexp.MustCompile(`^(==|!=|<=|>=|[+\-*/=<>])`),
		Punctuation: regexp.MustCompile(`^[(),{};]`),
		Identifier:  regexp.MustCompile(`^[a-zA-Z_]\w*`),
	}

	commentPattern := regexp.MustCompile(`^(//|#).*`)

	for len(source_code) > 0 {
		if source_code[0] == '\n' {
			line++
			column = 1
			source_code = source_code[1:]
			continue
		} else if source_code[0] == ' ' || source_code[0] == '\t' {
			column++
			source_code = source_code[1:]
			continue
		}

		if loc := commentPattern.FindStringIndex(source_code); loc != nil {
			endOfLine := regexp.MustCompile(`\n`).FindStringIndex(source_code)
			if endOfLine != nil {
				source_code = source_code[endOfLine[1]:]
				line++
				column = 1
			} else {
				break
			}
			continue
		}

		var matched bool
		for tokenType, pattern := range tokenPatterns {
			if loc := pattern.FindStringIndex(source_code); loc != nil {
				text := source_code[loc[0]:loc[1]]
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
				source_code = source_code[loc[1]:]
				matched = true
				break
			}
		}

		if !matched {
			column++
			source_code = source_code[1:]
		}
	}

	return tokens
}
