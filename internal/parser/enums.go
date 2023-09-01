package parser

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/webrpc/webrpc/schema"
)

// CollectEnums collects ENUM definitions, ie.:
//
//	// approved = 0
//	// pending  = 1
//	// closed   = 2
//	// new      = 3
//	type Status gospeak.Enum[int]
func (p *Parser) CollectEnums() error {
	debug := spew.NewDefaultConfig()
	debug.DisableMethods = true
	debug.DisablePointerAddresses = true
	debug.Indent = "\t"
	debug.SortKeys = true

	gospeakImportFound := false
	for _, file := range p.Pkg.Syntax {
		for _, decl := range file.Decls {
			if typeDeclaration, ok := decl.(*ast.GenDecl); ok && typeDeclaration.Tok == token.IMPORT {
				for _, spec := range typeDeclaration.Specs {
					if importSpec, ok := spec.(*ast.ImportSpec); ok {
						if strings.Contains(importSpec.Path.Value, `"github.com/golang-cz/gospeak/enum"`) {
							gospeakImportFound = true
						}
					}
				}
			}
		}
	}
	if !gospeakImportFound {
		return nil
	}

	for _, file := range p.Pkg.Syntax {
		for _, decl := range file.Decls {
			if typeDeclaration, ok := decl.(*ast.GenDecl); ok && typeDeclaration.Tok == token.TYPE {
				for _, spec := range typeDeclaration.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						selExpr, ok := typeSpec.Type.(*ast.SelectorExpr)
						if !ok {
							continue
						}
						ident, ok := selExpr.X.(*ast.Ident)
						if !ok {
							continue
						}

						// type Status enum.Int64
						enumName := typeSpec.Name.Name   // Status
						pkgName := ident.Name            // enum
						enumTypeName := selExpr.Sel.Name // Int64
						if pkgName != "enum" || enumName == "" || enumTypeName == "" {
							continue
						}

						enumElemType, ok := schema.CoreTypeFromString[strings.ToLower(enumTypeName)]
						if !ok {
							return fmt.Errorf("unknown enum type %v", enumTypeName)
						}

						enumType := &schema.Type{
							Kind: schema.TypeKind_Enum,
							Name: enumName,
							Type: &schema.VarType{
								Expr: enumElemType.String(),
								Type: enumElemType,
							},
							Fields: []*schema.TypeField{}, // webrpc TODO: should be Enums
						}

						doc := typeDeclaration.Doc
						if doc != nil {
							// name       value
							// ----------------
							// approved = 0
							// pending  = 1
							// closed   = 2
							// new      = 3
							for i, comment := range doc.List {
								commentValue, _ := strings.CutPrefix(comment.Text, "//")
								name, value, found := strings.Cut(commentValue, "=") // approved = 0
								if !found {                                          // approved
									name = commentValue
									value = fmt.Sprintf("%v", i)
								}
								enumType.Fields = append(enumType.Fields, &schema.TypeField{
									Name: strings.TrimSpace(name),
									TypeExtra: schema.TypeExtra{
										Value: strings.TrimSpace(value),
									},
								})
							}
						}

						p.Schema.Types = append(p.Schema.Types, enumType)
						p.ParsedEnumTypes[fmt.Sprintf("%v.%v", p.Pkg.PkgPath, enumName)] = enumType
					}
				}
			}
		}
	}

	return nil
}
