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
				type Status gospeak.Enum[int]
			`,
			t: schema.T_Int,
			out: []*schema.TypeField{
				// TODO: webrpc name/value looks to be reversed
				&schema.TypeField{Name: "0", TypeExtra: schema.TypeExtra{Value: "approved"}},
				&schema.TypeField{Name: "1", TypeExtra: schema.TypeExtra{Value: "pending"}},
				&schema.TypeField{Name: "2", TypeExtra: schema.TypeExtra{Value: "closed"}},
				&schema.TypeField{Name: "3", TypeExtra: schema.TypeExtra{Value: "new"}},
			},
		},
	}

	for _, tc := range tt {
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
