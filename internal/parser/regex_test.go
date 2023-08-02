package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestJsonTagRegex(t *testing.T) {
	tt := []struct {
		in  string
		out JsonTag
	}{
		{in: ``},
		{in: `db:"id"`},
		{in: `json:"id"`, out: JsonTag{Name: "id", Value: "id"}},
		{in: `json:"id,whatever"`, out: JsonTag{Name: "id", Value: "id,whatever"}},
		{in: `json:"id,whatever,else"`, out: JsonTag{Name: "id", Value: "id,whatever,else"}},
		{in: `json:"id,string"`, out: JsonTag{Name: "id", Value: "id,string", IsString: true}},
		{in: `json:"id,string,omit"`, out: JsonTag{Name: "id", Value: "id,string,omit", IsString: true}},
		{in: `json:"id,string,omitempty"`, out: JsonTag{Name: "id", Value: "id,string,omitempty", IsString: true, Omitempty: true}},
		{in: `json:"id,omitempty,string"`, out: JsonTag{Name: "id", Value: "id,omitempty,string", IsString: true, Omitempty: true}},
		{in: `json:"id,string,omitempty"`, out: JsonTag{Name: "id", Value: "id,string,omitempty", IsString: true, Omitempty: true}},
		{in: `json:"ID,string,omitempty"`, out: JsonTag{Name: "ID", Value: "ID,string,omitempty", IsString: true, Omitempty: true}},
		{in: `json:"renamed_fieldName99"`, out: JsonTag{Name: "renamed_fieldName99", Value: "renamed_fieldName99"}},
		{in: `xxx:"X X X" json:"id,string" yyy:"Y Y Y"`, out: JsonTag{Name: "id", Value: "id,string", IsString: true}},
		{in: `db:"id,omitempty,pk" json:"id,string"`, out: JsonTag{Name: "id", Value: "id,string", IsString: true}},
		{in: `db:"id,omitempty,pk" json:"External_ID,string,omitempty" someOtherTag:"some,other:value"`, out: JsonTag{Name: "External_ID", Value: "External_ID,string,omitempty", IsString: true, Omitempty: true}},
	}
	for _, tc := range tt {
		jsonTag, ok := GetJsonTag(tc.in)
		if ok != (tc.out.Value != "") {
			t.Errorf("expected ok=%v", tc.out)
		}

		if !cmp.Equal(jsonTag, tc.out) {
			t.Errorf(cmp.Diff(jsonTag, tc.out))
		}
	}
}

func TestTextMarshalerRegex(t *testing.T) {
	tt := []string{
		"func (github.com/google/uuid.UUID).MarshalText() ([]byte, error)",
		"func (github.com/google/uuid.UUID).MarshalText() (data []byte, err error)",
		"func (github.com/golang-cz/gospeak/uuid.UUID).MarshalText() ([]byte, error)",
		"func (github.com/golang-cz/gospeak/uuid.UUID).MarshalText() (b []byte, err error)",
	}
	for _, tc := range tt {
		if !textMarshalerRegex.MatchString(tc) {
			t.Errorf("textMarshalerRegex didn't match %q", tc)
		}
	}
}

func TestTextUnmarshalerRegex(t *testing.T) {
	tt := []string{
		"func (github.com/google/uuid.UUID).UnmarshalText(data []byte) (err error)",
		"func (github.com/google/uuid.UUID).UnmarshalText(data []byte) error",
		"func (*github.com/golang-cz/gospeak/uuid.UUID).UnmarshalText(b []byte) (err error)",
		"func (*github.com/golang-cz/gospeak/uuid.UUID).UnmarshalText(b []byte) error",
	}
	for _, tc := range tt {
		if !textUnmarshalerRegex.MatchString(tc) {
			t.Errorf("textUnmarshalerRegex didn't match %q", tc)
		}
	}
}

func TestJsonMarshalerRegex(t *testing.T) {
	tt := []string{
		"func (github.com/golang-cz/gospeak/data.Person).MarshalJSON() ([]byte, error)",
		"func (github.com/golang-cz/gospeak/data.Person).MarshalJSON() (data []byte, err error)",
	}
	for _, tc := range tt {
		if !jsonMarshalerRegex.MatchString(tc) {
			t.Errorf("jsonMarshalerRegex didn't match %q", tc)
		}
	}
}

func TestJsonUnmarshalerRegex(t *testing.T) {
	tt := []string{
		"func (*github.com/golang-cz/gospeak/data.Person).UnmarshalJSON(data []byte) error",
		"func (*github.com/golang-cz/gospeak/data.Person).UnmarshalJSON(b []byte) (err error)",
	}
	for _, tc := range tt {
		if !jsonUnmarshalerRegex.MatchString(tc) {
			t.Errorf("jsonUnmarshalerRegex didn't match %q", tc)
		}
	}
}
