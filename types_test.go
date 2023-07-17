package gospeak

import (
	"fmt"
	"go/types"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"github.com/webrpc/webrpc/schema"
	"golang.org/x/tools/go/packages"
)

func TestStructFieldJsonTags(t *testing.T) {
	t.Parallel()

	type webrpcType struct {
		name        string
		expr        string
		t           schema.CoreType
		goFieldName string
		goFieldType string
		optional    bool
	}

	tt := []struct {
		in  string
		out *webrpcType
	}{
		{
			in:  "ID int64", // default name
			out: &webrpcType{name: "ID", expr: "int64", t: schema.T_Int64, goFieldName: "ID", goFieldType: "int64"},
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
			out: &webrpcType{name: "ID", expr: "*int64", t: schema.T_Int64, goFieldName: "ID", goFieldType: "*int64", optional: true},
		},
		{
			in:  "ID int64 `json:\"renamed_id\"`", // renamed in JSON
			out: &webrpcType{name: "renamed_id", expr: "int64", t: schema.T_Int64, goFieldName: "ID", goFieldType: "int64"},
		},
		{
			in:  "ID int64 `json:\",string\"`", // string type in JSON
			out: &webrpcType{name: "ID", expr: "string", t: schema.T_String, goFieldName: "ID", goFieldType: "int64"},
		},
		{
			in:  "ID int64 `json:\"id,string\"`", // string type in JSON
			out: &webrpcType{name: "id", expr: "string", t: schema.T_String, goFieldName: "ID", goFieldType: "int64"},
		},
		{
			in:  "ID int64 `json:\",omitempty\"`", // optional in JSON
			out: &webrpcType{name: "ID", expr: "int64", t: schema.T_Int64, goFieldName: "ID", goFieldType: "int64", optional: true},
		},
		{
			in:  "ID int64 `json:\"id,string,omitempty\"`", // optional string type in JSON
			out: &webrpcType{name: "id", expr: "string", t: schema.T_String, goFieldName: "ID", goFieldType: "int64", optional: true},
		},
		{
			in:  "ID uuid.UUID", // uuid implements encoding.TextMarshaler interface, expect string in JSON
			out: &webrpcType{name: "ID", expr: "uuid.UUID", t: schema.T_String, goFieldName: "ID", goFieldType: "uuid.UUID"},
		},
		{
			in:  "ID uuid.UUID `json:\",string\"`", // string type in JSON
			out: &webrpcType{name: "ID", expr: "string", t: schema.T_String, goFieldName: "ID", goFieldType: "uuid.UUID"},
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
							{"go.field.name": tc.out.goFieldName},
							{"go.field.type": tc.out.goFieldType},
						},
					},
				},
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

	type webrpcType struct {
		name        string
		elemExpr    string          // element
		elemT       schema.CoreType // element
		goFieldName string
		goFieldType string
		optional    bool
		imports     []string
	}

	tt := []struct {
		in  string
		out *webrpcType
	}{
		{
			in:  "ID []int64",
			out: &webrpcType{name: "ID", elemExpr: "int64", elemT: schema.T_Int64, goFieldName: "ID", goFieldType: "[]int64"},
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
								{"go.field.name": tc.out.goFieldName},
								{"go.field.type": tc.out.goFieldType},
							},
						},
					},
				},
			},
		)
	}
}

// TODO: Test extra struct tags.
// {
// 	in:  "ID int64 `db:\"id\", json:\"id\"`", // test extra tags
// 	out: &webrpcType{name: "id", expr: "int64", t: schema.T_Int64, goFieldName: "ID", goFieldType: "int64"},
// },

func testStruct(t *testing.T, inputFields string, want *schema.Type) {
	t.Helper()

	srcCode := fmt.Sprintf(`	package gospeak

	import (
		"context"
		
		"github.com/golang-cz/gospeak/uuid"
	)

	type TestStruct struct {
		%s
	}
	
	//go:webrpc json -out=/dev/null
	type TestAPI interface{
		TestStruct(ctx context.Context) (tst *TestStruct, err error)
	}


	// Ensure all the imports are used.
	var _ uuid.UUID
	`, inputFields)

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	package1Path := filepath.Join(wd, "proto.go")
	package2Path := filepath.Join(wd, "uuid/uuid.go")

	cfg := &packages.Config{
		Dir:  wd,
		Mode: packages.NeedName | packages.NeedSyntax | packages.NeedTypes | packages.LoadImports,
		Overlay: map[string][]byte{
			package1Path: []byte(srcCode),
			package2Path: []byte(`
				package uuid

				type UUID [16]byte

				// MarshalText implements encoding.TextMarshaler.
				func (uuid UUID) MarshalText() ([]byte, error) {
					return []byte{}, nil
				}

				// UnmarshalText implements encoding.TextUnmarshaler.
				func (uuid *UUID) UnmarshalText(data []byte) error {
					return nil
				}
			`),
		},
	}

	pkgs, err := packages.Load(cfg, "file="+package1Path, "file="+package2Path)
	if err != nil {
		t.Fatal(inputFields, fmt.Errorf("loading Go packages: %w", err))
	}

	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			t.Log(spew.Sdump(pkg.Errors))
		}
	}

	if len(pkgs) != 2 {
		t.Fatal(inputFields, fmt.Errorf("expected 2 Go packages, got %v\n%s", len(pkgs), spew.Sdump(pkgs)))
	}

	pkg := pkgs[0]

	if len(pkg.Errors) > 0 {
		t.Fatal(inputFields, fmt.Sprintf("%+v\n%s", pkg.Errors, prefixLinesWithLineNumber(srcCode)))
	}

	_, _ = collectInterfaces(pkg)

	scope := pkg.Types.Scope()

	obj := scope.Lookup("TestStruct")
	if obj == nil {
		t.Fatal(inputFields, errors.Errorf("type TestStruct not defined"))
	}

	testStruct, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		t.Fatal(inputFields, errors.Errorf("type TestStruct is %T", obj.Type().Underlying()))
	}

	p := &parser{
		schema: &schema.WebRPCSchema{
			WebrpcVersion: "v1",
			SchemaName:    "TestAPI",
			SchemaVersion: "",
		},
		schemaPkgName:   pkg.Name,
		parsedTypes:     map[types.Type]*schema.VarType{},
		parsedTypeNames: map[string]struct{}{},

		// TODO: Change this to map[*types.Package]string so we can rename duplicated pkgs?
		importedPaths: map[string]struct{}{
			// Initial schema file's package name artificially set by golang.org/x/tools/go/packages.
			"command-line-arguments": {},
		},
	}

	_, err = p.parseStruct("TestStruct", testStruct)
	if err != nil {
		t.Fatal(inputFields, errors.Wrapf(err, "failed to parse struct TestStruct"))
	}

	if len(p.schema.Types) != 1 {
		t.Fatalf(inputFields, "expected one struct type, got %+v", p.schema.Types)
	}

	got := p.schema.Types[0]

	if !cmp.Equal(want, got) {
		t.Errorf("%s\n%s\n", inputFields, coloredDiff(want, got))
	}

	return
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

// exported := packagestest.Export(t, exporter, []packagestest.Module{{
// 	Name: "fake",
// 	Files: map[string]interface{}{
// 		"api.go":      "package foo\nfunc g(){}\n",
// 	},
// 	Overlay: map[string][]byte{
// 		"a.go":      []byte("package foox\nfunc g(){}\n"),
// 		"a_test.go": []byte("package foox\nfunc f(){}\n"),
// 	},
// }})
// defer exported.Cleanup()

func coloredDiff(x, y interface{}, opts ...cmp.Option) string {
	escapeCode := func(code int) string {
		return fmt.Sprintf("\x1b[%dm", code)
	}
	diff := cmp.Diff(x, y, opts...)
	if diff == "" {
		return ""
	}
	ss := strings.Split(diff, "\n")
	for i, s := range ss {
		switch {
		case strings.HasPrefix(s, "-"):
			ss[i] = escapeCode(31) + s + escapeCode(0)
		case strings.HasPrefix(s, "+"):
			ss[i] = escapeCode(32) + s + escapeCode(0)
		}
	}
	return strings.Join(ss, "\n")
}

func prefixLinesWithLineNumber(input string) string {
	lines := strings.Split(input, "\n")
	for i := range lines {
		lines[i] = fmt.Sprintf("%d: %s", i+1, lines[i])
	}
	return strings.Join(lines, "\n")
}
