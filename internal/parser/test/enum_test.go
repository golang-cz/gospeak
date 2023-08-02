package test

import (
	"fmt"
	"testing"

	"github.com/webrpc/webrpc/schema"
)

func TestStructFieldEnum(t *testing.T) {
	t.Parallel()

	type enum struct {
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
		out *enum
	}{
		{
			in: `
				// approved = 0
				// pending  = 1
				// closed   = 2
				// new      = 3
				type Status gospeak.Enum[int]
			`,
			out: &enum{name: "Status", expr: "Status", t: schema.T_Struct, goName: "Status", goType: "Status"},
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

		srcCode := fmt.Sprintf(`package test
		
			import (
				"context"
			
				"github.com/golang-cz/gospeak"
			)

			%s
		
			type TestStruct struct {
				Status Status
			}
			
			//go:webrpc json -out=/dev/null
			type TestAPI interface{
				TestStruct(ctx context.Context) (tst *TestStruct, err error)
			}

			// Ensure all the imports are used.
			//var _ = gospeak.Enum[int]{}
			//var _ = Status{}
			`, tc.in)

		p, err := testParser(srcCode)
		if err != nil {
			t.Fatal(fmt.Errorf("parsing: %w", err))
		}

		if err := p.CollectEnums(); err != nil {
			t.Fatalf("collecting enums: %v", err)
		}
	}
}
