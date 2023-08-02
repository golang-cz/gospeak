package test

import (
	"testing"

	"github.com/webrpc/webrpc/schema"
)

func TestStructFieldEnum(t *testing.T) {
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
			in:  "Status Status", // enum.Status
			out: &field{name: "Status", expr: "Status", t: schema.T_Struct, goName: "Status", goType: "Status"},
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
						Struct: &schema.VarStructType{
							Name: tc.out.name,
							Type: &schema.Type{
								Kind: schema.TypeKind_Enum,
								Name: tc.out.name,
								Type: &schema.VarType{
									Expr: "int",
									Type: schema.T_Int,
								},
								Fields: []*schema.TypeField{
									// 0 = approved
									// 1 = pending
									// 2 = closed
									// 3 = new
									&schema.TypeField{Name: "approved", TypeExtra: schema.TypeExtra{Value: "0"}},
									&schema.TypeField{Name: "pending", TypeExtra: schema.TypeExtra{Value: "1"}},
									&schema.TypeField{Name: "closed", TypeExtra: schema.TypeExtra{Value: "2"}},
									&schema.TypeField{Name: "new", TypeExtra: schema.TypeExtra{Value: "3"}},
								},
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
