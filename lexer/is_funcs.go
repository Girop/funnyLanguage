package lexer

import (
	"regexp"
	"strings"
)

func isChar(char string) bool {
	res, err := regexp.MatchString("[A-z_]", char)
	if err != nil {
		return false
	}
	return res
}

func isDigit(char string) bool {
    val, _ := regexp.MatchString("[0-9]", char)
    return val
}

func isOpChar(char string) bool {
	return strings.Contains("+-*/%=&|<>!", char)
}

func isPunct(char string) bool {
	return strings.Contains(",(){}[]", char)
}

func isCommentStart(char string) bool {
	return char == "#"
}

func isStringStart(char string) bool {
	return char == "\""
}

func canInsertSemicolon(char  string) bool {
    return !strings.Contains("{(;,",char)
}

var keywords = [...]string{
	"fn",
	"if",
	"else",
    "return",
}

func isKeyword(word string) bool {
	for _, val := range keywords {
		if word == val {
			return true
		}
	}
	return false
}

var booleans = [2]string{
	"true",
	"false",
}

func isBoolean(word string) bool {
	for _, val := range booleans {
		if val == word {
			return true
		}
	}
	return false
}

var annotations = [...]string{
	"number",
	"string",
	"any",
}

func isTypeAnnotation(word string) bool {
	for _, val := range annotations {
		if val == word {
			return true
		}
	}
	return false
}
