package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
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

func (i *InputStream) CurrentChar() string{
    return string(i.chars[i.position])
}

func (i *InputStream) GetNext() string {
    i.position++
    return i.CurrentChar()
}

func (i InputStream) Peek() string {
    return string(i.chars[i.position + 1])
}

func (i *InputStream) skipWhitespaces() {
    for unicode.IsSpace(rune(i.chars[i.position])) {
        i.position++
    }
}

func (i *InputStream) readNumber() int{
    numberString := ""
    for isDigit(i.Peek()) {
        numberString += i.GetNext()
    }

    numberValue, err := strconv.Atoi(numberString)
    if err != nil {
        panic(fmt.Errorf("Couln't parse number: line %d position %d", i.Line, i.Column))
    }
    return numberValue
}

func (i *InputStream) ReadToDelimiter(delimiter string) string {
    word := ""
    for !strings.Contains(i.Peek(), delimiter) {
        word += i.GetNext()
    }
    return word
}

func (i *InputStream) SkipComment() {
    for i.CurrentChar() != "\n" || i.CurrentChar() == "#" {
        i.position++
    }
}

func (i *InputStream) combainOpChar() string{
    //+-*/%=&|<>!
    combinations := []string{"<=", ">=", "!=", "||", "&&"}
    current := i.GetNext()
    comb := current + i.Peek()
    
    for _, possible := range combinations {
        if comb == possible {
            return comb
        }
    }
    return current
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

// Create Many tokens and token interfaces
type Token struct {
    class TokenType
    value string
}

func isWord(char string) bool{
    res,err := regexp.MatchString("w+", char)
    if err != nil {
        panic("Erro during matching string")
    }
    return res
}

func (i *InputStream) TokenizeNext() Token {
    i.skipWhitespaces()
    char := i.Peek()
    switch {
    case char == "#":
        i.SkipComment()
        return i.TokenizeNext()
    case isDigit(char): 
        return Token{NUMBER, i.readNumber()}
    case char == "\"":
        return Token{STRING, i.ReadToDelimiter("\"")}
    case isPunct(char):
        return Token{PUNC, i.GetNext()}
    case isOpChar(char):
        return Token{OPERATOR, i.combainOpChar()}
    case isWord(char):
        word := ""
        for {
            if !isWord(i.Peek()) { // Check for eof
                break
            }
            word += i.GetNext()
        }
        token := Token{IDENTYFIER, word}
        if isKeyword(word) {
            token.class = KEYWORD
        }
        return token
    }
    errMsg := fmt.Sprintf("Unexpected character %s", i.Peek())
    panic(FormatedError(errMsg, i.Line, i.Column))
}

func main(){
    fileData, err := GetFile()
    if err != nil {
        panic(err)
    } 
    // TODO check eof
    // TODO different Token types / interfaces

    stream := new(InputStream)
    stream.chars = fileData

}
