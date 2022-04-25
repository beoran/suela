package suela

import "fmt"

type AstType rune

const (
	AstTypeNone      AstType = 0
	AstTypeError     AstType = 'E'
	AstTypeArg       AstType = 'A'
	AstTypeStatement AstType = 'S'
	AstTypeScript    AstType = 'T'
	AstTypeCall      AstType = 'C'
	AstTypeComment   AstType = '#'
)

type Ast struct {
	AstType
	Token
	Children []*Ast
}

func AstFromToken(typ AstType, tok Token, sub ...*Ast) *Ast {
	return &Ast{typ, tok, sub}
}

// Describe returns a string that describes the AST as an S-expression.
func (a Ast) Describe() string {
	res := fmt.Sprintf("(%c (%s)", a.AstType, a.Token.Data.String())
	for _, c := range a.Children {
		res += " "
		res += c.Describe()
	}
	res += ")"
	return res
}

func (a *Ast) RunArg(s *Suela, data ...Data) Data {
	if len(a.Children) > 0 {
		var res Data = Nil{}
		for _, child := range a.Children {
			res = child.Run(s, data...)
		}
		return res
	}
	return a.Token.Data
}

func (a *Ast) RunStatement(s *Suela, data ...Data) Data {
	var res Data = Nil{}

	for _, child := range a.Children {
		res = child.Run(s, data...)
	}
	return res
}

func (a *Ast) RunCall(s *Suela, data ...Data) Data {
	args := []Data{}
	var sub Data = Nil{}
	for _, child := range a.Children {
		sub = child.Run(s, data...)
		args = append(args, sub)
	}
	fun, ok := s.Funcs[a.Token.Data.String()]
	if !ok {
		return Nil{}
	}
	return fun.Impl(s, args...)
}

func (a *Ast) RunScript(s *Suela, data ...Data) Data {
	var res Data = Nil{}
	for _, child := range a.Children {
		res = child.Run(s, data...)
	}
	return res
}

func (a *Ast) Run(s *Suela, data ...Data) Data {
	switch a.AstType {
	case AstTypeError, AstTypeComment:
		return a.Token.Data
	case AstTypeArg:
		return a.RunArg(s, data...)
	case AstTypeStatement:
		return a.RunStatement(s, data...)
	case AstTypeScript:
		return a.RunScript(s, data...)
	case AstTypeCall:
		return a.RunCall(s, data...)
	default:
		return Nil{}
	}
}
