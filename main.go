package main

import (
	"compiler/ir"
	"compiler/irgenerator"
	"compiler/parser"
	"compiler/tokenizer"
	"compiler/typechecker"
	"compiler/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Input struct {
	Command string `json:"command"`
	Code    string `json:"code,omitempty"`
}

type Result struct {
	Program string `json:"program,omitempty"`
	Error   string `json:"error,omitempty"`
}

func callCompiler(sourceCode string, file string) any {
	fmt.Println(sourceCode)
	fmt.Println("=================================================")
	tokens := tokenizer.Tokenize(sourceCode, file)
	fmt.Println(tokens)
	fmt.Println("=================================================")
	p := parser.New(tokens)
	res := p.Parse()
	typechecker.Type(res)
	fmt.Println(res)
	fmt.Println("=================================================")
	var rootTypes = make(map[irgenerator.IRVar]utils.Type)
	rootTypes["+"] = utils.Int{}
	rootTypes["*"] = utils.Int{}
	var instructions []ir.Instruction

	instructions = irgenerator.Generate(rootTypes, res[0])

	fmt.Println(instructions)
	return instructions
}

func handler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var input Input
	if err := json.Unmarshal(body, &input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var result Result

	switch input.Command {
	/*case "compile":*/
	/*compiled := callCompiler(input.Code, "test")*/
	/*result.Program = base64.StdEncoding.EncodeToString(compiled)*/
	case "ping":
		// no operation
	default:
		result.Error = fmt.Sprintf("Unknown command: %s", input.Command)
	}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func runServer(host string, port int) {
	address := fmt.Sprintf("%s:%d", host, port)
	fmt.Println("Server running on: ", address)
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		fmt.Println("Error: ", err)
	}

}

func main() {
	var command string
	var inputFile string
	var outputFile string
	var host string = "127.0.0.1"
	var port int = 3000
	var err error

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
		callCompiler("1 + 2 * 3;", inputFile)
		fmt.Print(outputFile)
	} else if command == "serve" {
		runServer(host, port)
	} else {
		fmt.Fprintln(os.Stderr, "Error: Unknown command")
	}
}
