package suela

import "testing"
import "github.com/matryer/is"

func TestParserPeekGet(t *testing.T) {
	is := is.New(t)
	par := NewParserFromString("parser_test.go", `1 2 3 4`)
	tok := par.Peek()
	is.Equal(tok.Data, Int(1))
	tok = par.Peek()
	is.Equal(tok.Data, Int(1))
	tok = par.Get()
	is.Equal(tok.Data, Int(1))
	tok = par.Peek()
	is.Equal(tok.Data, Int(2))
	tok = par.Get()
	is.Equal(tok.Data, Int(2))
	tok = par.Get()
	is.Equal(tok.Data, Int(3))
	tok = par.Get()
	is.Equal(tok.Data, Int(4))
	tok = par.Get()
	is.Equal(tok.Data, Nil{})
}

/*
type parserTestCase struct {
	input  string
	expect Data
}

func (ltc parserTestCase) run(is *is.I) {
	lex := MakeLexerFromString("parser_test.go", ltc.input)
	res := lex.Lex()
	if ltc.expect == nil { // special case for expected nils
		is.Equal(res, nil)
	} else { // special case for expected errors
		if _, isErr := ltc.expect.(Error); isErr {
			_, alsoErr := res.Data.(Error)
			is.True(alsoErr)
		} else {
			is.Equal(ltc.expect, res.Data)
		}
	}
}

func TestParserCases(t *testing.T) {
	is := is.New(t)
	cases := []parserTestCase{
		{``, End{}},
		{`,,`, Rune(',')},
		{`)(`, Rune(')')},
		{`()`, Rune('(')},
		{`foo`, FieldName("foo")},
		{"  \t  foo", FieldName("foo")},
		{`@foo`, FuncName("foo")},
		{`foo.bar.baz`, FieldName("foo.bar.baz")},
		{`@foo.bar.baz`, FuncName("foo.bar.baz")},
		{"#foo\n", Comment("foo")},
		{`'hello\'\n日本'`, String("hello'\n日本")},
		{`1234`, Int(1234)},
		{`0.1234`, Float(0.1234)},
		{`-567890`, Int(-567890)},
		{`-5678.90`, Float(-5678.90)},
		{`-5678.90.0`, Error{}},
		{`{"hello":"world"}`, Json([]byte(`{"hello":"world"}`))},
		{`{"hello}}}}":"wo{{{{\"\"rld", "foo": 3}`, Json([]byte(`{"hello}}}}":"wo{{{{\"\"rld", "foo": 3}`))},
	}
	for i, cas := range cases {
		t.Logf("Test case %d: %s", i, cas.input)
		cas.run(is)
	}
}
*/
