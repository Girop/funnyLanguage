package lexer

import (
	"regexp"
)

var digitRe = regexp.MustCompile("[0-9]")
var opCharRe = regexp.MustCompile(`[+\-*/%=&|<>!]`)
var punctRe = regexp.MustCompile(`[,\(\)\{\}\[\]]`)
var semicolonRe = regexp.MustCompile(`[{\(;]`)
var charRe = regexp.MustCompile("[A-z_]")

func isChar(char string) bool {
	return charRe.MatchString(char)
}

func isDigit(char string) bool {
	return digitRe.MatchString(char)
}

func isOpChar(char string) bool {
	return opCharRe.MatchString(char)
}

func isPunct(char string) bool {
	return punctRe.MatchString(char)
}

func isCommentStart(char string) bool {
	return char == "#"
}

func isStringStart(char string) bool {
	return char == "\""
}

func canInsertSemicolon(char string) bool {
	return !semicolonRe.MatchString(char)
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
