package main

import (
	"fmt"
	"language/lexer"
	"language/parser"
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
		panic(err)
	}

	tokens := lexer.NewStream(fileData).Tokenize()
	ast := parser.NewParser(tokens).Parse()
	fmt.Println(ast)
}
