package suela

type AstType int

const (
	AstTypeNone AstType = iota
	AstTypeError
	AstTypeArg
	AstTypeStatement
	AstTypeScript
	AstTypeCall
	AstTypeComment
)

type Ast struct {
	AstType
	*Token
	Children []*Ast
}
