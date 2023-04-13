package parser
import (
	"fmt"
	"language/lexer"
)

// TODO: make sure only one field exists / delete
// TODO: swap somehow for TokenTypes from lexer
type PrimitiveUnion struct {
	string  string
	int     int
	boolean bool
}

type Expresion interface {
	evaluate() PrimitiveUnion
}

type Statement interface {
	evaluate()
}

type BinaryExpr struct {
	operator *lexer.Token
	left     *Expresion
	right    *Expresion
}

type UnaryExpre struct {
	operator *lexer.Token
	body     *Expresion
}

// Needed?
type LiteralExpr struct {
	value *lexer.Token
}

type FunctionExpr struct {
	name string
	body []*Expresion
}

type Parser struct {
	tokens   []*lexer.Token
	position int
	errors   []error
}

func (p *Parser) peek() *lexer.Token {
	if p.isEOT() {
		return nil
	}
	return p.tokens[p.position]
}

func (p *Parser) getNext() *lexer.Token {
	if p.isEOT() {
		return nil
	}
	p.position++
    return p.peekPrev()
}

func (p *Parser) peekPrev() *lexer.Token {
    if p.position < 1 || p.isEOT() {
        return nil
    }
    return p.tokens[p.position - 1]
}

func (p *Parser) match(types ...lexer.TokenType) bool {
    for _, tokenType := range types {
        if tokenType == p.peek().Type {
            p.getNext()
            return true
        }
    }
    return false
}

func (p *Parser) matchSequence(types ...lexer.TokenType) bool {
    for _, nextType := range types {
        if p.peek().Type == nextType {
            p.getNext()
        } else {
            return false
        }
    }
    return true
}

var binaryOperators = [...]string {
    ">", ">=", "<", "<=", "==", "!=", "+", "-", "&&", "||", "*", "**",
}

func (p *Parser) function() *FunctionExpr {
    newFunc := new(FunctionExpr)
    if p.matchSequence(lexer.FUNC_DECLARATION, lexer.IDENTYFIER, lexer.PUNC) {

    }
    p.errors = append(p.errors, unexpectedToken(p.peekPrev()))
    return nil
}

func (p *Parser) isEOT() bool {
	return p.tokens[p.position].Type == lexer.LAST_TOKEN
}

func unexpectedToken(token *lexer.Token) error {
	return fmt.Errorf("Parsing error: unexpected token %s at line %d, position %d", token.Value, token.Line, token.ColumnPos)
}

func (p *Parser) parseNext() {
    switch p.peek().Type {
    case lexer.FUNC_DECLARATION:
        p.function()
    }
}

func NewParser(tokens []*lexer.Token) *Parser {
	return &Parser{
		tokens,
		0,
		make([]error, 0),
	}
}
