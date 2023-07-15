package gospeak

import (
	"fmt"
	"go/types"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"github.com/webrpc/webrpc/schema"
	"golang.org/x/tools/go/packages"
)

func TestStructFields(t *testing.T) {
	type webrpcType struct {
		name        string
		expr        string
		t           schema.CoreType
		goFieldName string
		goFieldType string
	}

	tt := []struct {
		in  string
		out *webrpcType
	}{
		{
			in:  "ID int64",
			out: &webrpcType{name: "ID", expr: "int64", t: schema.T_Int64, goFieldName: "ID", goFieldType: "int64"},
		},
		{
			in:  "ID int64 `json:\"id\"`",
			out: &webrpcType{name: "id", expr: "int64", t: schema.T_Int64, goFieldName: "ID", goFieldType: "int64"},
		},
		{
			in:  "ID int64 `json:\"id,string\"`",
			out: &webrpcType{name: "id", expr: "string", t: schema.T_String, goFieldName: "ID", goFieldType: "int64"},
		},
		{
			in:  "ID int64 `json:\",string\"`",
			out: &webrpcType{name: "ID", expr: "string", t: schema.T_String, goFieldName: "ID", goFieldType: "int64"},
		},
		{
			in:  "id int64", // unexported field
			out: nil,
		},
		{
			in:  "ID int64 `json:\"-\"`", // ignored by json:"-" tag
			out: nil,
		},
		{
			in:  "ID int64 `db:\"id\", json:\"id\"`", // test more tags
			out: &webrpcType{name: "id", expr: "int64", t: schema.T_Int64, goFieldName: "ID", goFieldType: "int64"},
		},
		{
			in:  "Id int32",
			out: &webrpcType{name: "Id", expr: "int32", t: schema.T_Int32, goFieldName: "Id", goFieldType: "int32"},
		},
		{
			in:  "Name string",
			out: &webrpcType{name: "Name", expr: "string", t: schema.T_String, goFieldName: "Name", goFieldType: "string"},
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
						Meta: []schema.TypeFieldMeta{
							{"go.field.name": tc.out.goFieldName},
							{"go.field.type": tc.out.goFieldType},
						},
					},
				}}
		}

		testStruct(t,
			tc.in,
			&schema.Type{
				Kind:   "struct",
				Name:   "TestStruct",
				Fields: fields,
			})
	}
}

func testStruct(t *testing.T, input string, want *schema.Type) {
	t.Helper()

	f, _ := os.CreateTemp("", "")
	defer os.Remove(f.Name()) // clean up
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	cfg := &packages.Config{
		Dir: ".",
		Mode: packages.NeedName | packages.NeedModule |
			packages.NeedImports | packages.NeedDeps |
			packages.NeedTypes | packages.NeedTypesInfo |
			packages.NeedFiles | packages.NeedCompiledGoFiles |
			packages.NeedSyntax | packages.LoadSyntax,
		Overlay: map[string][]byte{
			"/api.go": []byte(fmt.Sprintf(`
				package proto

				import "context"

				type TestedStruct struct {
					%s
				}
				
				//go:webrpc json -out=%v
				type ExampleAPI interface{
					TestStruct(ctx context.Context) (tst *TestedStruct, err error)
				}
				`, input, f.Name())),
		},
	}

	pkgs, err := packages.Load(cfg, "/api.go")
	if err != nil {
		t.Fatal(fmt.Errorf("loading Go packages"))
	}
	if len(pkgs) != 1 {
		t.Fatal(fmt.Errorf("expected one Go package, got %v", len(pkgs)))
	}
	pkg := pkgs[0]

	if len(pkg.Errors) > 0 {
		t.Fatal(fmt.Sprintf("%+v", pkg.Errors))
	}

	_, _ = collectInterfaces(pkg)

	scope := pkg.Types.Scope()

	obj := scope.Lookup("TestedStruct")
	if obj == nil {
		t.Fatal(errors.Errorf("type TestedStruct not defined"))
	}

	testedStruct, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		t.Fatal(errors.Errorf("type TestedStruct is %T", obj.Type().Underlying()))
	}

	p := &parser{
		schema: &schema.WebRPCSchema{
			WebrpcVersion: "v1",
			SchemaName:    "ExampleAPI",
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

	_, err = p.parseStruct("TestStruct", testedStruct)
	if err != nil {
		t.Fatal(errors.Wrapf(err, "failed to parse struct TestStruct"))
	}

	if len(p.schema.Types) != 1 {
		t.Fatalf("expected one struct type, got %+v", p.schema.Types)
	}

	got := p.schema.Types[0]

	if !cmp.Equal(want, got) {
		t.Errorf("%s\n%s\n", input, cmp.Diff(want, got))
	}

	return
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
