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

	// Cache for already parsed types, to improve performance & so we can traverse circular dependencies.
	ParsedTypes     map[types.Type]*schema.VarType
	ParsedTypeNames map[string]struct{}

	InlineMode    bool // When traversing `json:",inline"`, we don't want to store the struct type as WebRPC message.
	ImportedPaths map[string]struct{}

	SchemaPkgName string // Schema file's package name.

	Pkg *packages.Package
}
