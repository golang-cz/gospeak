package parser

import (
	"go/types"

	"github.com/webrpc/webrpc/schema"
	"golang.org/x/tools/go/packages"
)

// Parser walks the Go AST tree of a given package and returns WebRPC Schema.
//
// Walks the AST tree sequentially, without concurrency, to handle circular and
// recursive types. Aggressively caches parsed types to improve performance.
type Parser struct {
	Schema *schema.WebRPCSchema

	// ParsedTypes is a cache to improve performance and so we can traverse circular dependencies.
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
