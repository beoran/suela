package suela

import "io"
import "fmt"
import "strings"
import "strconv"

const LEXER_EOF = -1
const LEXER_ERROR = -2

type Position struct {
	Source string
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
	ReadError    error
}

func (l Lexer) Peek() rune {
	r, _, err := l.RuneScanner.ReadRune()
	if err != nil {
		if err == io.EOF {
			return LEXER_EOF
		} else {
			l.ReadError = err
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
			l.ReadError = err
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

func (l Lexer) MakeToken(d Data) *Token {
	return &Token{l.Position, d}
}

func (l Lexer) Error(msg string, args ...interface{}) *Token {
	loc := fmt.Sprintf(l.Source, l.Line, l.Column)
	err := fmt.Errorf(loc+msg, args...)
	return l.MakeToken(Error{err})
}

func (l Lexer) LexComment() *Token {
	buf := strings.Builder{}
	l.Get() // skip #
	for r := l.Get(); r != '\n'; r = l.Get() {
		if r < 0 {
			return l.Error("Unexpected EOF or read error.")
		}
		buf.WriteRune(r)
	}
	return l.MakeToken(String(buf.String()))
}

func (l *Lexer) lexField(buf *strings.Builder) *Token {
	for r := l.Peek(); (r >= '0' && r <= '9') ||
		(r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		r == '.' ||
		r == '_' ||
		r < 0; r = l.Peek() {
		if r < 0 {
			return l.Error("Unexpected EOF or read error.")
		}
		buf.WriteRune(l.Get())
	}
	return nil // means no error here
}

func (l *Lexer) LexField() *Token {
	buf := strings.Builder{}
	err := l.lexField(&buf)
	if err != nil {
		return l.Error("Could not parse func call: %w", err)
	}
	return l.MakeToken(FieldName(buf.String()))
}

func (l *Lexer) LexFuncall() *Token {
	buf := strings.Builder{}
	l.Get() // skip @
	err := l.lexField(&buf)
	if err != nil {
		return l.Error("Could not parse func call: %w", err)
	}
	return l.MakeToken(FuncName(buf.String()))
}

func (l *Lexer) LexString() *Token {
	buf := strings.Builder{}
	l.Get() // skip '
	for r := l.Get(); r != '\''; r = l.Get() {
		if r < 0 {
			return l.Error("Unexpected EOF or read error.")
		}
		if r == '\\' { // handle simple escapes in string
			e := l.Get()
			if e < 0 {
				return l.Error("Unexpected EOF or read error after \\ escape.")
			}
			switch e {
			case 'n':
				r = '\n'
			case 'r':
				r = '\r'
			default:
				r = e
			}
		}
		buf.WriteRune(r)
	}
	return l.MakeToken(String(buf.String()))
}

func (l *Lexer) LexNumber() *Token {
	buf := strings.Builder{}
	buf.WriteRune(l.Get())
	isInt := true
	for r := l.Peek(); (r >= '0' && r <= '9') || r == '.' || r < 0; r = l.Peek() {
		if r < 0 {
			return l.Error("Unexpected EOF or read error.")
		}
		if r == '.' {
			if !isInt {
				return l.Error("More than one floating point in number.")
			}
			isInt = false
		}
		buf.WriteRune(l.Get())
	}
	var res Data
	var err error
	if isInt {
		var i int64
		i, err = strconv.ParseInt(buf.String(), 0, 64)
		res = Int(i)
	} else {
		var f float64
		f, err = strconv.ParseFloat(buf.String(), 64)
		res = Float(f)
	}
	if err != nil {
		return l.Error("Not a number: %w", err)
	}
	return l.MakeToken(res)
}

func (l *Lexer) Lex() *Token {
	r := l.Peek()
	if r == ' ' || r == '\t' {
		for r = l.Get(); r == ' ' || r == '\t'; r = l.Get() {
		}
	}
	switch r {
	case '#':
		return l.LexComment()
	case '@':
		return l.LexFuncall()
	case '\'':
		return l.LexString()
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return l.LexNumber()
	case LEXER_EOF:
		return nil
	default:
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' {
			return l.LexField()
		}
		return l.Error("Unexpected character %c", r)
	}

	return &Token{}
}
