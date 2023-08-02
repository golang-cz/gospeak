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
	//panic(debug.Sdump(p.Pkg.Syntax))

	gospeakImportFound := false
	for _, file := range p.Pkg.Syntax {
		for _, decl := range file.Decls {
			if typeDeclaration, ok := decl.(*ast.GenDecl); ok && typeDeclaration.Tok == token.IMPORT {
				for _, spec := range typeDeclaration.Specs {
					if importSpec, ok := spec.(*ast.ImportSpec); ok {
						if strings.Contains(importSpec.Path.Value, `"github.com/golang-cz/gospeak"`) {
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
						var enumName, enumTypeName string
						if iExpr, ok := typeSpec.Type.(*ast.IndexExpr); ok {
							if selExpr, ok := iExpr.X.(*ast.SelectorExpr); ok {
								pkgName, ok := selExpr.X.(*ast.Ident)
								if ok && pkgName.Name == "gospeak" && selExpr.Sel.Name == "Enum" {
									if id, ok := iExpr.Index.(*ast.Ident); ok {
										enumName = typeSpec.Name.Name
										enumTypeName = id.Name
									}
								}
							}
						}
						if enumName == "" || enumTypeName == "" {
							continue
						}

						enumElemType, ok := schema.CoreTypeFromString[enumTypeName]
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
								before, after, found := strings.Cut(commentValue, "=") // approved = 0
								if !found {                                            // approved
									before = commentValue
									after = fmt.Sprintf("%v", i)
								}
								// This looks reversed. TODO: webrpc enum type
								enumType.Fields = append(enumType.Fields, &schema.TypeField{
									Name: strings.TrimSpace(after),
									TypeExtra: schema.TypeExtra{
										Value: strings.TrimSpace(before),
									},
								})
							}
						}

						p.Schema.Types = append(p.Schema.Types, enumType)
					}
				}
			}
		}
	}

	return nil
}
