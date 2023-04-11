package main

import (
	"fmt"
	"language/lexer"
	// "language/parser"
	"os"
)

func GetFile() (string, error) {
	if len(os.Args) < 2 {
		return "", fmt.Errorf("No file specified")
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		return "", fmt.Errorf("Can't read file: %s", os.Args[1])
	}

	return string(data), nil
}

func main() {
	fileData, err := GetFile()
	if err != nil {
        fmt.Println("File not specified")
        os.Exit(1)
	}

	tokens := lexer.Tokenize(fileData)

    for _, token := range tokens {
        fmt.Println(token)
    }

	// ast := parser.NewParser(tokens).Parse()
	// fmt.Println(ast)
}
