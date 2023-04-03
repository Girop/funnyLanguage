package main

import (
	"fmt"
	"os"
)

func GetFile() (string, error) {
	if len(os.Args) < 2 {
		return "", fmt.Errorf("No file specified")
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		return "", fmt.Errorf("Can't read file: %s", os.Args[1])
	}

	return string(data), nil
}

// Ast
func FormatedUnexpectedTokenError(name string) error {
	return fmt.Errorf("Parser error, unexpected token: ", name)
}

func getPrecedence(op Token) int {
	precedence := map[string]int{
		"=":  1,
		"||": 2, "&&": 2,
		"<": 3, "<=": 3, ">": 3, ">=": 3,
		"+": 4, "-": 4,
		"*": 5, "/": 5, "%": 5,
	}
	if op.Type != OPERATOR {
		return 0
	}
	return precedence[op.Value]
}

type AstType int

const (
	FUNCTION AstType = iota
)

type AstElement interface {
	getType() AstType
}

type FunctionAst struct {
	Name string
	Args *AstElement
	body *AstElement
}

func (f FunctionAst) getType() AstType {
	return FUNCTION
}

// TODO: Do i need to operate on position ?
// Can't i just pass next tokens as parameters ?
type Parser struct {
	tokens   []*Token
	position int
}

func (p Parser) peek() (*Token, error) {
	if p.endOfTokens() {
		return nil, fmt.Errorf("EOT")
	}
	return p.tokens[p.position+1], nil
}

func (p *Parser) getNext() *Token {
	p.position++
	return p.tokens[p.position+1]
}

func (p *Parser) endOfTokens() bool {
	return p.position >= len(p.tokens)
}

// TODO:
func (p *Parser) parseArgs() {

}

func (p *Parser) parseDelimiters(start Token, end Token) (*AstElement, error) {
	token := p.getNext()
	if *token != start {
		return nil, FormatedUnexpectedTokenError(token.Value)
	}
	// TODO: Do i need one func for parsing args and one for body?
	// Or use this as subprocedure of both, addontional struct for
	// holding arg data, and still somwhow implement those interfaces
	body := new()
	for !p.endOfTokens() {
		nextTok, err := p.peek()
		if err != nil {
			return err
		}
		if *nextTok == end {
			return
		}
	}
}

func (p *Parser) parseFunction() (*FunctionAst, error) {
	funcNode := new(FunctionAst)
	funcName := p.getNext()

	if funcName.Type != IDENTYFIER {
		return nil, FormatedUnexpectedTokenError(funcName.Value)
	}
	funcNode.Name = funcName.Value
	opening, ending := Token{PUNC, "("}, Token{PUNC, ")"}
	args, err := p.parseDelimiters(opening, ending)
	if err != nil {
		return nil, err
	}
	funcNode.Args = args

	opening, ending = Token{PUNC, "{"}, Token{PUNC, "}"}
	body, err := p.parseDelimiters(opening, ending)
	if err != nil {
		return nil, err
	}
	funcNode.body = body

	return funcNode, nil
}

func isMain(token1 *Token, token2 *Token) bool {
	return token1.Type == KEYWORD &&
		token1.Value == "fn" &&
		token2.Type == IDENTYFIER &&
		token2.Value == "main"
}

func (p *Parser) parseEntryPoint() (*AstElement, error) {
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

func (p *Parser) Parse() *AstElement {
	topNode, err := p.parseEntryPoint()
	if err != nil {
		panic(err)
	}
	return topNode
}

func main() {
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
