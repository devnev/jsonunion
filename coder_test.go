package jsonunion

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			title: "tag_value_after_empty_object_field",
			coder: coderWithValueStructATag,
			input: `{"value": {}, "type": "a"}`,
		},
		{
			title: "tag_value_after_array_field",
			coder: coderWithValueStructATag,
			input: `{"value": [1], "type": "a"}`,
		},
		{
			title: "tag_value_after_null_field",
			coder: coderWithValueStructATag,
			input: `{"value": null, "type": "a"}`,
		},
		{
			title: "tag_value_after_object_field_with_complex_sub_object",
			coder: coderWithValueStructATag,
			input: `{"value": {"a": 1, "b": [2, 3], "c": {"d": 4, "e": null}}, "type": "a"}`,
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

func TestCoder_EncodeDecodeRoundtripsOK(t *testing.T) {
	for _, tc := range []struct {
		title  string
		coder  *Coder
		input  interface{}
		output interface{}
	}{
		{
			title: "empty_struct",
			coder: coderWithEmpyStructATag,
			input: EmpyStruct{},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			buf, err := tc.coder.Encode(tc.input)
			require.NoError(t, err)
			t.Logf("encoded is %q", string(buf))
			val, err := tc.coder.Decode(buf)
			require.NoError(t, err)
			if tc.output != nil {
				assert.Equal(t, tc.output, val)
			} else {
				assert.Equal(t, tc.input, val)
			}
		})
	}
}

func TestCoder_InsertingTagPreservesCustomSpacing(t *testing.T) {
	for _, tc := range []struct {
		title        string
		valueJSON    string
		expectedJSON string
	}{
		{
			title:        "empty_object_extra_spaces",
			valueJSON:    `   {       }       `,
			expectedJSON: `   {       "type": "a"       }       `,
		},
		{
			title:        "number_object_extra_spaces",
			valueJSON:    `   {    "b": 1      }       `,
			expectedJSON: `   {    "type": "a",    "b": 1      }       `,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			buf, err := coderWithEmpyStructATag.InsertTag(EmpyStruct{}, []byte(tc.valueJSON))
			require.NoError(t, err)
			assert.Equal(t, tc.expectedJSON, string(buf))
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
			title:  "object_with_value_instead_of_properties",
			coder:  coderWithEmpyStructATag,
			input:  "{null}",
			errstr: "invalid character 'n'",
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
			errval: ErrTagType,
		},
		{
			title:  "null_tag",
			coder:  coderWithEmpyStructATag,
			input:  `{"type": null}`,
			errstr: "tag value must be a string",
			errval: ErrTagType,
		},
		{
			title:  "tag_value_after_number_but_required_at_start",
			coder:  coderWithNumberStructAndRequiringATagAtStart,
			input:  `{"number": 1, "type": "a"}`,
			errstr: "missing tag property or not at start",
			errval: ErrTagMissing,
		},
		{
			title:  "tag_value_after_null_but_required_at_start",
			coder:  coderWithNumberStructAndRequiringATagAtStart,
			input:  `{"value": null, "type": "a"}`,
			errstr: "missing tag property or not at start",
			errval: ErrTagMissing,
		},
		{
			title:  "null_tag_and_required_at_sart",
			coder:  coderWithNumberStructAndRequiringATagAtStart,
			input:  `{"type": null}`,
			errstr: "tag value must be a string",
			errval: ErrTagType,
		},
		{
			title:  "number_tag_and_required_at_sart",
			coder:  coderWithNumberStructAndRequiringATagAtStart,
			input:  `{"type": 1}`,
			errstr: "tag value must be a string",
			errval: ErrTagType,
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

type (
	EmpyStruct   struct{}
	NumberStruct struct {
		Number float32 `json:"number"`
	}
	ValueStruct struct {
		Value interface{} `json:"value"`
	}
)

var (
	coderWithNoTags = &Coder{
		TagKey: "type",
	}
	coderWithEmpyStructATag = &Coder{
		TagKey: "type",
		Tags:   []string{"a"},
		Types:  []interface{}{EmpyStruct{}},
	}
	coderWithNumberStructATag = &Coder{
		TagKey: "type",
		Tags:   []string{"a"},
		Types:  []interface{}{NumberStruct{}},
	}
	coderWithValueStructATag = &Coder{
		TagKey: "type",
		Tags:   []string{"a"},
		Types:  []interface{}{ValueStruct{}},
	}
	coderWithNumberStructAndRequiringATagAtStart = &Coder{
		TagKey:          "type",
		Tags:            []string{"a"},
		Types:           []interface{}{NumberStruct{}},
		RequireTagFirst: true,
	}
)
