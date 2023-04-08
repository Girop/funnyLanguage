package parser

import (
	"fmt"
	"language/lexer"
)

func FormatedUnexpectedTokenError(token lexer.Token) error {
	return fmt.Errorf("Parser error, unexpected token: ", token.Value)
}

func getPrecedence(op lexer.Token) int {
	precedence := map[string]int{
		"=":  1,
		"||": 2, "&&": 2,
		"<": 3, "<=": 3, ">": 3, ">=": 3,
		"+": 4, "-": 4,
		"*": 5, "/": 5, "%": 5,
	}
	if op.Type != lexer.OPERATOR {
		return 10
	}
	return precedence[op.Value]
}

type Parser struct {
	tokens   []*lexer.Token
	position int
}

func (p Parser) peek() (*lexer.Token, error) {
	if p.endOfTokens() {
		return nil, fmt.Errorf("EOT")
	}
	return p.tokens[p.position+1], nil
}

func (p *Parser) getNext() (*lexer.Token, error) {
	p.position++
	if p.position >= len(p.tokens) {
		return nil, fmt.Errorf("EOT")
	}
	return p.tokens[p.position], nil
}

func (p *Parser) endOfTokens() bool {
	return p.position >= len(p.tokens)
}

func (p *Parser) parseFunction() (*FunctionAst, error) {
	funcNode := new(FunctionAst)

	funcName, err := p.getNext()
    if err != nil {
        return nil, err
    }

	if funcName.Type != lexer.IDENTYFIER {
		return nil, FormatedUnexpectedTokenError(*funcName)
	}
	funcNode.Name = funcName.Value

	args, err := p.parseArgs()
	if err != nil {
		return nil, err
	}
	funcNode.Args = args

	body, err := p.parseBody()
	if err != nil {
		return nil, err
	}
	funcNode.body = body

	return funcNode, nil
}

func (p *Parser) parseArgs() (ArgsAst, error) {
	tokensToParse, err := p.tokensInsideDelimiters(lexer.Token{lexer.PUNC, "("}, lexer.Token{lexer.PUNC, ")"})
	if err != nil {
		return nil, err
	}

	args := make(ArgsAst, 0)
	newVar := new(VarAst)
	for _, token := range tokensToParse {
		switch token.Type {
		case lexer.IDENTYFIER:
			newVar.Name = token.Value
		case lexer.KEYWORD: // TODO: better defence against other keywords
			newVar.PrimType = convertTokenTypeToPrimitve(token.Type)
			args = append(args, newVar)
		case lexer.PUNC:
			newVar = new(VarAst)
		default:
			return nil, FormatedUnexpectedTokenError(*token)
		}
	}
	return args, nil
}

func (p *Parser) parseBody() (*Prog, error) { 
    // TODO:
    // opening, ending := lexer.Token{lexer.PUNC, "{"}, lexer.Token{lexer.PUNC, "}"}
    // _, err := p.tokensInsideDelimiters(opening, ending)
    // if err != nil {
    //     return nil, err
    // }

    body := new(Prog)
    // TODO: finsh from here
    return body, nil
}

func (p *Parser) tokensInsideDelimiters(start lexer.Token, end lexer.Token) ([]*lexer.Token, error) {
	nextTok, err := p.getNext()
	if err != nil {
		return nil, err
	}

	if *nextTok != start {
		return nil, FormatedUnexpectedTokenError(*nextTok)
	}
	tokensInside := make([]*lexer.Token, 0)
	// FIXME: possible bug
	for nextTok, err := p.getNext(); *nextTok != end; nextTok, err = p.getNext() {
		if err != nil {
			return nil, fmt.Errorf("EOT")
		}
		tokensInside = append(tokensInside, nextTok)
	}
	return tokensInside, nil
}

func isMain(token1 *lexer.Token, token2 *lexer.Token) bool {
	return token1.Type == lexer.KEYWORD &&
		token1.Value == "fn" &&
		token2.Type == lexer.IDENTYFIER &&
		token2.Value == "main"
}

func (p *Parser) parseEntryPoint() (*FunctionAst, error) {
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

func (p *Parser) Parse() *FunctionAst {
	topNode, err := p.parseEntryPoint()
	if err != nil {
		panic(err)
	}
	return topNode
}

func NewParser(tokens []*lexer.Token) *Parser {
	return &Parser{
		tokens,
		0,
	}
}

