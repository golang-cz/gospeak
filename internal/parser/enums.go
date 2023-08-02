package parser

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/webrpc/webrpc/schema"
)

func (p *Parser) CollectEnums() error {
	enumValues := []*schema.TypeField{}
	for _, file := range p.Pkg.Syntax {
		for _, decl := range file.Decls {
			if typeDeclaration, ok := decl.(*ast.GenDecl); ok && typeDeclaration.Tok == token.TYPE {
				for _, spec := range typeDeclaration.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						typeName := typeSpec.Name.Name
						if typeName == "Enum" {
							fmt.Println(typeDeclaration, typeSpec.Name.Obj, typeName, typeSpec.Type)
							if indent, ok := typeSpec.Type.(*ast.Ident); ok {
								typeName := indent.Name
								//fmt.Println(typeName)
								if typeName == "Enum" {
									doc := typeDeclaration.Doc
									if doc != nil {
										for i, comment := range doc.List {
											commentValue, _ := strings.CutPrefix(comment.Text, "//")
											name, value, found := strings.Cut(commentValue, "=") // approved = 0
											if !found {                                          // approved
												value = fmt.Sprintf("%v", i)
												name = commentValue
											}
											enumValues = append(enumValues, &schema.TypeField{
												Name: strings.TrimSpace(name),
												TypeExtra: schema.TypeExtra{
													Value: strings.TrimSpace(value),
												},
											})
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	//enumType := &Schema.Type{
	//	Kind:   Schema.TypeKind_Enum,
	//	Name:   name,
	//	Type:   enumElemType,
	//	Fields: enumValues, // webrpc TODO: should be Enums
	//}
	//
	//p.Schema.Types = append(p.Schema.Types, enumType)

	return nil
}
