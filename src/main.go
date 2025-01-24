package main

import (
	"compiler/src/tokenizer"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func call_compiler(sourceCode string, file string) []tokenizer.Token {
	tokens := tokenizer.Tokenize(sourceCode, file)
	return tokens
}

func main() {
	var command string
	var inputFile string
	var outputFile string
	var host string = "127.0.0.1"
	var port int = 3000
	var err error
	var tokens []tokenizer.Token

	for _, arg := range os.Args[1:] {
		if matched, _ := regexp.MatchString(`^--output=(.+)`, arg); matched {
			re := regexp.MustCompile(`^--output=(.+)`)
			matches := re.FindStringSubmatch(arg)
			if len(matches) > 1 {
				outputFile = matches[1]
			}
		} else if matched, _ := regexp.MatchString(`^--host=(.+)`, arg); matched {
			re := regexp.MustCompile(`^--host=(.+)`)
			matches := re.FindStringSubmatch(arg)
			if len(matches) > 1 {
				host = matches[1]
			}
		} else if matched, _ := regexp.MatchString(`^--port=(.+)`, arg); matched {
			re := regexp.MustCompile(`^--port=(.+)`)
			matches := re.FindStringSubmatch(arg)
			if len(matches) > 1 {
				port, err = strconv.Atoi(matches[1])
				if err != nil {
					fmt.Println("Error: Invalid port value")
					return
				}
			}
		} else if strings.HasPrefix(arg, "-") {
			fmt.Printf("Error: Unknown argument: %s\n", arg)
			return
		} else if command == "" {
			command = arg
		} else if inputFile == "" {
			inputFile = arg
		} else {
			fmt.Println("Error: Multiple input files not supported")
			return
		}
	}

	if command == "" {
		fmt.Fprintln(os.Stderr, "Error: command argument missing")
		return
	}

	if command == "compile" {
		tokens = call_compiler("if a <= bee then print_int(123)", inputFile)
		fmt.Println(tokens)
		fmt.Println(outputFile)
	} else if command == "serve" {

	}
}
