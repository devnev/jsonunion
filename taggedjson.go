package taggedjson

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrInputType  = errors.New("expected an object")
	ErrTagMissing = errors.New("missing tag property")
	ErrTagType    = errors.New("tag value must be a string")
	ErrTagValue   = errors.New("unknown tag value")
)

type Coder struct {
	TagKey          string
	Tags            []string
	Types           []reflect.Type
	RequireTagFirst bool
}

func (c *Coder) Decode(data []byte) (interface{}, error) {
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
		if tok == nil || err != nil {
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
			dstType = c.Types[i]
			break
		}
	}
	if dstType == nil {
		return nil, fmt.Errorf("%w %q", ErrTagValue, tagValue)
	}

	dst := reflect.New(dstType)
	err = json.Unmarshal(data, dst.Interface())
	return dst.Elem().Interface(), err
}

func (c *Coder) Encode(v interface{}) ([]byte, error) {
	reflected := reflect.ValueOf(v)
	for reflected.Kind() == reflect.Interface {
		reflected = reflected.Elem()
	}
	if !reflected.IsValid() {
		return json.Marshal(v)
	}

	reflectedType := reflected.Type()
	var tagValue string
	for i := range c.Types {
		if c.Types[i] == reflectedType {
			tagValue = c.Tags[i]
			break
		}
	}

	if tagValue == "" {
		panic("bad destination type")
	}

	buf, err := json.Marshal(reflected.Interface())
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(bytes.NewReader(buf))
	firstToken, _ := dec.Token()
	if firstToken == nil {
		return buf, nil
	}
	if firstToken != json.Delim('{') {
		panic("cannot add tag to non-object")
	}
	var empty bool
	if tok, _ := dec.Token(); tok == json.Delim('}') {
		empty = true
	}

	objStart := bytes.IndexRune(buf, '{')
	keyOffset := len(bytes.TrimSpace(buf)) - len(bytes.TrimSpace(buf[objStart+1:]))
	spaceBytes := keyOffset - 1
	const delimBytes = 3 // colon, space, comma
	jsonTag, _ := json.Marshal(c.TagKey)
	jsonTagValue, _ := json.Marshal(tagValue)

	dstBuf := bytes.NewBuffer(nil)
	dstBuf.Grow(len(buf) + spaceBytes + delimBytes + len(jsonTag) + len(jsonTagValue))
	dstBuf.Write(buf[:objStart+keyOffset])
	dstBuf.Write(jsonTag)
	dstBuf.WriteRune(':')
	dstBuf.WriteRune(' ')
	dstBuf.Write(jsonTagValue)
	if !empty {
		dstBuf.WriteRune(',')
	}
	dstBuf.Write(buf[objStart+1:])

	return dstBuf.Bytes(), nil
}
