package suela

import "io"

const LEXER_EOF = -1
const LEXER_ERROR = -2

type Position struct {
	Line   int
	Column int
}

type Token struct {
	Position
	Data
}

type Lexer struct {
	Position
	io.RuneScanner
	prevLineSize int
	Error        error
}

func (l Lexer) Peek() rune {
	r, _, err := l.RuneScanner.ReadRune()
	if err != nil {
		if err == io.EOF {
			return LEXER_EOF
		} else {
			l.Error = err
			return LEXER_ERROR
		}
	}
	l.RuneScanner.UnreadRune()
	return r
}

func (l Lexer) Get() rune {
	r, _, err := l.RuneScanner.ReadRune()
	if err != nil {
		if err == io.EOF {
			return LEXER_EOF
		} else {
			l.Error = err
			return LEXER_ERROR
		}
	}
	if r == '\n' {
		l.Line++
		l.Column = 1
	} else {
		l.Column++
	}
	return r
}

func (l Lexer) LexString() Token {

}

func (l Lexer) Lex() *Token {
	r := l.Peek()
	switch r {
	case LEXER_EOF:
		return nil
	}

	return &Token{}
}
