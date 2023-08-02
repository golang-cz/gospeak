package test

import (
	"fmt"
	"go/types"
	"os"
	"path/filepath"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang-cz/gospeak"
	"github.com/golang-cz/gospeak/internal/parser"
	"github.com/google/go-cmp/cmp"
	"github.com/webrpc/webrpc/schema"
	"golang.org/x/tools/go/packages"
)

func testStruct(t *testing.T, inputFields string, want *schema.Type) {
	t.Helper()

	srcCode := fmt.Sprintf(`	package test

	import (
		"context"
		"time"
		
		"github.com/golang-cz/gospeak"
		"github.com/golang-cz/gospeak/internal/parser/test/uuid"
	)

	type TestStruct struct {
		%s
	}
	
	//go:webrpc json -out=/dev/null
	type TestAPI interface{
		TestStruct(ctx context.Context) (tst *TestStruct, err error)
	}

	type Number int // should be number over JSON

	type Locale int // implements MarshalText(), should be string over JSON

	// MarshalText implements encoding.TextMarshaler.
	func (locale Locale) MarshalText() ([]byte, error) {
		return []byte{}, nil
	}

	// UnmarshalText implements encoding.TextUnmarshaler.
	func (locale *Locale) UnmarshalText(data []byte) error {
		return nil
	}

	// approved = 0
	// pending  = 1
	// closed   = 2
	// new      = 3
	type Status gospeak.Enum[int]

	// Ensure all the imports are used.
	var _ time.Time
	var _ uuid.UUID
	var _ Number
	var _ Locale
	var _ Status
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
		Mode: packages.NeedName | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedImports,
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

	_, _ = gospeak.CollectInterfaces(pkg)

	scope := pkg.Types.Scope()

	obj := scope.Lookup("TestStruct")
	if obj == nil {
		t.Fatal(inputFields, fmt.Errorf("type TestStruct not defined"))
	}

	testStruct, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		t.Fatal(inputFields, fmt.Errorf("type TestStruct is %T", obj.Type().Underlying()))
	}

	p := &parser.Parser{
		Schema: &schema.WebRPCSchema{
			WebrpcVersion: "v1",
			SchemaName:    "TestAPI",
			SchemaVersion: "v0.0.1",
		},
		SchemaPkgName:   pkg.PkgPath,
		ParsedTypes:     map[types.Type]*schema.VarType{},
		ParsedTypeNames: map[string]struct{}{},
		Pkg:             pkg,

		// TODO: Change this to map[*types.Package]string so we can rename duplicated pkgs?
		ImportedPaths: map[string]struct{}{
			// Initial Schema file's package name artificially set by golang.org/x/tools/go/packages.
			"command-line-arguments": {},
		},
	}

	if err := p.CollectEnums(); err != nil {
		t.Fatal(inputFields, fmt.Errorf("collecting enums: %w", err))
	}

	_, err = p.ParseStruct("TestStruct", testStruct)
	if err != nil {
		t.Fatal(inputFields, fmt.Errorf("failed to parse struct TestStruct: %w", err))
	}

	for _, got := range p.Schema.Types {
		if got.Name != "TestStruct" {
			continue
		}

		if !cmp.Equal(want, got) {
			t.Errorf("%s\n%s\n", inputFields, coloredDiff(want, got))
		}

		return // success
	}

	t.Fatalf("%s\nexpected one struct type, got %s", inputFields, spew.Sdump(p.Schema.Types))
	return
}
