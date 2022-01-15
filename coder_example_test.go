package jsonunion_test

import (
	"fmt"

	"github.com/devnev/jsonunion"
)

func ExampleCoder_Encode() {
	coder := &jsonunion.Coder{
		TagKey: "type",
		Tags:   []string{"hello", "goodbye"},
		Types:  []interface{}{&HelloAction{}, &GoodbyeAction{}},
	}

	buf, _ := coder.Encode(&GoodbyeAction{UntilWhen: "soon"})
	fmt.Println(string(buf))
	// Output: {"type": "goodbye","untilWhen":"soon"}
}

func ExampleCoder_Decode() {
	coder := &jsonunion.Coder{
		TagKey: "type",
		Tags:   []string{"hello", "goodbye"},
		Types:  []interface{}{&HelloAction{}, &GoodbyeAction{}},
	}

	action, _ := coder.Decode([]byte(`{"type":"hello","target":"world"}`))
	fmt.Printf("%#v\n", action)
	// Output: &jsonunion_test.HelloAction{Target:"world"}
}
