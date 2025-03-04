package main

import (
	"compiler/asmgenerator"
	"compiler/irgenerator"
	"compiler/parser"
	"compiler/tokenizer"
	"compiler/typechecker"
	"compiler/utils"
	"encoding/base64"
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

func callCompiler(sourceCode string, file string) string {
	tokens := tokenizer.Tokenize(sourceCode, file)
	p := parser.New(tokens)
	res := p.Parse()
	typechecker.Type(res)
	var rootTypes = make(map[irgenerator.IRVar]utils.Type)
	rootTypes["+"] = utils.Int{}
	rootTypes["*"] = utils.Int{}
	rootTypes[">"] = utils.Bool{}
	rootTypes["%"] = utils.Int{}
	rootTypes["=="] = utils.Int{}
	rootTypes["/"] = utils.Int{}
	rootTypes["<="] = utils.Int{}
	rootTypes["<"] = utils.Int{}
	rootTypes[">="] = utils.Int{}
	rootTypes["!="] = utils.Int{}

	gen := irgenerator.NewIRGenerator(rootTypes)
	instructions := gen.Generate(res)

	asm := asmgenerator.GenerateASM(instructions)
	return asm
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
	case "compile":
		compiled := callCompiler(input.Code, "test")
		result.Program = base64.StdEncoding.EncodeToString([]byte(compiled))
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
	fmt.Println("Server running on:", address)

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("pong\n"))
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/", handler)

	if err := http.ListenAndServe(address, nil); err != nil {
		fmt.Println("Error:", err)
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
		asm := callCompiler("var n: Int = read_int();while n < 100 do {print_int(n);n = n + 1;}", inputFile)
		os.WriteFile(outputFile, []byte(asm), 0644)
	} else if command == "serve" {
		runServer(host, port)
	} else {
		fmt.Fprintln(os.Stderr, "Error: Unknown command")
	}
}
