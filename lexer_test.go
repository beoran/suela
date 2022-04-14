package suela

import "testing"
import "github.com/matryer/is"

type lexerTestCase struct {
	input string
	TokenKind
	expect Data
}

func (ltc lexerTestCase) run(is *is.I, t *testing.T) {
	lex := MakeLexerFromString("lexer_test.go", ltc.input)
	res := lex.Lex()
	t.Logf("res: %v", res)
	if ltc.expect == nil { // special case for expected nils
		is.Equal(res, nil)
	} else {
		is.Equal(ltc.TokenKind, res.TokenKind) // Kind must be the same.
		// special case for expected errors
		if _, isErr := ltc.expect.(Error); isErr {
			_, alsoErr := res.Data.(Error)
			is.True(alsoErr)
		} else {
			is.Equal(ltc.expect, res.Data)
		}
	}
}

func TestLexerCases(t *testing.T) {
	is := is.New(t)
	cases := []lexerTestCase{
		{``, TokenKindEnd, Nil{}},
		{`,,`, TokenKindComma, Nil{}},
		{`)(`, TokenKindCp, Nil{}},
		{`()`, TokenKindOp, Nil{}},
		{`foo`, TokenKindField, MakeFieldName("foo")},
		{"  \t  foo", TokenKindField, MakeFieldName("foo")},
		{`@foo`, TokenKindFunc, FuncName("foo")},
		{`foo.bar.baz`, TokenKindField, MakeFieldName("foo.bar.baz")},
		{`foo.[1].bar.baz`, TokenKindField, MakeFieldName("foo.[1].bar.baz")},
		{`@foo.bar.baz`, TokenKindFunc, FuncName("foo.bar.baz")},
		{"#foo\n", TokenKindComment, Comment("foo")},
		{`'hello\'\n日本'`, TokenKindLiteral, String("hello'\n日本")},
		{`1234`, TokenKindLiteral, Int(1234)},
		{`0.1234`, TokenKindLiteral, Float(0.1234)},
		{`-567890`, TokenKindLiteral, Int(-567890)},
		{`-5678.90`, TokenKindLiteral, Float(-5678.90)},
		{`-5678.90.0`, TokenKindError, Error{}},
		{`{"hello":"world"}`, TokenKindLiteral, Json([]byte(`{"hello":"world"}`))},
		{`{"hello}}}}":"wo{{{{\"\"rld", "foo": 3}`,
			TokenKindLiteral,
			Json([]byte(`{"hello}}}}":"wo{{{{\"\"rld", "foo": 3}`))},
	}
	for i, cas := range cases {
		t.Logf("Test case %d: <<%s>>", i, cas.input)
		cas.run(is, t)
	}
}
