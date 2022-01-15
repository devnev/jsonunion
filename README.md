# JSON tagged union encoding and decoding for Go

In JavaScript and JSON, unions are often represented using an internal tag, e.g.

```json
{
  "type": "hello",
  "target": "world"
}
```

Or written as TypeScript types:

```ts
type Actions =
  | { type: "hello"; target: string }
  | { type: "goodbye"; untilWhen?: string };
```

This package supports easy encoding and decoding such types in Go, including
integrating with `encoding/json.Marshal` and `encoding/json.Unmarshal`.

## Installation

Add this to your module with

```sh
go get github.com/devnev/jsonunion
```

## Usage

See examples in the Go package documentation.
