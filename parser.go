package suela

import "fmt"

/*
SCRIPT -> STATEMENTS .
STATEMENTS -> STATEMENT nl STATEMENTS | .
STATEMENT -> comment | EXPR .
EXPR -> funcname op ARGS cp .
ARGS -> ARG comma ARGS | .
ARG -> string | int | float | field | json | CALL .
*/

type Parser struct {
	Lexer
	Lookahead *Token
}

func (p *Parser) Peek() *Token {
	if p.Lookahead == nil {
		p.Lookahead = p.Lexer.Lex()
	}
	return p.Lookahead
}

func (p *Parser) Get() *Token {
	if p.Lookahead == nil {
		return p.Lexer.Lex()
	} else {
		res := p.Lookahead
		p.Lookahead = nil
		return res
	}
}

// ParseFunc return nil in case of "no match", an an ast with type
// AstTypeError in case of errors.
type ParseFunc func() *Ast

func AstFromToken(typ AstType, tok *Token, sub ...*Ast) *Ast {
	return &Ast{typ, tok, sub}
}

func (p *Parser) AcceptToken(kind TokenKind, typ AstType) ParseFunc {
	return func() *Ast {
		tok := p.Peek()
		if tok.TokenKind != kind {
			return nil
		}
		return AstFromToken(typ, tok)
	}
}

func (p *Parser) RequireToken(kind TokenKind, typ AstType) ParseFunc {
	return func() *Ast {
		tok := p.Peek()
		if tok.TokenKind != kind {
			return nil
		}
		return AstFromToken(typ, tok)
	}
}

func parseAlternates(alternates ...ParseFunc) ParseFunc {
	return func() *Ast {
		for _, alt := range alternates {
			subAst := alt()
			if subAst != nil {
				return subAst
			}
		}
		return nil
	}
}

func parseSequence(typ AstType, sequence ...ParseFunc) ParseFunc {
	return func() *Ast {
		ast := &Ast{AstType: typ}
		for _, seq := range sequence {
			subAst := seq()
			if subAst == nil {
				return nil
			}
			ast.Children = append(ast.Children, subAst)
		}
		return ast
	}
}

func parseList(typ AstType, start, item, sep, end ParseFunc) ParseFunc {
	return func() *Ast {
		ast := &Ast{AstType: typ}
		var subAst *Ast
		if start != nil {
			subAst = start()
			if subAst == nil {
				return nil
			}
		}
		for {
			subAst = item()
			if subAst == nil {
				return ast
			}
			ast.Children = append(ast.Children, subAst)
			if sep != nil {
				subAst = sep()
				if subAst == nil {
					return nil
				}
			}
			if end != nil {
				subAst = end()
				if subAst != nil {
					return nil
				}
			}
		}
		return ast
	}
}

func NewParserFromString(source, input string) *Parser {
	parser := &Parser{}
	parser.Lexer = MakeLexerFromString(source, input)
	parser.Lookahead = nil
	return parser
}

func (p *Parser) ParseExpr() *Ast {
	tok := p.Get()
	if tok.TokenKind != TokenKindFunc {
		return p.Error("Expected @function name, got %s", tok.String())
	}
	ast := &Ast{AstTypeCall, tok, nil}

	tok = p.Get()
	if tok.TokenKind != TokenKindOp {
		return p.Error("Expected (, got %s", tok.String())
	}

	for {
		tok = p.Peek()
		if tok.TokenKind == TokenKindCp {
			return ast // end of call arguments.
		}
		arg := p.ParseArg()
		if arg.AstType == AstTypeError {
			ast.Children = append(ast.Children, arg)
			return ast
		}
		tok = p.Peek()
		if tok.TokenKind == TokenKindCp {
			p.Get()
			return ast // end of call arguments.
		} else if tok.TokenKind == TokenKindComma {
			p.Get() // skip comma
		} else {
			return p.Error("Expected , or ) got %s", tok.String())
		}
	}
	return nil
}

func (p *Parser) ParseArg() *Ast {
	tok := p.Get()
	switch tok.TokenKind {
	case TokenKindField, TokenKindLiteral:
		return &Ast{AstTypeArg, tok, nil}
	default:
		sub := p.ParseExpr()
		if sub == nil {
			return p.Error("Expected a field, string, int, float, json, or expression. Got:%s.", tok)
		}
		return &Ast{AstTypeArg, nil, []*Ast{sub}}
	}
	return nil
}

func (p *Parser) Error(msg string, args ...interface{}) *Ast {
	loc := fmt.Sprintf("%s:%d:%d:", p.Lexer.Source, p.Lexer.Line, p.Lexer.Column)
	err := fmt.Errorf(loc+msg, args...)
	tok := p.Lexer.MakeToken(TokenKindError, Error{err})
	return &Ast{AstTypeError, tok, nil}
}

func (p *Parser) ParseStatement() *Ast {
	tok := p.Peek()
	switch tok.TokenKind {
	case TokenKindComment:
		return &Ast{AstTypeComment, p.Get(), nil}
	case TokenKindFunc:
		return p.ParseExpr()
	default:
		return p.Error("Expected a comment, or a function call expression Got:%s.", tok)
	}
}

func (p *Parser) Parse() *Ast {
	ast := &Ast{AstTypeScript, nil, nil}

	for sub := p.ParseStatement(); sub != nil; sub = p.ParseStatement() {
		ast.Children = append(ast.Children, sub)
		tok := p.Get()
		switch tok.TokenKind {
		case TokenKindEol:
			continue
		case TokenKindEnd:
			break
		default:
			return p.Error("Expected new line, found: s.", tok.Data.String())
		}
	}

	return ast
}
