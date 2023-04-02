package main

import (
	"fmt"
	"os"
    "unicode"
    "strings"
    "strconv"
)


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

func (i *InputStream) readWord(ending string) string {
    word := ""
    for !strings.Contains(i.Peek(), ending) {
        word += i.GetNext()
    }
    return word
}

func (i *InputStream) SkipComment() {
    for i.CurrentChar() != "\n" || i.CurrentChar() == "#" {
        i.position++
    }
}

func Tokenize(stream InputStream) InputStream{
    stream.skipWhitespaces()
    switch stream.Peek() {
    case "#":
        stream.SkipComment()
        return Tokenize(stream)
        // TODO: Finish there

    }
}

func main(){
    fileData, err := GetFile()
    if err != nil {
        panic(err)
    } 

    stream := new(InputStream)
    stream.chars = fileData

}
