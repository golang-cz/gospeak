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

	"github.com/webrpc/webrpc/schema"
	"golang.org/x/tools/go/packages"
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
			out: &field{name: "renamed_id", expr: "int64", t: schema.T_Int64, jsonTag: "renamed_id", goName: "ID", goType: "int64"},
		},
		{
			in:  "ID int64 `json:\",string\"`", // string type in JSON
			out: &field{name: "ID", expr: "string", t: schema.T_String, jsonTag: ",string", goName: "ID", goType: "int64"},
		},
		{
			in:  "ID int64 `json:\"id,string\"`", // string type in JSON
			out: &field{name: "id", expr: "string", t: schema.T_String, jsonTag: "id,string", goName: "ID", goType: "int64"},
		},
		{
			in:  "ID int64 `json:\",omitempty\"`", // optional in JSON
			out: &field{name: "ID", expr: "int64", t: schema.T_Int64, jsonTag: ",omitempty", goName: "ID", goType: "int64", optional: true},
		},
		{
			in:  "ID int64 `json:\"id,string,omitempty\"`", // optional string type in JSON
			out: &field{name: "id", expr: "string", t: schema.T_String, jsonTag: "id,string,omitempty", goName: "ID", goType: "int64", optional: true},
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
			in:  "ID uuid.UUID", // uuid implements encoding.TextMarshaler interface, expect string in JSON
			out: &field{name: "ID", expr: "string", t: schema.T_String, goName: "ID", goType: "uuid.UUID", goImport: "github.com/golang-cz/gospeak/uuid"},
		},
		{
			in:  "ID uuid.UUID `json:\",string\"`", // string type in JSON
			out: &field{name: "ID", expr: "string", t: schema.T_String, jsonTag: ",string", goName: "ID", goType: "uuid.UUID", goImport: "github.com/golang-cz/gospeak/uuid"},
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

func testStruct(t *testing.T, inputFields string, want *schema.Type) {
	t.Helper()

	srcCode := fmt.Sprintf(`	package gospeak

	import (
		"context"
		"time"
		
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
	var _ time.Time
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
		t.Fatal(inputFields, fmt.Errorf("type TestStruct not defined"))
	}

	testStruct, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		t.Fatal(inputFields, fmt.Errorf("type TestStruct is %T", obj.Type().Underlying()))
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
		t.Fatal(inputFields, fmt.Errorf("failed to parse struct TestStruct: %w", err))
	}

	if len(p.schema.Types) != 1 {
		t.Fatalf("%s\nexpected one struct type, got %+v", inputFields, p.schema.Types)
	}

	got := p.schema.Types[0]

	if !cmp.Equal(want, got) {
		t.Errorf("%s\n%s\n", inputFields, coloredDiff(want, got))
	}

	return
}

func TestJsonTagRegex(t *testing.T) {
	tt := []struct {
		in  string
		out jsonTag
	}{
		{in: ``},
		{in: `db:"id"`},
		{in: `json:"id"`, out: jsonTag{Name: "id", Value: "id"}},
		{in: `json:"id,whatever"`, out: jsonTag{Name: "id", Value: "id,whatever"}},
		{in: `json:"id,whatever,else"`, out: jsonTag{Name: "id", Value: "id,whatever,else"}},
		{in: `json:"id,string"`, out: jsonTag{Name: "id", Value: "id,string", IsString: true}},
		{in: `json:"id,string,omit"`, out: jsonTag{Name: "id", Value: "id,string,omit", IsString: true}},
		{in: `json:"id,string,omitempty"`, out: jsonTag{Name: "id", Value: "id,string,omitempty", IsString: true, Omitempty: true}},
		{in: `json:"id,omitempty,string"`, out: jsonTag{Name: "id", Value: "id,omitempty,string", IsString: true, Omitempty: true}},
		{in: `json:"id,string,omitempty"`, out: jsonTag{Name: "id", Value: "id,string,omitempty", IsString: true, Omitempty: true}},
		{in: `json:"ID,string,omitempty"`, out: jsonTag{Name: "ID", Value: "ID,string,omitempty", IsString: true, Omitempty: true}},
		{in: `json:"renamed_fieldName99"`, out: jsonTag{Name: "renamed_fieldName99", Value: "renamed_fieldName99"}},
		{in: `xxx:"X X X" json:"id,string" yyy:"Y Y Y"`, out: jsonTag{Name: "id", Value: "id,string", IsString: true}},
		{in: `db:"id,omitempty,pk" json:"id,string"`, out: jsonTag{Name: "id", Value: "id,string", IsString: true}},
		{in: `db:"id,omitempty,pk" json:"External_ID,string,omitempty" someOtherTag:"some,other:value"`, out: jsonTag{Name: "External_ID", Value: "External_ID,string,omitempty", IsString: true, Omitempty: true}},
	}
	for _, tc := range tt {
		jsonTag, ok := getJsonTag(tc.in)
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
