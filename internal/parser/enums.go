package parser

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/webrpc/webrpc/schema"
	"golang.org/x/tools/go/packages"
)

func (p *Parser) ExtractEnumConsts(pkg *packages.Package) error {
	enumMap := map[string]*schema.Type{}

	// First pass: find all enum types with //gospeak:enum comment.
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec := spec.(*ast.TypeSpec)

				comments := []string{}

				// Look for last comment in form of //gospeak:enum
				if genDecl.Doc == nil || len(genDecl.Doc.List) == 0 {
					continue
				}
				if genDecl.Doc.List[len(genDecl.Doc.List)-1].Text != "//gospeak:enum" {
					continue
				}
				for _, comment := range genDecl.Doc.List[:(len(genDecl.Doc.List) - 1)] {
					comments = append(comments, strings.TrimSpace(strings.TrimPrefix(comment.Text, "//")))
				}

				enumName := typeSpec.Name.Name
				enumElemType := pkg.TypesInfo.TypeOf(typeSpec.Type)
				if enumElemType == nil {
					continue
				}

				enumType := &schema.Type{
					Kind: schema.TypeKind_Enum,
					Name: enumName,
					Type: &schema.VarType{
						Expr: enumElemType.String(),
						Type: schema.T_Enum,
					},
					Fields:   []*schema.TypeField{},
					Comments: comments,
				}

				enumImportTypeName := fmt.Sprintf("%v.%v", p.Pkg.PkgPath, enumName)

				// Save for second pass
				enumMap[enumImportTypeName] = enumType

				// Save to schema
				p.Schema.Types = append(p.Schema.Types, enumType)
				p.ParsedEnumTypes[enumImportTypeName] = enumType
			}
		}
	}

	// Second pass: collect consts
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.CONST {
				continue
			}

			for _, spec := range genDecl.Specs {
				valSpec := spec.(*ast.ValueSpec)

				for _, ident := range valSpec.Names {
					obj := pkg.TypesInfo.Defs[ident]
					if obj == nil {
						continue
					}
					constObj, ok := obj.(*types.Const)
					if !ok {
						continue
					}

					enumType, ok := enumMap[constObj.Type().String()]
					if !ok {
						continue
					}

					enumName := constObj.Type().String()
					fmt.Println(enumName)

					// Get value from trailing comment, e.g. // "some value"
					var value string
					if valSpec.Comment != nil {
						value = strings.TrimSpace(valSpec.Comment.Text())
						if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
							// Parse quoted string, handling escaped quotes
							value = value[1 : len(value)-1]              // Remove outer quotes
							value = strings.ReplaceAll(value, `\"`, `"`) // Unescape quotes
						}
					}
					if value == "" {
						return fmt.Errorf(`Enum %v: Missing value comment, e.g. // "value"`, enumName)
					}

					// TODO: how can we pass a custom key value (uint8) to the webrpc enum?
					enumType.Fields = append(enumType.Fields, &schema.TypeField{
						Name: ident.Name,
						TypeExtra: schema.TypeExtra{
							Value: fmt.Sprintf("%q", value), // hmm, webrpc requires quotes for string enums values
						},
					})

					p.Schema.Types = append(p.Schema.Types, enumType)
					p.ParsedEnumTypes[fmt.Sprintf("%v.%v", p.Pkg.PkgPath, enumName)] = enumType
				}
			}
		}
	}

	return nil
}

// CollectEnums collects ENUM definitions, ie.:
//
//	// approved = 0
//	// pending  = 1
//	// closed   = 2
//	// new      = 3
//	type Status gospeak.Enum[int]
//
// Deprecated: We have switche to ExtractEnumConsts instead. Left here for now to print error to users.
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

						typeName := fmt.Sprintf("%v.%v", p.Pkg.PkgPath, enumName)

						return fmt.Errorf(`Obsolete ENUM definition for type %v.

						Please, migrate to this new ENUM format:
	
						//gospeak:enum
						type Status uint8
						
						const (
							StatusUnknown Status = iota // "unknown"
							StatusActive                // "active"
						)
					`, typeName)

						// p.Schema.Types = append(p.Schema.Types, enumType)
						// p.ParsedEnumTypes[typeName] = enumType
					}
				}
			}
		}
	}

	return nil
}
