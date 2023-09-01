package parser

import (
	"go/types"

	"github.com/webrpc/webrpc/schema"
	"golang.org/x/tools/go/packages"
)

// Parses Go source file and returns WebRPC Schema.
//
// This Parser was designed to run sequentially, without any concurrency, so we can leverage
// maps to cache the already parsed types, while not having to deal with sync primitives.
type Parser struct {
	Schema *schema.WebRPCSchema

	// Cache parsed types to improve performance and so we can traverse circular dependencies.
	ParsedTypes map[types.Type]*schema.VarType

	ParsedEnumTypes map[string]*schema.Type // Helps lookup enum types by pkg easily.

	InlineMode    bool // When traversing `json:",inline"`, we don't want to store the struct type as WebRPC message.
	ImportedPaths map[string]struct{}

	SchemaPkgName string // Schema file's package name.

	Pkg *packages.Package
}

func New(pkg *packages.Package) *Parser {
	return &Parser{
		Schema: &schema.WebRPCSchema{
			WebrpcVersion: "v1",
			SchemaName:    "",
			SchemaVersion: "",
		},
		SchemaPkgName:   pkg.PkgPath,
		ParsedTypes:     map[types.Type]*schema.VarType{},
		Pkg:             pkg,
		ParsedEnumTypes: map[string]*schema.Type{},

		// TODO: Change this to map[*types.Package]string so we can rename duplicated pkgs?
		ImportedPaths: map[string]struct{}{
			// Initial schema file's package name artificially set by golang.org/x/tools/go/packages.
			"command-line-arguments": {},
		},
	}
}
