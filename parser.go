package go2webrpc

import (
	"go/types"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/webrpc/webrpc/schema"
	"golang.org/x/tools/go/packages"
)

// Parse Go source file or package folder and return WebRPC schema.
func Parse(filePath string, goInterfaceName string) (*schema.WebRPCSchema, error) {
	file, err := os.Stat(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open -schema=%q", filePath)
	}
	if file.Mode().IsRegular() {
		// Parse all files in the given schema file's directory, so the parser can see all the pkg types.
		filePath = path.Dir(filePath)
	}

	cfg := &packages.Config{
		Dir:  filePath,
		Mode: packages.NeedName | packages.NeedImports | packages.NeedTypes | packages.NeedFiles | packages.NeedDeps | packages.NeedSyntax,
	}

	schemaPkg, err := packages.Load(cfg, filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load packages")
	}
	if len(schemaPkg) != 1 {
		return nil, errors.Errorf("failed to load initial package (len=%v)", len(schemaPkg))
	}

	scope := schemaPkg[0].Types.Scope()

	obj := scope.Lookup(goInterfaceName)
	if obj == nil {
		return nil, errors.Errorf("interface %q not found", goInterfaceName)
	}

	iface, ok := obj.Type().Underlying().(*types.Interface)
	if !ok {
		return nil, errors.Errorf("%q is not an interface", goInterfaceName)
	}

	p := &parser{
		schema:          &schema.WebRPCSchema{},
		schemaPkgName:   schemaPkg[0].Name,
		parsedTypes:     map[types.Type]*schema.VarType{},
		parsedTypeNames: map[string]struct{}{},

		// TODO: Change this to map[*types.Package]string so we can rename duplicated pkgs?
		importedPaths: map[string]struct{}{
			// Initial schema file's package name artificially set by golang.org/x/tools/go/packages.
			"command-line-arguments": {},
		},
	}

	if err := p.parseInterfaceMethods(iface, goInterfaceName); err != nil {
		return nil, errors.Wrapf(err, "failed to parser interface %q", goInterfaceName)
	}

	return p.schema, nil
}

// Parses Go source file and returns WebRPC schema.
//
// This parser was designed to run sequentially, without any concurrency, so we can leverage
// maps to cache the already parsed types, while not having to deal with sync primitives.
type parser struct {
	schema *schema.WebRPCSchema

	// Cache for already parsed types, to improve performance & so we can traverse circular dependencies.
	parsedTypes     map[types.Type]*schema.VarType
	parsedTypeNames map[string]struct{}

	inlineMode    bool // When traversing `json:",inline"`, we don't want to store the struct type as WebRPC message.
	importedPaths map[string]struct{}

	schemaPkgName string // Schema file's package name.
}
