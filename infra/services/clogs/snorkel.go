package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
)

const (
	IntType = 1
	StrType = 2
	SetType = 3
)

type EventData map[string]interface{}

type event struct {
	Table string    `json:"table"`
	Token string    `json:"token"`
	Data  EventData `json:"data"`
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
			intv, err := toInt(v)
			if err != nil {
				return fmt.Errorf("int type expected, key=%v, value=%v, %w", k, v, err)
			}
			out[k] = intv
		case StrType:
			_, ok := v.(string)
			if !ok {
				return fmt.Errorf("string type expected, key=%v, value=%v", k, v)
			}
			out[k] = v
		case SetType:
			_, ok := v.([]string)
			if !ok {
				return fmt.Errorf("expected an array of string type, but key=%v, value=%v", k, v)
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
	bsln := append(bs, []byte("\n")...)
	if nwritten, err := writer.Write(bsln); err != nil {
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

func toInt(v interface{}) (int64, error) {
	if v == nil {
		return 0, errors.New("not compatiable")
	}

	rv := reflect.ValueOf(v)
	rt := rv.Type()
	if !rt.ConvertibleTo(reflect.TypeOf(int64(0))) {
		return 0, errors.New("not compatiable")
	}
	value := rv.Convert(reflect.TypeOf(int64(0)))
	return value.Int(), nil
}
