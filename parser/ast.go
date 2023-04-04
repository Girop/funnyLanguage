package parser

type AstType int

const (
	FUNCTION AstType = iota
	VAR
	ASSIGNMENT
	CALL
	ARGS
	PROG
)

type AstElement interface {
	getType() AstType
}

type FunctionAst struct {
	Name string
	Args ArgsAst
	body *Prog
}

func (f FunctionAst) getType() AstType {
	return FUNCTION
}

type ArgsAst []*VarAst

type VarAst struct {
	Name     string
	PrimType Primitve
}

func (a *VarAst) getType() AstType {
	return VAR
}

type CallAst struct {
	Name string
	Args ArgsAst
}

func (c *CallAst) getType() AstType {
	return CALL
}

type Prog struct {
	NodesAst []*AstElement
}

func (p *Prog) getType() AstType {
	return PROG
}
