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
				type Enum enum.Int
			`,
			t: schema.T_Int,
			out: []*schema.TypeField{
				// TODO: webrpc name/value looks to be reversed
				{Name: "approved", TypeExtra: schema.TypeExtra{Value: "0"}},
				{Name: "pending", TypeExtra: schema.TypeExtra{Value: "1"}},
				{Name: "closed", TypeExtra: schema.TypeExtra{Value: "2"}},
				{Name: "new", TypeExtra: schema.TypeExtra{Value: "3"}},
			},
		},
		{
			in: `
				// approved
				// pending
				// closed
				// new
				type Enum enum.Uint64
			`,
			t: schema.T_Uint64,
			out: []*schema.TypeField{
				{Name: "approved", TypeExtra: schema.TypeExtra{Value: "0"}},
				{Name: "pending", TypeExtra: schema.TypeExtra{Value: "1"}},
				{Name: "closed", TypeExtra: schema.TypeExtra{Value: "2"}},
				{Name: "new", TypeExtra: schema.TypeExtra{Value: "3"}},
			},
		},
		{
			// TODO: Can we also support "cs-CZ"?
			in: `
				// en    = 0
				// zh_CN = 1
				// zh_HK = 2
				// da_DK = 3
				// nl_NL = 4
				// en_AU = 5
				// en_CA = 6
				// en_GB = 7
				// fi_FI = 8
				// fr_CA = 9
				// fr_FR = 10
				// de_DE = 11
				// id_ID = 12
				// it_IT = 13
				// ja_JP = 14
				// ko_KR = 15
				// ms_MY = 16
				// nb_NO = 17
				// pl_PL = 18
				// pt_BR = 19
				// ru_RU = 20
				// es_MX = 21
				// sv_SE = 22
				// tl_PH = 23
				// th_TH = 24
				// vi_VN = 25
				// el_GR = 29
				// es_ES = 30
				// hi_IN = 31
				// hu_HU = 32
				// sk_SK = 33
				// tr_TR = 34
				// cs_CZ = 35
				// en_US = 36
				// ro_RO = 37
				// pt_PT = 38
				// zh_TW = 39
				// es_US = 40
				// hr_HR = 41
				// zh_SG = 42
				// ar_SA = 43
				// he_IL = 44
				// ca_ES = 45
				type Enum enum.Uint32
			`,
			t: schema.T_Uint32,
			out: []*schema.TypeField{
				{Name: "en", TypeExtra: schema.TypeExtra{Value: "0"}},
				{Name: "zh_CN", TypeExtra: schema.TypeExtra{Value: "1"}},
				{Name: "zh_HK", TypeExtra: schema.TypeExtra{Value: "2"}},
				{Name: "da_DK", TypeExtra: schema.TypeExtra{Value: "3"}},
				{Name: "nl_NL", TypeExtra: schema.TypeExtra{Value: "4"}},
				{Name: "en_AU", TypeExtra: schema.TypeExtra{Value: "5"}},
				{Name: "en_CA", TypeExtra: schema.TypeExtra{Value: "6"}},
				{Name: "en_GB", TypeExtra: schema.TypeExtra{Value: "7"}},
				{Name: "fi_FI", TypeExtra: schema.TypeExtra{Value: "8"}},
				{Name: "fr_CA", TypeExtra: schema.TypeExtra{Value: "9"}},
				{Name: "fr_FR", TypeExtra: schema.TypeExtra{Value: "10"}},
				{Name: "de_DE", TypeExtra: schema.TypeExtra{Value: "11"}},
				{Name: "id_ID", TypeExtra: schema.TypeExtra{Value: "12"}},
				{Name: "it_IT", TypeExtra: schema.TypeExtra{Value: "13"}},
				{Name: "ja_JP", TypeExtra: schema.TypeExtra{Value: "14"}},
				{Name: "ko_KR", TypeExtra: schema.TypeExtra{Value: "15"}},
				{Name: "ms_MY", TypeExtra: schema.TypeExtra{Value: "16"}},
				{Name: "nb_NO", TypeExtra: schema.TypeExtra{Value: "17"}},
				{Name: "pl_PL", TypeExtra: schema.TypeExtra{Value: "18"}},
				{Name: "pt_BR", TypeExtra: schema.TypeExtra{Value: "19"}},
				{Name: "ru_RU", TypeExtra: schema.TypeExtra{Value: "20"}},
				{Name: "es_MX", TypeExtra: schema.TypeExtra{Value: "21"}},
				{Name: "sv_SE", TypeExtra: schema.TypeExtra{Value: "22"}},
				{Name: "tl_PH", TypeExtra: schema.TypeExtra{Value: "23"}},
				{Name: "th_TH", TypeExtra: schema.TypeExtra{Value: "24"}},
				{Name: "vi_VN", TypeExtra: schema.TypeExtra{Value: "25"}},
				{Name: "el_GR", TypeExtra: schema.TypeExtra{Value: "29"}},
				{Name: "es_ES", TypeExtra: schema.TypeExtra{Value: "30"}},
				{Name: "hi_IN", TypeExtra: schema.TypeExtra{Value: "31"}},
				{Name: "hu_HU", TypeExtra: schema.TypeExtra{Value: "32"}},
				{Name: "sk_SK", TypeExtra: schema.TypeExtra{Value: "33"}},
				{Name: "tr_TR", TypeExtra: schema.TypeExtra{Value: "34"}},
				{Name: "cs_CZ", TypeExtra: schema.TypeExtra{Value: "35"}},
				{Name: "en_US", TypeExtra: schema.TypeExtra{Value: "36"}},
				{Name: "ro_RO", TypeExtra: schema.TypeExtra{Value: "37"}},
				{Name: "pt_PT", TypeExtra: schema.TypeExtra{Value: "38"}},
				{Name: "zh_TW", TypeExtra: schema.TypeExtra{Value: "39"}},
				{Name: "es_US", TypeExtra: schema.TypeExtra{Value: "40"}},
				{Name: "hr_HR", TypeExtra: schema.TypeExtra{Value: "41"}},
				{Name: "zh_SG", TypeExtra: schema.TypeExtra{Value: "42"}},
				{Name: "ar_SA", TypeExtra: schema.TypeExtra{Value: "43"}},
				{Name: "he_IL", TypeExtra: schema.TypeExtra{Value: "44"}},
				{Name: "ca_ES", TypeExtra: schema.TypeExtra{Value: "45"}},
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
				Enum Enum
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
			Name: "Enum",
			Type: &schema.VarType{
				Expr: tc.t.String(),
				Type: tc.t,
			},
			Fields: tc.out,
		}

		var got *schema.Type
		for _, schemaType := range p.Schema.Types {
			if schemaType.Name == "Enum" {
				got = schemaType
			}
		}

		if !cmp.Equal(want, got) {
			t.Errorf("%s\n%s\n", tc.in, coloredDiff(want, got))
		}

	}
}
