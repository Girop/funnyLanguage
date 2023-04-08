package lexer

import (
	"fmt"
	"strings"
)

type TokenType int

const (
	// Primitves
	NUMBER TokenType = iota
	STRING
	BOOLEAN

	// Words
	IDENTYFIER
	KEYWORD

	// Other
	PUNC
	OPERATOR
	TYPE_ANNOTATION
)

type Token struct {
	Type  TokenType
	Value string
}

func formatedPositonedError(msg string, stream InputStream) error {
	return fmt.Errorf("Error: %s, at line %d, position %d", msg, stream.line, stream.position)
}

type InputStream struct {
	chars    string
	position int
	line     int
	column   int
    // Note: Anywhere error occurs, if is serious it should be explicitly
    // set as 'terminating' by appending to this slice, allowing program 
    // to find another errors but, preventing from entering next phase of parsing
	errors []error
}

// Important: use only if sure, that next token exists FIXME
func (i *InputStream) getNext() string {
    i.position++
	return string(i.chars[i.position])
}

func (i InputStream) peek() (string, error) {
	if i.isEof() {
		err := formatedPositonedError("EOF", i)
		return "", err
	}
	return string(i.chars[i.position+1]), nil
}

func (i *InputStream) skipWhitespaces() {
	for {
		char, err := i.peek()

		switch {
		case err != nil:
			return
		case strings.ContainsAny(char, " \t\r\f\v"):
			i.position++
			i.column++
		case char == "\n":
			i.position++
			i.line++
			i.column = 0
		default:
			return
		}
	}
}

func (i *InputStream) readNumber() string {
	numberString := string(i.getNext())
	sepratorOccured := false

	for {
		nextChar, err := i.peek()
		switch {
		case err != nil:
			i.errors = append(i.errors, err)
			return ""
		case isDigit(nextChar):
			numberString += i.getNext()
		case nextChar == "." && !sepratorOccured:
			sepratorOccured = true
			numberString += i.getNext()
		case nextChar == ".":
			i.errors = append(i.errors, formatedPositonedError("Two or more floating point separtors", *i))
			fallthrough
		default:
			return numberString
		}
	}
}

func (i *InputStream) readToDelimiter(delimiter string) string {
	word := ""
	for {
		nextChar, err := i.peek()
		if err != nil || strings.Contains(nextChar, delimiter) {
			break
		}
		word += i.getNext()
	}
	return word
}

// TODO test if works this way
func (i InputStream) peekToDelimiter(delimiter string) string { 
	word := ""
	for {
		nextChar, err := i.peek()
		if err != nil || strings.Contains(nextChar, delimiter) {
			break
		}
		word += i.getNext()
	}
	return word
}

func (i *InputStream) skipLine() {
	for nextChar, err := i.peek(); nextChar != "\n" || err != nil; nextChar, err = i.peek() {
		i.position++
	}
}

func (i *InputStream) handleOpChar() *Token {
    nextChar := i.combainOpChar()
    if nextChar == "<" && isTypeAnnotation(i.peekToDelimiter(">")){
        annotation := i.getNext() + i.readToDelimiter(">")
        return &Token{TYPE_ANNOTATION, annotation}
    } 
	return &Token{OPERATOR, nextChar}
}

func (i *InputStream) combainOpChar() string {
	combinations := []string{"<=", ">=", "!=", "||", "&&"}
	current := i.getNext()

	nextChar, err := i.peek()
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
	return i.position+1 >= len(i.chars)
}

func (i *InputStream) parseWord() string {
	word := ""
	for {
		nextChar, err := i.peek()
		if err != nil || !isChar(nextChar) {
			break
		}
		word += i.getNext()
	}
	return word
}

func (i *InputStream) tokenizeWord() *Token {
	word := i.parseWord()
	var type_ TokenType
	switch {
	case isKeyword(word):
		type_ = KEYWORD
	case isBoolean(word):
		type_ = BOOLEAN
	default:
		type_ = IDENTYFIER
	}

	return &Token{type_, word}
}

func (i *InputStream) tokenizeNext() *Token {
	i.skipWhitespaces()
	char, _ := i.peek()
	switch {
	case isCommentStart(char):
		i.skipLine()
		return i.tokenizeNext()
	case isDigit(char):
		return &Token{NUMBER, i.readNumber()}
	case isPunct(char):
		return &Token{PUNC, i.getNext()}
	case isOpChar(char):
		return i.handleOpChar()
	case isStringStart(char):
		return &Token{STRING, i.readToDelimiter("\"")}
	case isChar(char):
		return i.tokenizeWord()
	}

	err := formatedPositonedError(fmt.Sprint("Unexpected character: ", char), *i)
	i.errors = append(i.errors, err)
	return nil
}

func (i *InputStream) Tokenize() []*Token {
	tokens := make([]*Token, 0)

	for !i.isEof() {
		newToken := i.tokenizeNext()
		tokens = append(tokens, newToken)
	}
	if len(i.errors) > 0 {
		i.ShowErrorMsg()
	}
	return tokens
}

func (i *InputStream) ShowErrorMsg() {
	msg := fmt.Sprintf("\nParsing errors: encounterd %d errors: \n", len(i.errors))
	for index, errMsg := range i.errors {
		msg += fmt.Sprintf("%d. %s\n", index, errMsg)
	}
	panic(msg)
}

func NewStream(file string) *InputStream {
	stream := new(InputStream)
	stream.chars = file
	stream.line = 1
	stream.line = 1
	return stream
}

func Tokenize(file string) []*Token {
    return NewStream(file).Tokenize()
}
