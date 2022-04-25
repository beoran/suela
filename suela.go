package suela

type Impl func(*Suela, ...Data) Data

type Func struct {
	Impl
}

type Suela struct {
	Funcs  map[string]Func
	Vars   map[string]Data
	Lookup []func(string) Data
	Stack  []Data
}

func New(lookup ...func(string) Data) *Suela {
	return &Suela{make(map[string]Func), make(map[string]Data), lookup, []Data{}}
}

func (s *Suela) Func(name string, impl Impl) {
	s.Funcs[name] = Func{impl}
}

func (s *Suela) Var(name string, data Data) {
	s.Vars[name] = data
}
