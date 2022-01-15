package jsonunion_test

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/devnev/jsonunion"
)

func ExampleMarshaler() {
	fmt.Println()

	message := MessageWithUnionProp{
		Action: &HelloAction{Target: "world"},
	}
	buf, _ := json.Marshal(message)
	fmt.Println("marshaled:", string(buf))

	var message2 MessageWithUnionProp
	_ = json.Unmarshal(buf, &message2)
	fmt.Printf("unmarshaled: %#v\n", message2.Action)

	// Output:
	// marshaled: {"action":{"type":"hello","target":"world"}}
	// unmarshaled: &jsonunion_test.HelloAction{Target:"world"}
}

type MessageWithUnionProp struct {
	Action Action `json:"-"`
}

func (m MessageWithUnionProp) MarshalJSON() ([]byte, error) {
	type WithoutMarshal MessageWithUnionProp
	return json.Marshal(struct {
		WithoutMarshal
		Action jsonunion.Marshaler `json:"action"`
	}{
		WithoutMarshal: WithoutMarshal(m),
		Action: jsonunion.Marshaler{
			Value: m.Action,
			Coder: actionCoder,
		},
	})
}

func (m *MessageWithUnionProp) UnmarshalJSON(data []byte) error {
	type WithoutMarshal MessageWithUnionProp
	return json.Unmarshal(data, &struct {
		*WithoutMarshal
		Action jsonunion.Marshaler `json:"action"`
	}{
		WithoutMarshal: (*WithoutMarshal)(m),
		Action: jsonunion.Marshaler{
			Value: &m.Action,
			Coder: actionCoder,
		},
	})
}

var actionCoder = &jsonunion.Coder{
	TagKey: "type",
	Tags:   []string{"hello", "goodbye"},
	Types:  []reflect.Type{reflect.TypeOf(&HelloAction{}), reflect.TypeOf(&GoodbyeAction{})},
}
