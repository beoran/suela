package suela

import "fmt"
import "strings"

type Type string

const (
	TypeType      = Type("type")
	TypeString    = Type("string")
	TypeFuncName  = Type("funcname")
	TypeFieldName = Type("fieldname")
	TypeComment   = Type("comment")
	TypeInt       = Type("int")
	TypeFloat     = Type("float")
	TypeList      = Type("list")
	TypeMap       = Type("map")
	TypeJson      = Type("json")
	TypeError     = Type("error")
	TypeAst       = Type("ast")
	TypeToken     = Type("token")
	TypeFunc      = Type("func")
	TypeNil       = Type("nil")
	TypeEnd       = Type("end")
)

/*
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

type Nil struct{}

func (n Nil) Type() Type {
	return TypeNil
}

func (n Nil) String() string {
	return fmt.Sprintf("")
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
	return "@" + string(d)
}

type FieldName []string

func (d FieldName) Type() Type {
	return TypeFieldName
}

func (d FieldName) String() string {
	return strings.Join([]string(d), ".")
}

func MakeFieldName(s string) FieldName {
	return FieldName(strings.Split(s, "."))
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

var _ Data = Int(0)
var _ Data = String("")
*/

type Data struct {
	Type
	Str        *string
	Int        *int64
	Float      *float64
	Error      error
	Json       []byte
	StringList []string
	List       []*Data
	Map        map[string]*Data
	User       interface{ String() string }
}

func (d Data) String() string {
	switch d.Type {
	case TypeType:
		return fmt.Sprintf("%s", d.Type)
	case TypeString:
		return fmt.Sprintf("%s", d.Str)
	case TypeFuncName:
		return fmt.Sprintf("%s", d.Str)
	case TypeFieldName:
		return fmt.Sprintf("%s", d.StringList)
	case TypeComment:
		return fmt.Sprintf("%s", d.Str)
	case TypeInt:
		return fmt.Sprintf("%s", d.Int)
	case TypeFloat:
		return fmt.Sprintf("%s", d.Float)
	case TypeList:
		return fmt.Sprintf("%v", d.List)
	case TypeMap:
		return fmt.Sprintf("%v", d.Map)
	case TypeJson:
		return fmt.Sprintf("%s", d.Json)
	case TypeError:
		return fmt.Sprintf("%s", d.Error)
	case TypeNil:
		return "nil"
	case TypeEnd:
		return "eof"
	default:
		return d.User.String()
	}
}
