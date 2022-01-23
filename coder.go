package jsonunion

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

var (
	Err           = errors.New("")
	ErrInputType  = fmt.Errorf("expected an object%w", Err)
	ErrTagMissing = fmt.Errorf("missing tag property%w", Err)
	ErrTagType    = fmt.Errorf("tag value must be a string%w", Err)
	ErrTagValue   = fmt.Errorf("unknown tag value%w", Err)
)

type Coder struct {
	TagKey          string
	Tags            []string
	Types           []interface{}
	RequireTagFirst bool
}

func (c *Coder) Encode(v interface{}) ([]byte, error) {
	buf, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return c.InsertTag(v, buf)
}

func (c *Coder) Decode(data []byte) (interface{}, error) {
	dstType, err := c.DecodeTag(data)
	if dstType == nil || err != nil {
		return nil, err
	}
	dst := reflect.New(dstType)
	err = json.Unmarshal(data, dst.Interface())
	return dst.Elem().Interface(), err
}

func (c *Coder) DecodeTag(data []byte) (reflect.Type, error) {
	dec := json.NewDecoder(bytes.NewReader(data))

	tok, err := dec.Token()
	if tok == nil || err != nil {
		return nil, err
	}
	if tok != json.Delim('{') {
		return nil, ErrInputType
	}

	var tagValue string

	for {
		tok, err = dec.Token()
		if err != nil {
			return nil, err
		}

		if tok == json.Delim('}') {
			return nil, ErrTagMissing
		}

		key := tok.(string)

		if c.RequireTagFirst {
			if key != c.TagKey {
				return nil, fmt.Errorf("%w or not at start", ErrTagMissing)
			}
			tok, err = dec.Token()
			if tok == nil || err != nil {
				return nil, err
			}
		} else {
			for depth := 0; true; {
				tok, err = dec.Token()
				if tok == nil || err != nil {
					return nil, err
				}
				if tok == json.Delim('[') || tok == json.Delim('{') {
					depth++
				} else if tok == json.Delim(']') || tok == json.Delim('}') {
					depth--
				}
				if depth == 0 {
					break
				}
			}

			if key != c.TagKey {
				continue
			}
		}

		switch t := tok.(type) {
		case string:
			tagValue = t
		default:
			return nil, ErrTagType
		}

		break
	}

	var dstType reflect.Type
	for i := range c.Tags {
		if c.Tags[i] == tagValue {
			dstType = reflect.TypeOf(c.Types[i])
			break
		}
	}
	if dstType == nil {
		return nil, fmt.Errorf("%w %q", ErrTagValue, tagValue)
	}

	return dstType, nil
}

func (c *Coder) InsertTag(v interface{}, encoded []byte) ([]byte, error) {
	reflected := reflect.ValueOf(v)
	if !reflected.IsValid() {
		return encoded, nil
	}

	// allow passing a pointer to the union field so that the same pointer can
	// be passed to both DecodeTag and InsertTag, as in the combined example
	if reflected.Kind() == reflect.Ptr && reflected.Elem().Kind() == reflect.Interface {
		reflected = reflected.Elem().Elem()
	}

	reflectedType := reflected.Type()
	var tagValue string
	for i := range c.Types {
		if reflect.TypeOf(c.Types[i]) == reflectedType {
			tagValue = c.Tags[i]
			break
		}
	}

	if tagValue == "" {
		panic("bad destination type")
	}

	dec := json.NewDecoder(bytes.NewReader(encoded))
	firstToken, _ := dec.Token()
	if firstToken == nil {
		return encoded, nil
	}
	if firstToken != json.Delim('{') {
		panic("cannot add tag to non-object")
	}
	var empty bool
	if tok, _ := dec.Token(); tok == json.Delim('}') {
		empty = true
	}

	objStart := bytes.IndexRune(encoded, '{')
	keyOffset := len(bytes.TrimSpace(encoded)) - len(bytes.TrimSpace(encoded[objStart+1:]))
	spaceBytes := keyOffset - 1
	const delimBytes = 3 // colon, space, comma
	jsonTag, _ := json.Marshal(c.TagKey)
	jsonTagValue, _ := json.Marshal(tagValue)

	dstBuf := bytes.NewBuffer(nil)
	dstBuf.Grow(len(encoded) + spaceBytes + delimBytes + len(jsonTag) + len(jsonTagValue))
	dstBuf.Write(encoded[:objStart+keyOffset])
	dstBuf.Write(jsonTag)
	dstBuf.WriteRune(':')
	dstBuf.WriteRune(' ')
	dstBuf.Write(jsonTagValue)
	if !empty {
		dstBuf.WriteRune(',')
	}
	dstBuf.Write(encoded[objStart+1:])

	return dstBuf.Bytes(), nil
}
