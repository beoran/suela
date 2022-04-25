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
	ast := AstFromToken(AstTypeCall, *tok)

	tok = p.Get()
	if tok.TokenKind != TokenKindOp {
		return p.Error("Expected (, got %s", tok.String())
	}

	tok = p.Peek()
	if tok.TokenKind == TokenKindCp {
		p.Get()    // skip )
		return ast // Empty list case.
	}

	for {
		arg := p.ParseArg() // get argument
		ast.Children = append(ast.Children, arg)
		if arg.AstType == AstTypeError {
			return ast
		}

		tok = p.Peek() // get ) or ,
		if tok.TokenKind == TokenKindCp {
			p.Get()    // skip )
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
	tok := p.Peek()
	switch tok.TokenKind {
	case TokenKindField, TokenKindLiteral:
		tok = p.Get() // make progress
		return AstFromToken(AstTypeArg, *tok)
	default:
		sub := p.ParseExpr()
		if sub == nil {
			return p.Error("Expected a field, string, int, float, json, or expression. Got:%s.", tok)
		}
		return &Ast{AstTypeArg, *tok, []*Ast{sub}}
	}
	return nil
}

func (p *Parser) Error(msg string, args ...interface{}) *Ast {
	loc := fmt.Sprintf("%s:%d:%d:", p.Lexer.Source, p.Lexer.Line, p.Lexer.Column)
	err := fmt.Errorf(loc+msg, args...)
	tok := p.Lexer.MakeToken(TokenKindError, Error{err})
	return &Ast{AstTypeError, *tok, nil}
}

func (p *Parser) ParseStatement() *Ast {
	tok := p.Peek()
	switch tok.TokenKind {
	case TokenKindComment:
		comm := p.Get()
		return AstFromToken(AstTypeComment, *comm)
	case TokenKindFunc:
		return p.ParseExpr()
	default:
		return p.Error("Expected a comment, or a function call expression, found: %s.", tok)
	}
}

func (p *Parser) Parse() *Ast {
	tok := p.Lexer.MakeToken(TokenKindLiteral, Nil{})
	ast := AstFromToken(AstTypeScript, *tok)

	for sub := p.ParseStatement(); sub != nil; sub = p.ParseStatement() {
		ast.Children = append(ast.Children, sub)
		if sub.AstType == AstTypeError {
			return ast
		}
		tok := p.Get()
		switch tok.TokenKind {
		case TokenKindEol:
			continue
		case TokenKindEnd:
			return ast
		default:
			return p.Error("Expected new line, found: %s (%c).", tok, tok.TokenKind)
		}
	}

	return ast
}
