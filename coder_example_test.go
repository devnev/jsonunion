package jsonunion_test

import (
	"fmt"
	"reflect"

	"github.com/devnev/jsonunion"
)

func ExampleCoder() {
	fmt.Println()

	coder := &jsonunion.Coder{
		TagKey: "type",
		Tags:   []string{"hello", "goodbye"},
		Types:  []reflect.Type{reflect.TypeOf(&HelloAction{}), reflect.TypeOf(&GoodbyeAction{})},
	}

	action, _ := coder.Decode([]byte(`{"type":"hello","target":"world"}`))
	fmt.Printf("decoded: %#v\n", action)

	buf, _ := coder.Encode(&GoodbyeAction{UntilWhen: "soon"})
	fmt.Println("encoded:", string(buf))

	// Output:
	// decoded: &jsonunion_test.HelloAction{Target:"world"}
	// encoded: {"type": "goodbye","untilWhen":"soon"}
}
