package taggedjson

import (
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type (
	EmpyStruct   struct{}
	NumberStruct struct {
		Number float32 `json:"number"`
	}
)

var (
	coderWithNoTags = &Coder{
		TagKey: "type",
	}
	coderWithEmpyStructATag = &Coder{
		TagKey: "type",
		Tags:   []string{"a"},
		Types:  []reflect.Type{reflect.TypeOf(EmpyStruct{})},
	}
	coderWithNumberStructATag = &Coder{
		TagKey: "type",
		Tags:   []string{"a"},
		Types:  []reflect.Type{reflect.TypeOf(NumberStruct{})},
	}
	coderWithNumberStructAndRequiringATagAtStart = &Coder{
		TagKey:          "type",
		Tags:            []string{"a"},
		Types:           []reflect.Type{reflect.TypeOf(NumberStruct{})},
		RequireTagFirst: true,
	}
)

func TestCoder_DecodeEncodeRoundtripsOK(t *testing.T) {
	for _, tc := range []struct {
		title string
		coder *Coder
		input string
	}{
		{
			title: "empty_can_handle_null",
			coder: &Coder{},
			input: "null",
		},
		{
			title: "no_tag_values_can_handle_null",
			coder: coderWithNoTags,
			input: "null",
		},
		{
			title: "tag_value_no_fields",
			coder: coderWithEmpyStructATag,
			input: `{"type": "a"}`,
		},
		{
			title: "tag_value_with_number_field",
			coder: coderWithNumberStructATag,
			input: `{"type": "a", "number": 1}`,
		},
		{
			title: "tag_value_after_number_field",
			coder: coderWithNumberStructATag,
			input: `{"number": 1, "type": "a"}`,
		},
		{
			title: "tag_value_required_at_start_with_number_field",
			coder: coderWithNumberStructAndRequiringATagAtStart,
			input: `{"type": "a", "number": 1}`,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			val, err := tc.coder.Decode([]byte(tc.input))
			require.NoError(t, err)
			t.Logf("value is %#v", val)
			buf, err := tc.coder.Encode(val)
			require.NoError(t, err)
			t.Logf("output is %q", string(buf))
			assert.JSONEq(t, tc.input, string(buf))
		})
	}
}

func TestCoder_DecodeFails(t *testing.T) {
	for _, tc := range []struct {
		title  string
		coder  *Coder
		input  string
		errstr string
		errval error
	}{
		{
			title:  "empty_input",
			coder:  coderWithEmpyStructATag,
			input:  "",
			errstr: "EOF",
			errval: io.EOF,
		},
		{
			title:  "bad_first_token",
			coder:  coderWithEmpyStructATag,
			input:  "bad",
			errstr: "invalid character 'b' looking for beginning of value",
		},
		{
			title:  "string_instead_of_object",
			coder:  coderWithEmpyStructATag,
			input:  `""`,
			errstr: "expected an object",
			errval: ErrInputType,
		},
		{
			title:  "number_instead_of_object",
			coder:  coderWithEmpyStructATag,
			input:  `1`,
			errstr: "expected an object",
			errval: ErrInputType,
		},
		{
			title:  "missing_tag",
			coder:  coderWithEmpyStructATag,
			input:  "{}",
			errstr: "missing tag property",
		},
		{
			title:  "unknown_tag",
			coder:  coderWithEmpyStructATag,
			input:  `{"type": "b"}`,
			errstr: `unknown tag value "b"`,
		},
		{
			title:  "number_tag",
			coder:  coderWithEmpyStructATag,
			input:  `{"type": 0}`,
			errstr: "tag value must be a string",
		},
		{
			title:  "tag_value_not_at_start_as_Required",
			coder:  coderWithNumberStructAndRequiringATagAtStart,
			input:  `{"number": 1, "type": "a"}`,
			errstr: "missing tag property or not at start",
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			_, err := tc.coder.Decode([]byte(tc.input))
			assert.EqualError(t, err, tc.errstr)
			if tc.errval != nil {
				assert.ErrorIs(t, err, tc.errval)
			}
		})
	}
}
