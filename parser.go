package gospeak

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-cz/gospeak/internal/parser"
	"github.com/webrpc/webrpc/schema"
	"golang.org/x/tools/go/packages"
)

type Target struct {
	Schema        *schema.WebRPCSchema
	Generator     string
	InterfaceName string
	OutFile       string
	Opts          map[string]interface{}
}

// Parse Go source file or package folder and return WebRPC schema.
func Parse(filePath string) ([]*Target, error) {
	path, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get directory from %q: %w", path, err)
	}

	file, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open %q", path)
	}
	if file.Mode().IsRegular() {
		// Parse all files in the given schema file's directory, so the parser can see all the pkg files.
		path = filepath.Dir(path)
	}

	cfg := &packages.Config{
		Dir:  path,
		Mode: packages.NeedName | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedImports,
	}

	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		return nil, fmt.Errorf("failed to load Go packages from %q: %w", path, err)
	}

	// Print all errors.
	for _, pkg := range pkgs {
		for _, pkgErr := range pkg.Errors {
			fmt.Fprintln(os.Stderr, pkgErr)
		}

		for _, typeErr := range pkg.TypeErrors {
			fmt.Fprintln(os.Stderr, typeErr)
		}
	}

	if len(pkgs) != 1 {
		return nil, fmt.Errorf("failed to load Go package (len=%v) from %q", len(pkgs), path)
	}
	pkg := pkgs[0]

	if len(pkg.Errors) > 0 || len(pkg.TypeErrors) > 0 {
		return nil, fmt.Errorf("%v errors", len(pkg.Errors)+len(pkg.TypeErrors))
	}

	// Collect Go interfaces with `//go:webrpc` comments.
	targets, err := CollectInterfaces(pkg)
	if err != nil {
		return nil, fmt.Errorf("collecting Go interfaces: %w", err)
	}

	cache := map[string]*schema.WebRPCSchema{}
	for _, target := range targets {
		if interfaceSchema, ok := cache[target.InterfaceName]; ok {
			// Hit.
			target.Schema = interfaceSchema
		}

		// Miss.
		p := &parser.Parser{
			Schema: &schema.WebRPCSchema{
				WebrpcVersion: "v1",
				SchemaName:    target.InterfaceName,
				SchemaVersion: "vTODO",
			},
			SchemaPkgName:   pkg.PkgPath,
			ParsedTypes:     map[types.Type]*schema.VarType{},
			ParsedTypeNames: map[string]struct{}{},
			Pkg:             pkg,

			// TODO: Change this to map[*types.Package]string so we can rename duplicated pkgs?
			ImportedPaths: map[string]struct{}{
				// Initial schema file's package name artificially set by golang.org/x/tools/go/packages.
				"command-line-arguments": {},
			},
		}

		if err := p.CollectEnums(); err != nil {
			return nil, fmt.Errorf("collecting enums: %w", err)
		}

		obj := pkg.Types.Scope().Lookup(target.InterfaceName)
		if obj == nil {
			return nil, fmt.Errorf("type interface %v{} not found", target.InterfaceName)
		}

		iface, ok := obj.Type().Underlying().(*types.Interface)
		if !ok {
			return nil, fmt.Errorf("type %v{} is %T", target.InterfaceName, obj.Type().Underlying())
		}

		if err := p.ParseInterfaceMethods(iface, target.InterfaceName); err != nil {
			return nil, fmt.Errorf("failed to parse interface %q: %w", target.InterfaceName, err)
		}

		target.Schema = p.Schema
		cache[target.InterfaceName] = p.Schema
	}

	return targets, nil
}

// Find all Go interfaces with the special //go:webrpc comments.
func CollectInterfaces(pkg *packages.Package) ([]*Target, error) {
	var targets []*Target

	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			if typeDeclaration, ok := decl.(*ast.GenDecl); ok && typeDeclaration.Tok == token.TYPE {
				for _, spec := range typeDeclaration.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if _, ok := typeSpec.Type.(*ast.InterfaceType); ok {
							doc := typeDeclaration.Doc
							if doc != nil {
								for _, comment := range doc.List {
									if webrpcCmd, hasPrefix := strings.CutPrefix(comment.Text, "//go:webrpc "); hasPrefix {
										target, err := parseWebrpcCommand(webrpcCmd)
										if err != nil {
											return nil, fmt.Errorf("failed to parse %s", comment.Text)
										}
										target.InterfaceName = typeSpec.Name.Name
										targets = append(targets, target)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return targets, nil
}

// Parses webrpc CLI command into a target, ie. webrpc typescript@v0.11.0 -client -out=./videoAuthoringClient.gen.ts.
func parseWebrpcCommand(cmd string) (*Target, error) {
	target := &Target{
		Opts: map[string]interface{}{},
	}

	for _, arg := range strings.Split(cmd, " ") {
		name, value, _ := strings.Cut(arg, "=")

		if strings.HasPrefix(name, "-") {
			name = strings.TrimLeft(name, "-")

			// target options
			if name == "out" {
				target.OutFile = value
			} else {
				target.Opts[name] = value
			}
		} else {
			if target.Generator != "" {
				return nil, fmt.Errorf("unexpected argument %v", name)
			}
			target.Generator = name
		}
	}

	if target.OutFile == "" {
		return nil, fmt.Errorf("-out=<path> flag is required")
	}

	return target, nil
}
