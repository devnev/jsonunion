# JSON tagged union encoding and decoding for Go

In JavaScript and JSON, unions are often represented using an internal tag, e.g. (written as TypeSsript types):

```ts
type Actions =
  | { type: "hello"; target: string }
  | { type: "goodbye"; untilWhen?: string };
```

This package supports easy encoding and decoding such types in Go, including
integrating with `encoding/json.Marshal` and `encoding/json.Unmarshal`, such
that JSON matching the above TypeScript type can easily be mapped to something
like the following in Go:

```go
type Action interface {
	isAction()
}

type HelloAction struct {
	Target string `json:"target"`
}

type GoodbyeAction struct {
	UntilWhen string `json:"untilWhen,omitempty"`
}

func (*HelloAction) isAction()   {}
func (*GoodbyeAction) isAction() {}
```

## Installation

Add this to your module with

```sh
go get github.com/devnev/jsonunion
```

## Usage

See examples in the Go package documentation or [in the source code](marshaler_example_combined_test.go).
