package jsonunion

import "reflect"

type Marshaler struct {
	Value interface{}
	Coder *Coder
}

func (m Marshaler) MarshalJSON() ([]byte, error) {
	return m.Coder.Encode(m.Value)
}

func (m *Marshaler) UnmarshalJSON(data []byte) error {
	value, err := m.Coder.Decode(data)
	if err != nil {
		return err
	}
	reflect.ValueOf(m.Value).Elem().Set(reflect.ValueOf(value))
	return nil
}
