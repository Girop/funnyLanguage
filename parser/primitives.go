package parser

import (
	"language/lexer"
)

type Primitve int

const (
	Number Primitve = iota
	String
	Boolean
)

func convertTokenTypeToPrimitve(value lexer.TokenType) Primitve {
	switch value {
	case lexer.NUMBER:
		return Number
	case lexer.STRING:
		return String
	default: // TODO: fix that
		return Boolean
	}
}
