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

	// Words & keywords
	IDENTYFIER
	FUNC_DECLARATION
	CONTROL_FLOW
	TYPE_ANNOTATION
	OTHER_KEYWORD

	// Punctuation
	PUNC
	END_OF_LINE
	OPERATOR
	LAST_TOKEN
)

// TODO fix Column positioning
// TODO write some tests

type Token struct {
	Type      TokenType
	Value     string
	Line      int
	ColumnPos int
}

func formatedPositonedError(msg string, stream InputStream) error {
	return fmt.Errorf("Error: %s, at line %d, position %d", msg, stream.line, stream.position)
}

type InputStream struct {
	chars    string
	position int
	line     int
	column   int
	// Anywhere error occurs, if is serious it should be explicitly
	// set as 'terminating' by appending to this slice, allowing program
	// to find another errors but, preventing from entering next phase of parsing
	errors []error
}

func (i *InputStream) getNext() string {
	i.position++
	if i.isEof() {
		return ""
	}
	return string(i.chars[i.position])
}

func (i InputStream) peek() string {
	if i.isEof() {
		return ""
	}
	return string(i.chars[i.position+1])
}

func (i *InputStream) skipWhitespaces() {
	for {
		char := i.peek()

		switch {
		case strings.ContainsAny(char, " \t\r\f\v"):
			i.position++
			i.column++
		case char == "\n":
			i.position++
			i.line++
			i.column = 1
		default:
			return
		}
	}
}

func (i *InputStream) readNumber() string {
	numberString := string(i.getNext())
	sepratorOccured := false

	for {
		nextChar := i.peek()
		switch {
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
		nextChar := i.getNext()
		if strings.Contains(nextChar, delimiter) {
			break
		}
		word += nextChar
	}

	return word
}

func (i *InputStream) skipLine() {
	for nextChar := i.peek(); nextChar != "\n" && nextChar != ""; nextChar = i.peek() {
		i.position++
	}
}

func (i *InputStream) handleOpChar() *Token {
	nextChar := i.combainOpChar()
	if len(nextChar) == 1 && nextChar == "<" && i.peekAfterWord() == ">" {
		return i.newToken(TYPE_ANNOTATION, i.getNext()+i.readToDelimiter(">"))
	}
	return i.newToken(OPERATOR, nextChar)
}

func (i InputStream) peekAfterWord() string {
	i.getNext()
	i.parseWord()
	return i.getNext()
}

var combinations = []string{"<=", ">=", "!=", "==", "||", "&&"}

func (i *InputStream) combainOpChar() string {
	current := i.getNext()
	comb := current + i.peek()

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

	for char := i.peek(); isChar(char); char = i.peek() {
		word += i.getNext()
	}
	return word
}

func (i *InputStream) tokenizeWord() *Token {
	word := i.parseWord()
	var type_ TokenType
	switch {
	case isControlFlow(word):
		type_ = CONTROL_FLOW
	case isFuncDecl(word):
		type_ = FUNC_DECLARATION
	case isBoolean(word):
		type_ = BOOLEAN
	default:
		type_ = IDENTYFIER
	case isKeyword(word):
		type_ = OTHER_KEYWORD
	}
	return i.newToken(type_, word)
}

func (i InputStream) newToken(type_ TokenType, value string) *Token {
	return &Token{
		type_,
		value,
		i.line,
		i.column,
	}
}

func (i *InputStream) tokenizeNext() *Token {
	i.skipWhitespaces()
	char := i.peek()

	switch {
	case isCommentStart(char):
		i.skipLine()
		return i.tokenizeNext()
	case isDigit(char):
		return i.newToken(NUMBER, i.readNumber())
	case isPunct(char):
		return i.newToken(PUNC, i.getNext())
	case isOpChar(char):
		return i.handleOpChar()
	case isStringStart(char):
		i.getNext() // Skipping "
		return i.newToken(STRING, i.readToDelimiter("\""))
	case isChar(char):
		return i.tokenizeWord()
	}

	if !i.isEof() {
		err := formatedPositonedError(fmt.Sprint("Unexpected character: ", char), *i)
		i.errors = append(i.errors, err)
	}
	return i.newToken(LAST_TOKEN, "END")
}

func (i *InputStream) Tokenize() []*Token {
	tokens := make([]*Token, 0)
	prevToken := i.newToken(END_OF_LINE, ";")

	for !i.isEof() {
		newToken := i.tokenizeNext()
		if newToken.Line != prevToken.Line && canInsertSemicolon(prevToken.Value) {
			tokens = append(tokens, i.newToken(END_OF_LINE, ";"))
		}

		prevToken = newToken
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
}

func NewStream(file string) *InputStream {
	stream := new(InputStream)
	stream.chars = file
	stream.line = 1
	stream.column = 1
	return stream
}

func Tokenize(file string) []*Token {
	return NewStream(file).Tokenize()
}
