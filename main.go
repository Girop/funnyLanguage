package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)


func FormatedError(msg string, line int, position int) error {
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
    res,err := regexp.MatchString("[A-z]", char)
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
    panic(FormatedError(errMsg, i.Line, i.Column))
}

func (i *InputStream) Tokenize() []*Token {
    tokens := make([]*Token, 0)
    for !i.isEof() {
        newToken := i.tokenizeNext()
        fmt.Println(newToken)
        tokens = append(tokens, &newToken)
    }
    return tokens
}

func main(){
    fileData, err := GetFile()
    if err != nil {
        panic(err)
    } 
    
    stream := InputStream{fileData, 0, 1, 0}
    tokens := stream.Tokenize()

    for _, token := range tokens {
        fmt.Println(token)
    }

}
