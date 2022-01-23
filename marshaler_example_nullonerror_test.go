package jsonunion_test

import (
	"encoding/json"
	"fmt"

	"github.com/devnev/jsonunion"
)

func ExampleMarshaler_nullOnError() {
	fmt.Println()

	message := MessageWithUnionPropAndNullForErrors{
		Action: &HelloAction{Target: "world"},
	}
	buf, _ := json.Marshal(message)
	fmt.Println("marshaled:", string(buf))

	var message2 MessageWithUnionPropAndNullForErrors
	_ = json.Unmarshal(buf, &message2)
	fmt.Printf("unmarshaled: %#v\n", message2.Action)

	// unknown tag value
	var message3 MessageWithUnionPropAndNullForErrors
	err := json.Unmarshal([]byte(`{"action":{"type":"unknown"}}`), &message3)
	fmt.Printf("unmarshaled: %#v (%v)\n", message3.Action, err)

	// not an object
	var message4 MessageWithUnionPropAndNullForErrors
	err = json.Unmarshal([]byte(`{"action":"woopsy"}`), &message4)
	fmt.Printf("unmarshaled: %#v (%v)\n", message4.Action, err)

	// missing tag
	var message5 MessageWithUnionPropAndNullForErrors
	err = json.Unmarshal([]byte(`{"action":{"foo":"bar"}}`), &message5)
	fmt.Printf("unmarshaled: %#v (%v)\n", message5.Action, err)

	// Output:
	// marshaled: {"action":{"type":"hello","target":"world"}}
	// unmarshaled: &jsonunion_test.HelloAction{Target:"world"}
	// unmarshaled: <nil> (<nil>)
	// unmarshaled: <nil> (<nil>)
	// unmarshaled: <nil> (<nil>)
}

type MessageWithUnionPropAndNullForErrors struct {
	Action Action `json:"-"`
}

func (m *MessageWithUnionPropAndNullForErrors) wrapForJSON() interface{} {
	// In this example we can combine everything into one definition inside this
	// method, with some potential performance downsides due to extra allocation
	// and reflect calls.
	type WithoutMarshal MessageWithUnionPropAndNullForErrors
	return &struct {
		*WithoutMarshal
		Action jsonunion.Marshaler `json:"action"`
	}{
		WithoutMarshal: (*WithoutMarshal)(m),
		Action: jsonunion.Marshaler{
			Value: &m.Action,
			Coder: &jsonunion.Coder{
				TagKey: "type",
				Tags:   []string{"hello", "goodbye"},
				Types:  []interface{}{&HelloAction{}, &GoodbyeAction{}},
			},
			NilOnErrors: true,
		},
	}
}

func (m MessageWithUnionPropAndNullForErrors) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.wrapForJSON())
}

func (m *MessageWithUnionPropAndNullForErrors) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, m.wrapForJSON())
}
