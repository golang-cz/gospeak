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

	srcCode := fmt.Sprintf(`package test

	import (
		"context"
		"time"

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

	// Ensure all the imports are used.
	var _ time.Time
	var _ uuid.UUID
	var _ Number
	var _ Locale
	`, inputFields)

	p, err := testStructParser(srcCode)
	if err != nil {
		t.Fatal(fmt.Errorf("error parsing: %q: %w", inputFields, err))
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

func testParser(srcCode string) (*parser.Parser, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("getting working directory: %w", err)
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
		return nil, fmt.Errorf("error loading Go packages: %v\n%s", err, prefixLinesWithLineNumber(srcCode))
	}

	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			return nil, fmt.Errorf("%v\n%s", spew.Sdump(pkg.Errors), prefixLinesWithLineNumber(srcCode))
		}
	}

	if len(pkgs) != 2 {
		return nil, fmt.Errorf("expected 2 Go packages, got %v\n%s", len(pkgs), spew.Sdump(pkgs))
	}

	pkg := pkgs[0]

	_, _ = gospeak.CollectInterfaces(pkg)

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

	return p, nil
}

// Parses code with TestStruct type.
func testStructParser(srcCode string) (*parser.Parser, error) {
	p, err := testParser(srcCode)
	if err != nil {
		return nil, err
	}

	scope := p.Pkg.Types.Scope()

	obj := scope.Lookup("TestStruct")
	if obj == nil {
		return nil, fmt.Errorf("type TestStruct not defined")
	}

	testStruct, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		return nil, fmt.Errorf("type TestStruct is %T", obj.Type().Underlying())
	}

	_, err = p.ParseStruct("TestStruct", testStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to parse struct TestStruct: %w", err)
	}

	return p, nil
}
