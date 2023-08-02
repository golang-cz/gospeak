package parser

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/webrpc/webrpc/schema"
)

func (p *Parser) ParseStruct(typeName string, structTyp *types.Struct) (*schema.VarType, error) {
	structType := &schema.Type{
		Kind: "struct",
		Name: typeName,
	}

	for i := 0; i < structTyp.NumFields(); i++ {
		structField := structTyp.Field(i)
		if !structField.Exported() {
			continue
		}
		structTags := structTyp.Tag(i)

		if structField.Embedded() || strings.Contains(structTags, `json:",inline"`) {
			varType, err := p.ParseNamedType("", structField.Type())
			if err != nil {
				return nil, fmt.Errorf("parsing var %v: %w", structField.Name(), err)
			}

			if varType.Type == schema.T_Struct {
				for _, embeddedField := range varType.Struct.Type.Fields {
					structType.Fields = appendOrOverrideExistingField(structType.Fields, embeddedField)
				}
			}
			continue
		}

		field, err := p.parseStructField(typeName+"Field", structField, structTags)
		if err != nil {
			return nil, fmt.Errorf("parsing struct field %v: %w", i, err)
		}
		if field != nil {
			structType.Fields = appendOrOverrideExistingField(structType.Fields, field)
		}
	}

	p.Schema.Types = append(p.Schema.Types, structType)

	return &schema.VarType{
		Expr: typeName,
		Type: schema.T_Struct,
		Struct: &schema.VarStructType{
			Name: typeName,
			Type: structType,
		},
	}, nil
}

// parses single Go struct field
// if the field is embedded, ie. `json:",inline"`, parse recursively
func (p *Parser) parseStructField(structTypeName string, field *types.Var, structTags string) (*schema.TypeField, error) {
	fieldName := field.Name()
	fieldType := field.Type()

	jsonFieldName := fieldName
	goFieldType := p.GoTypeName(fieldType)
	optional := false

	goFieldImport := p.GoTypeImport(fieldType)

	jsonTag, ok := GetJsonTag(structTags)
	if ok {
		if jsonTag.Name == "-" { // struct field ignored by `json:"-"` struct tag
			return nil, nil
		}

		if jsonTag.Name != "" {
			jsonFieldName = jsonTag.Name
		}

		optional = jsonTag.Omitempty
		if optional {
			goFieldType = "*" + goFieldType
		}

		if jsonTag.IsString { // struct field forced to be string by `json:",string"`

			structField := &schema.TypeField{
				Name: jsonFieldName,
				Type: &schema.VarType{
					Expr: "string",
					Type: schema.T_String,
				},
				TypeExtra: schema.TypeExtra{
					Meta: []schema.TypeFieldMeta{
						{"go.field.name": fieldName},
						{"go.field.type": goFieldType},
					},
					Optional: optional,
				},
			}
			if goFieldImport != "" {
				structField.TypeExtra.Meta = append(structField.TypeExtra.Meta,
					schema.TypeFieldMeta{"go.type.import": goFieldImport},
				)
			}
			structField.TypeExtra.Meta = append(structField.TypeExtra.Meta,
				schema.TypeFieldMeta{"go.tag.json": jsonTag.Value},
			)

			return structField, nil
		}
	}

	if _, ok := field.Type().Underlying().(*types.Pointer); ok {
		optional = true
		goFieldType = "*" + goFieldType
	}

	if _, ok := field.Type().Underlying().(*types.Struct); ok {
		// Anonymous struct fields.
		// Example:
		//   type Something struct {
		// 	   AnonymousField struct { // no explicit struct type name
		//       Name string
		//     }
		//   }
		structTypeName = /*structTypeName + */ "Anonymous" + field.Name()
	}

	varType, err := p.ParseNamedType(goFieldType, fieldType)
	if err != nil {
		return nil, fmt.Errorf("failed to parse var %v: %w", field.Name(), err)
	}

	structField := &schema.TypeField{
		Name: jsonFieldName,
		Type: varType,
		TypeExtra: schema.TypeExtra{
			Meta: []schema.TypeFieldMeta{
				{"go.field.name": fieldName},
				{"go.field.type": goFieldType},
			},
			Optional: optional,
		},
	}
	if goFieldImport != "" {
		structField.TypeExtra.Meta = append(structField.TypeExtra.Meta,
			schema.TypeFieldMeta{"go.type.import": goFieldImport},
		)
	}
	if jsonTag.Value != "" {
		structField.TypeExtra.Meta = append(structField.TypeExtra.Meta, schema.TypeFieldMeta{"go.tag.json": jsonTag.Value})
	}

	return structField, nil
}

// Appends message field to the given slice, while also removing any previously defined field of the same name.
// This lets us overwrite embedded fields, exactly how Go does it behind the scenes in the JSON marshaller.
func appendOrOverrideExistingField(slice []*schema.TypeField, newItem *schema.TypeField) []*schema.TypeField {
	// Let's try to find an existing item of the same name and delete it.
	for i, item := range slice {
		if item.Name == newItem.Name {
			// Delete item.
			copy(slice[i:], slice[i+1:])
			slice = slice[:len(slice)-1]
		}
	}
	// And then append the new item at the end of the slice.
	return append(slice, newItem)
}
