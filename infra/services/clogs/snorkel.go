package main

import (
	"encoding/json"
	"fmt"
	"io"
)

const (
	IntType = 1
	StrType = 2
	SetType = 3
)

type event struct {
	Table string                 `json:"table"`
	Token string                 `json:"token"`
	Data  map[string]interface{} `json:"data"`
}

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

func (s *Snorkel) Write(writer io.Writer, data map[string]interface{}) error {
	out := map[string]interface{}{}
	for k, v := range data {
		ftype, ok := s.spec[k]
		if !ok {
			continue
		}
		switch ftype {
		case IntType:
			vv, ok := v.(int)
			if !ok {
				return fmt.Errorf("int type expected, %v", vv)
			}
			out[k] = v
		case StrType:
			vv, ok := v.(string)
			if !ok {
				return fmt.Errorf("string type expected, %v", vv)
			}
			out[k] = v
		case SetType:
			vv, ok := v.([]string)
			if !ok {
				return fmt.Errorf("expected an array of string type, but %v", vv)
			}
			out[k] = v
		default:
			return fmt.Errorf("unknown type is given, type=%v", ftype)
		}
	}

	evt := event{
		Table: s.table,
		Token: s.token,
		Data:  out,
	}

	bs := toJson(evt)
	if nwritten, err := writer.Write(bs); err != nil {
		return err
	} else if nwritten != len(bs) {
		return fmt.Errorf(
			"failed to write all data, written=%v, expected=%v",
			nwritten,
			len(bs))
	}
	return nil
}

func toJson(e event) []byte {
	bs, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return bs
}
