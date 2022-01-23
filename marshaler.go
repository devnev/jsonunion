package jsonunion

import (
	"encoding/json"
	"errors"
	"reflect"
)

type Marshaler struct {
	Value       interface{}
	Coder       *Coder
	NilOnErrors bool
}

func (m Marshaler) MarshalJSON() ([]byte, error) {
	buf, err := json.Marshal(m.Value)
	if err != nil {
		return nil, err
	}

	return m.Coder.InsertTag(m.Value, buf), nil
}

func (m *Marshaler) UnmarshalJSON(data []byte) error {
	dstType, err := m.Coder.DecodeTag(data)
	if dstType == nil || err != nil {
		if m.NilOnErrors && errors.Is(err, Err) {
			return nil
		}
		return err
	}

	dst := reflect.New(dstType)
	err = json.Unmarshal(data, dst.Interface())
	if err != nil {
		return err
	}

	reflect.ValueOf(m.Value).Elem().Set(dst.Elem())
	return nil
}
