package test

import (
	"testing"

	"github.com/webrpc/webrpc/schema"
)

func TestStructFieldJsonTags(t *testing.T) {
	t.Parallel()

	type field struct {
		name     string
		expr     string
		t        schema.CoreType
		jsonTag  string
		goName   string
		goType   string
		goImport string
		optional bool
	}

	tt := []struct {
		in  string
		out *field
	}{
		{
			in:  "ID int64", // default name
			out: &field{name: "ID", expr: "int64", t: schema.T_Int64, goName: "ID", goType: "int64"},
		},
		{
			in:  "id int64", // unexported field
			out: nil,
		},
		{
			in:  "ID int64 `json:\"-\"`", // ignored in JSON
			out: nil,
		},
		{
			in:  "ID *int64", // optional
			out: &field{name: "ID", expr: "int64", t: schema.T_Int64, goName: "ID", goType: "*int64", optional: true},
		},
		{
			in:  "ID int64 `json:\"renamed_id\"`", // renamed in JSON
			out: &field{name: "renamed_id", expr: "int64", t: schema.T_Int64, goName: "ID", goType: "int64", jsonTag: "renamed_id"},
		},
		{
			in:  "ID int64 `json:\",string\"`", // string type in JSON
			out: &field{name: "ID", expr: "string", t: schema.T_String, goName: "ID", goType: "int64", jsonTag: ",string"},
		},
		{
			in:  "ID int64 `json:\"id,string\"`", // string type in JSON
			out: &field{name: "id", expr: "string", t: schema.T_String, goName: "ID", goType: "int64", jsonTag: "id,string"},
		},
		{
			in:  "ID int64 `json:\",omitempty\"`", // optional in JSON
			out: &field{name: "ID", expr: "int64", t: schema.T_Int64, goName: "ID", goType: "int64", jsonTag: ",omitempty", optional: true},
		},
		{
			in:  "ID int64 `json:\"id,string,omitempty\"`", // optional string type in JSON
			out: &field{name: "id", expr: "string", t: schema.T_String, goName: "ID", goType: "int64", jsonTag: "id,string,omitempty", optional: true},
		},
		{
			in:  "CreatedAt time.Time",
			out: &field{name: "CreatedAt", expr: "timestamp", t: schema.T_Timestamp, goName: "CreatedAt", goType: "time.Time"},
		},
		{
			in:  "DeletedAt *time.Time",
			out: &field{name: "DeletedAt", expr: "timestamp", t: schema.T_Timestamp, goName: "DeletedAt", goType: "*time.Time", optional: true},
		},
		{
			in:  "Number Number",
			out: &field{name: "Number", expr: "int", t: schema.T_Int, goName: "Number", goType: "Number"},
		},
		{
			in:  "NumberString Number `json:\",string\"`",
			out: &field{name: "NumberString", expr: "string", t: schema.T_String, goName: "NumberString", goType: "Number", jsonTag: ",string"},
		},
		{
			in:  "LocaleString Locale",
			out: &field{name: "LocaleString", expr: "string", t: schema.T_String, goName: "LocaleString", goType: "Locale"},
		},
		{
			in:  "ID uuid.UUID", // uuid implements encoding.TextMarshaler interface, expect string in JSON
			out: &field{name: "ID", expr: "string", t: schema.T_String, goName: "ID", goType: "uuid.UUID", goImport: "github.com/golang-cz/gospeak/internal/parser/test/uuid"},
		},
		{
			in:  "ID uuid.UUID `json:\",string\"`", // string type in JSON
			out: &field{name: "ID", expr: "string", t: schema.T_String, jsonTag: ",string", goName: "ID", goType: "uuid.UUID", goImport: "github.com/golang-cz/gospeak/internal/parser/test/uuid"},
		},
	}

	for _, tc := range tt {
		var fields []*schema.TypeField
		if tc.out != nil {
			fields = []*schema.TypeField{
				&schema.TypeField{
					Name: tc.out.name,
					Type: &schema.VarType{
						Expr: tc.out.expr,
						Type: tc.out.t,
					},
					TypeExtra: schema.TypeExtra{
						Optional: tc.out.optional,
						Meta: []schema.TypeFieldMeta{
							{"go.field.name": tc.out.goName},
							{"go.field.type": tc.out.goType},
						},
					},
				},
			}
			if tc.out.goImport != "" {
				fields[0].TypeExtra.Meta = append(fields[0].TypeExtra.Meta, schema.TypeFieldMeta{"go.type.import": tc.out.goImport})
			}
			if tc.out.jsonTag != "" {
				fields[0].TypeExtra.Meta = append(fields[0].TypeExtra.Meta, schema.TypeFieldMeta{"go.tag.json": tc.out.jsonTag})
			}
		}

		testStruct(t,
			tc.in,
			&schema.Type{
				Kind:   "struct",
				Name:   "TestStruct",
				Fields: fields,
			},
		)
	}
}

func TestStructSliceField(t *testing.T) {
	t.Parallel()

	type field struct {
		name     string
		elemExpr string          // element
		elemT    schema.CoreType // element
		goName   string
		goType   string
		optional bool
		imports  []string
	}

	tt := []struct {
		in  string
		out *field
	}{
		{
			in:  "ID []int64",
			out: &field{name: "ID", elemExpr: "int64", elemT: schema.T_Int64, goName: "ID", goType: "[]int64"},
		},
	}

	for _, tc := range tt {
		testStruct(t,
			tc.in,
			&schema.Type{
				Kind: "struct",
				Name: "TestStruct",
				Fields: []*schema.TypeField{
					&schema.TypeField{
						Name: tc.out.name,
						Type: &schema.VarType{
							Expr: "[]" + tc.out.elemExpr,
							Type: schema.T_List,
							List: &schema.VarListType{
								Elem: &schema.VarType{
									Expr: tc.out.elemExpr,
									Type: tc.out.elemT,
								},
							},
						},
						TypeExtra: schema.TypeExtra{
							Optional: tc.out.optional,
							Meta: []schema.TypeFieldMeta{
								{"go.field.name": tc.out.goName},
								{"go.field.type": tc.out.goType},
							},
						},
					},
				},
			},
		)
	}
}
