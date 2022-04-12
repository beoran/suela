package suela

import "fmt"

type Type string

const TypeType = Type("type")
const TypeString = Type("string")
const TypeFuncName = Type("funcname")
const TypeFieldName = Type("fieldname")
const TypeComment = Type("comment")
const TypeInt = Type("int")
const TypeFloat = Type("float")
const TypeList = Type("list")
const TypeMap = Type("map")
const TypeJson = Type("Jjon")
const TypeError = Type("error")
const TypeAst = Type("ast")
const TypeToken = Type("token")
const TypeFunc = Type("func")
const TypeRune = Type("rune")

type Data interface {
	Type() Type
	String() string
}

type Int int64

func (d Int) Type() Type {
	return TypeInt
}

func (d Int) String() string {
	return fmt.Sprintf("%d", d)
}

// Rune is used only by the lexer for syntactical characters like (),
type Rune int64

func (d Rune) Type() Type {
	return TypeRune
}

func (d Rune) String() string {
	return fmt.Sprintf("%c", d)
}

type Float float64

func (d Float) Type() Type {
	return TypeFloat
}

func (d Float) String() string {
	return fmt.Sprintf("%f", d)
}

type Comment string

func (d Comment) Type() Type {
	return TypeComment
}

func (d Comment) String() string {
	return "#" + string(d) + "\n"
}

type String string

func (d String) Type() Type {
	return TypeString
}

func (d String) String() string {
	return "'" + string(d) + "'"
}

type FuncName string

func (d FuncName) Type() Type {
	return TypeFuncName
}

func (d FuncName) String() string {
	return "$" + string(d)
}

type FieldName string

func (d FieldName) Type() Type {
	return TypeFieldName
}

func (d FieldName) String() string {
	return string(d)
}

type List []Data

func (d List) Type() Type {
	return TypeList
}

func (d List) String() string {
	s := "@list("
	for i, e := range d {
		if i > 0 {
			s += ","
		}
		s += fmt.Sprintf("%s", e.String())
	}
	s += ")"
	return s
}

type Map map[string]Data

func (d Map) Type() Type {
	return TypeMap
}

func (d Map) String() string {
	s := "@map("
	i := 0
	for k, e := range d {
		if i > 0 {
			s += ","
		}
		s += fmt.Sprintf("%s,%s", String(k).String(), e.String())
		i++
	}
	s += ")"
	return s
}

type Json []byte

func (d Json) Type() Type {
	return TypeJson
}

func (d Json) String() string {
	return string(d)
}

type Error struct {
	Err error
}

func (d Error) Type() Type {
	return TypeError
}

func (d Error) String() string {
	return d.Err.Error()
}

func (d Error) Error() string {
	return d.Err.Error()
}

type Suela struct {
}

type Func struct {
	Func func(Suela, ...Data) Data
}

var _ Data = Int(0)
var _ Data = String("")
