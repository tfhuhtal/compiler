package main

import (
	"compiler/asmgenerator"
	"compiler/assembler"
	"compiler/irgenerator"
	"compiler/parser"
	"compiler/tokenizer"
	"compiler/typechecker"
	"compiler/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func callCompiler(sourceCode string, file string) []byte {
	var output []byte
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
			output = []byte(fmt.Sprintf("compiler error: %s", r))
		}
	}()

	tokens := tokenizer.Tokenize(sourceCode, file)
	p := parser.New(tokens)
	res := p.Parse()
	typechecker.Type(res)

	rootTypes := map[irgenerator.IRVar]utils.Type{
		"+":   utils.Int{},
		"*":   utils.Int{},
		"/":   utils.Int{},
		"%":   utils.Int{},
		"-":   utils.Int{},
		">":   utils.Bool{},
		"==":  utils.Bool{},
		"<=":  utils.Bool{},
		"<":   utils.Bool{},
		">=":  utils.Bool{},
		"!=":  utils.Bool{},
		"and": utils.Bool{},
		"or":  utils.Bool{},
	}

	gen := irgenerator.NewIRGenerator(rootTypes)
	instructions := gen.Generate(res)
	asm := asmgenerator.GenerateASM(instructions)
	output, _ = assembler.Assemble(asm, "")
	return output
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	body, err := io.ReadAll(conn)
	if err != nil {
		resp, _ := json.Marshal(map[string]string{"error": err.Error()})
		conn.Write(resp)
		return
	}

	var inputMap map[string]interface{}
	if err := json.Unmarshal(body, &inputMap); err != nil {
		resp, _ := json.Marshal(map[string]string{"error": err.Error()})
		conn.Write(resp)
		return
	}

	cmd, _ := inputMap["command"].(string)
	code, _ := inputMap["code"].(string)
	result := map[string]string{}

	switch cmd {
	case "compile":
		executable := callCompiler(code, "")
		if strings.HasPrefix(string(executable), "compiler error:") || len(executable) == 0 {
			resp, _ := json.Marshal(map[string]string{"error": string(executable)})
			conn.Write(resp)
			return
		} else {
			result["program"] = base64.StdEncoding.EncodeToString(executable)
		}
	case "ping":
	default:
		result["error"] = fmt.Sprintf("Unknown command: %s", cmd)
	}

	resp, err := json.Marshal(result)
	if err == nil {
		conn.Write(resp)
	}
}

func runServer(host string, port int) {
	address := fmt.Sprintf("%s:%d", host, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()
	fmt.Println("TCP server running on:", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		go handleConnection(conn)
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
		asm := callCompiler("read_int()", inputFile)
		os.WriteFile(outputFile, []byte(asm), 0644)
	} else if command == "serve" {
		runServer(host, port)
	} else {
		fmt.Fprintln(os.Stderr, "Error: Unknown command")
	}
}
