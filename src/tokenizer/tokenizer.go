package tokenizer

import (
	"fmt"
)

type Token struct {
	kind  string
	value string
}

func tokenize(source_code string) []string {
	fmt.Println("Tokenizing source code", source_code)
	return []string{}
}
