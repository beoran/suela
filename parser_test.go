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

type parserTestCase struct {
	input  string
	expect string
}

func (ltc parserTestCase) run(is *is.I, t *testing.T) {
	par := NewParserFromString("parser_test.go", ltc.input)
	res := par.Parse()
	t.Logf("parse result: %v: %s", res, res.Describe())
	is.Equal(ltc.expect, res.Describe())
}

func TestParserCases(t *testing.T) {
	is := is.New(t)
	cases := []parserTestCase{
		{`@nil()`, `(T () (C (@nil)))`},
		{`@printf('%s %s', 'Hello world', field)`, `(T () (C (@printf) (A ('%s %s')) (A ('Hello world')) (A (field))))`},
		{"@hello()\n@world()", `(T () (C (@hello)) (C (@world)))`},
		{`@set(a, @sum(2,3))`, `(T () (C (@set) (A (a)) (A (@sum) (C (@sum) (A (2)) (A (3))))))`},
	}
	for i, cas := range cases {
		t.Logf("Test case %d: %s", i, cas.input)
		cas.run(is, t)
	}
}
