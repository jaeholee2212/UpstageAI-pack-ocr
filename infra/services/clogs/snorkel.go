package main

const (
	IntType = 1
	StrType = 2
	SetType = 3
)

type Snorkel struct {
	table string
	token string
	spec  map[string]int
}

func NewSnorkel(table string, token string) *Snorkel {
	return &Snorkel{
		table: table,
		token: token,
		spec:  map[string]int{},
	}
}

func (s *Snorkel) AddIntField(name string) *Snorkel {
	return s.addField(name, IntType)
}

func (s *Snorkel) AddStrField(name string) *Snorkel {
	return s.addField(name, StrType)
}

func (s *Snorkel) AddSetField(name string) *Snorkel {
	return s.addField(name, SetType)
}

func (s *Snorkel) addField(name string, fieldType int) *Snorkel {
	s.spec[name] = fieldType
	return s
}
