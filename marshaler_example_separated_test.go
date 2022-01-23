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

	// Output:
	// marshaled: {"action":{"type":"hello","target":"world"}}
	// unmarshaled: &jsonunion_test.HelloAction{Target:"world"}
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
