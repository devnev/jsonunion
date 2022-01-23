package jsonunion_test

import (
	"encoding/json"
	"fmt"

	"github.com/devnev/jsonunion"
)

func ExampleMarshaler_separateFuncs() {
	fmt.Println()

	message := MessageWithUnionPropAndSeparateMethods{
		Action: &HelloAction{Target: "world"},
	}
	buf, _ := json.Marshal(message)
	fmt.Println("marshaled:", string(buf))

	var message2 MessageWithUnionPropAndSeparateMethods
	_ = json.Unmarshal(buf, &message2)
	fmt.Printf("unmarshaled: %#v\n", message2.Action)

	// unknown tag value
	var message3 MessageWithUnionPropAndSeparateMethods
	err := json.Unmarshal([]byte(`{"action":{"type":"unknown"}}`), &message3)
	fmt.Printf("unmarshaled: %#v (%v)\n", message3.Action, err)

	// not an object
	var message4 MessageWithUnionPropAndSeparateMethods
	err = json.Unmarshal([]byte(`{"action":"woopsy"}`), &message4)
	fmt.Printf("unmarshaled: %#v (%v)\n", message4.Action, err)

	// missing tag
	var message5 MessageWithUnionPropAndSeparateMethods
	err = json.Unmarshal([]byte(`{"action":{"foo":"bar"}}`), &message5)
	fmt.Printf("unmarshaled: %#v (%v)\n", message5.Action, err)

	// Output:
	// marshaled: {"action":{"type":"hello","target":"world"}}
	// unmarshaled: &jsonunion_test.HelloAction{Target:"world"}
	// unmarshaled: <nil> (unknown tag value "unknown")
	// unmarshaled: <nil> (expected an object)
	// unmarshaled: <nil> (missing tag property)
}

type MessageWithUnionPropAndSeparateMethods struct {
	Action Action `json:"-"`
}

func (m MessageWithUnionPropAndSeparateMethods) MarshalJSON() ([]byte, error) {
	type WithoutMarshal MessageWithUnionPropAndCombinedMethods
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

func (m *MessageWithUnionPropAndSeparateMethods) UnmarshalJSON(data []byte) error {
	type WithoutMarshal MessageWithUnionPropAndCombinedMethods
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
	Types:  []interface{}{&HelloAction{}, &GoodbyeAction{}},
}
