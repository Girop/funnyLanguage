package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)


func FormatedPositonedError(msg string, line int, position int) error {
    return fmt.Errorf("Error: %s, at line %d, position %d", msg, line, position)
}

func GetFile() (string, error){
    if len(os.Args) < 2 {
        return "", fmt.Errorf("No file specified")
    }
    
    data, err := os.ReadFile(os.Args[1])
    if err != nil {
        return "", fmt.Errorf("Can't read file: %s", os.Args[1])
    }
    
    return string(data), nil
}


func isKeyword(word string) bool {
    keyWords := []string{
        "if",
        "else",
        "fn",
        "true",
        "false",
    }

    for _, val := range keyWords {
        if word == val {
            return true
        }
    }
    return false
}

func isDigit(char string) bool {
    return strings.Contains("0123456789", char)
}

func isOpChar(char string) bool {
    return strings.Contains("+-*/%=&|<>!", char)
}

func isPunct(char string) bool {
    return strings.Contains(",;(){}[]", char)
}

type InputStream struct {
    chars string
    position int
    Line int
    Column int
}


func (i *InputStream) GetNext() string {
    i.position++
    return string(i.chars[i.position])
}

func (i InputStream) Peek() (string, error) {
    if i.isEof() {
        return "", fmt.Errorf("EOF")
    }
    return string(i.chars[i.position + 1]), nil
}

func (i *InputStream) skipWhitespaces() {
    for {
        char, err := i.Peek()
        if err != nil{
            return
        }
        if strings.ContainsAny(char, " \t\r\f\v") {
            i.position++
            i.Column++
        } else if char == "\n"{
            i.position++
            i.Line++
            i.Column = 0
        } else {
            return
        }
    }
}

func (i *InputStream) readNumber() int{
    numberString := ""
    for {
        nextChar, err := i.Peek()
        if err != nil {
            return 0
        }
        if isDigit(nextChar) {
            numberString += i.GetNext()
        } else {
            break
        }
    }

    numberValue, err := strconv.Atoi(numberString)
    if err != nil {
        panic(fmt.Errorf("Couln't parse number: line %d position %d", i.Line, i.Column))
    }
    return numberValue
}

func (i *InputStream) ReadToDelimiter(delimiter string) string {
    word := ""
    for {
        nextChar, err := i.Peek()
        if err != nil || strings.Contains(nextChar, delimiter){
            break
        }
        word += i.GetNext()
    }
    return word
}

func (i *InputStream) SkipLine() {
    for nextChar, err := i.Peek(); nextChar != "\n" || err != nil; nextChar, err = i.Peek(){
        i.position++
    }
}

func (i *InputStream) combainOpChar() string{
    combinations := []string{"<=", ">=", "!=", "||", "&&"}
    current := i.GetNext()
    nextChar, err := i.Peek()
    if err != nil {
        return ""
    }

    comb := current + nextChar
 
    for _, possible := range combinations {
        if comb == possible {
            return comb
        }
    }
    return current
}

func (i *InputStream) isEof() bool {
    return i.position + 1 >= len(i.chars)
}

type TokenType int

const (
    NUMBER TokenType = iota
    STRING
    IDENTYFIER
    KEYWORD
    PUNC
    OPERATOR
)

type Token struct {
    Type TokenType
    Value string
}

func isWord(char string) bool{
    res,err := regexp.MatchString("[A-z]_", char)
    if err != nil {
        panic("Error during matching string")
    }
    return res
}

func (i *InputStream) tokenizeWord() Token{ 
        word := ""
        for {
            nextChar, err := i.Peek()
            if err != nil || !isWord(nextChar) {
                break
            }
            word += i.GetNext()
        }
        type_ := IDENTYFIER
        if isKeyword(word) {
            type_ = KEYWORD
        }
    return Token{type_, word}
}

func (i *InputStream) tokenizeNext() Token {
    i.skipWhitespaces()
    char, _ := i.Peek()
    switch {
    case char == "#":
        i.SkipLine()
        return i.tokenizeNext()
    case isDigit(char): 
        return Token{NUMBER, fmt.Sprintf("%v", i.readNumber())}
    case char == "\"":
        return Token{STRING, i.ReadToDelimiter("\"")}
    case isPunct(char):
        return Token{PUNC, i.GetNext()}
    case isOpChar(char):
        return Token{OPERATOR, i.combainOpChar()}
    case isWord(char):
        return i.tokenizeWord()
    }

    errMsg := fmt.Sprintf("Unexpected character %s", char)
    panic(FormatedPositonedError(errMsg, i.Line, i.Column))
}

func (i *InputStream) Tokenize() []*Token {
    tokens := make([]*Token, 0)
    for !i.isEof() {
        newToken := i.tokenizeNext()
        tokens = append(tokens, &newToken)
    }
    return tokens
}

// Ast
func getPrecedence(op Token) int {
    precedence := map[string]int {
        "=": 1,
        "||": 2, "&&": 2,
        "<": 3, "<=":3, ">": 3, ">=": 3,
        "+": 4, "-": 4,
        "*": 5, "/": 5, "%": 5,
    }
    if op.Type != OPERATOR {
        return 0
    }
    return precedence[op.Value]
}

func FormatedUnexpectedTokenError(name string) error {
    return fmt.Errorf("Parser error, unexpected token: ", name)
}

type AstNode interface{
}

type AstToken struct {
    Type string
}

type FunctionAst struct {
    Args []Token
    body *AstNode
}

type Parser struct {
    tokens []*Token
    position int
}

func (p Parser) peek() (*Token, error) {
    if p.endOfTokens() {
        return nil, fmt.Errorf("EOT")
    }
    return p.tokens[p.position + 1], nil
}

func (p *Parser) getNext() *Token {
    p.position++
    return p.tokens[p.position + 1]
}

func (p *Parser) endOfTokens() bool{
    if p.position >= len(p.tokens) {
        return true
    }
    return false
}


func parseDelimiters(start Token, end Token) {

}

func (p *Parser) parseFunction() (*AstNode, error){
    newFunc := new(FunctionAst)
    funcName := p.getNext()
    if funcName.Type != IDENTYFIER {
        return nil, FormatedUnexpectedTokenError(funcName.Value)
    }
    return newFunc, nil
}

func isMain(token1 *Token, token2 *Token) bool{
    return token1.Type == KEYWORD && 
    token1.Value == "fn" && 
    token2.Type == IDENTYFIER && 
    token2.Value == "main"
}

func (p *Parser) parseEntryPoint() (*AstNode, error){
    for i, token := range p.tokens {
        p.position = i

        nextTok, err := p.peek()
        if err != nil {
            return nil, fmt.Errorf("No entry point (fn main{}) found")
        }
        if isMain(token, nextTok) {
            mainFunc, err := p.parseFunction()
            if err != nil {
                return nil, err
            }
            return mainFunc, nil
        }
    }
    return nil, fmt.Errorf("No idea what happend, parseEntryPoint func")
}

func (p *Parser) Parse() *AstNode{
    topNode, err := p.parseEntryPoint()
    if err != nil {
        panic(err)
    }
    return topNode
}

func main(){
    fileData, err := GetFile()
    if err != nil {
        panic(err)
    } 
    
    stream := InputStream{fileData, 0, 1, 0}
    tokens := stream.Tokenize()
    parser := Parser{tokens, 0}
    ast := parser.Parse()
    fmt.Println(ast)
}
