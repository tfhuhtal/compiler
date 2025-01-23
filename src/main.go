package main

import (
	"compiler/src/tokenizer"
	"fmt"
)

func main() {
	tokens := tokenizer.Tokenize("// bla \n if (kissa == true) {}", "test")
	fmt.Println(tokens)
}
