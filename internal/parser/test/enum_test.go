package test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/webrpc/webrpc/schema"
)

func TestStructFieldEnum(t *testing.T) {
	t.Parallel()

	tt := []struct {
		in  string
		t   schema.CoreType
		out []*schema.TypeField
	}{
		{
			in: `
				// approved = 0
				// pending  = 1
				// closed   = 2
				// new      = 3
				type Status enum.Int
			`,
			t: schema.T_Int,
			out: []*schema.TypeField{
				// TODO: webrpc name/value looks to be reversed
				&schema.TypeField{Name: "approved", TypeExtra: schema.TypeExtra{Value: "0"}},
				&schema.TypeField{Name: "pending", TypeExtra: schema.TypeExtra{Value: "1"}},
				&schema.TypeField{Name: "closed", TypeExtra: schema.TypeExtra{Value: "2"}},
				&schema.TypeField{Name: "new", TypeExtra: schema.TypeExtra{Value: "3"}},
			},
		},
		{
			in: `
				// approved
				// pending
				// closed
				// new
				type Status enum.Uint64
			`,
			t: schema.T_Uint64,
			out: []*schema.TypeField{
				&schema.TypeField{Name: "approved", TypeExtra: schema.TypeExtra{Value: "0"}},
				&schema.TypeField{Name: "pending", TypeExtra: schema.TypeExtra{Value: "1"}},
				&schema.TypeField{Name: "closed", TypeExtra: schema.TypeExtra{Value: "2"}},
				&schema.TypeField{Name: "new", TypeExtra: schema.TypeExtra{Value: "3"}},
			},
		},
	}

	for _, tc := range tt {
		srcCode := fmt.Sprintf(`package test
		
			import (
				"context"
			
				"github.com/golang-cz/gospeak/enum"
			)

			%s
		
			type TestStruct struct {
				Status Status
			}
			
			//go:webrpc json -out=/dev/null
			type TestAPI interface{
				Test(ctx context.Context) (tst *TestStruct, err error)
			}
			`, tc.in)

		p, err := testParser(srcCode)
		if err != nil {
			t.Fatal(fmt.Errorf("parsing: %w", err))
		}

		if err := p.CollectEnums(); err != nil {
			t.Fatalf("collecting enums: %v", err)
		}

		want := &schema.Type{
			Kind: schema.TypeKind_Enum,
			Name: "Status",
			Type: &schema.VarType{
				Expr: tc.t.String(),
				Type: tc.t,
			},
			Fields: tc.out,
		}

		var got *schema.Type
		for _, schemaType := range p.Schema.Types {
			if schemaType.Name == "Status" {
				got = schemaType
			}
		}

		if !cmp.Equal(want, got) {
			t.Errorf("%s\n%s\n", tc.in, coloredDiff(want, got))
		}

	}
}
