package parser

import (
	"fmt"
	"go/types"

	"github.com/webrpc/webrpc/schema"
)

func (p *Parser) ParseStruct(parent *types.Named, structTyp *types.Struct) (*schema.VarType, error) {
	webrpcTypeName := p.WebrpcTypeName(parent)

	structType := &schema.Type{
		Kind: "struct",
		Name: webrpcTypeName,
		// NOTE: Webrpc (templates) don't allow annotating structs right now -- perhaps we can change allow it?
		//
		// TypeExtra: schema.TypeExtra{
		// 	Meta: []schema.TypeFieldMeta{
		// 		{"go.field.name": goTypeName},
		// 		{"go.type": goTypeName},
		// 	},
		// 	//Optional: optional,
		// },
	}

	for i := 0; i < structTyp.NumFields(); i++ {
		structField := structTyp.Field(i)
		if !structField.Exported() {
			continue
		}
		structTags := structTyp.Tag(i)

		jsonTag, _ := GetJsonTag(structTags)
		if jsonTag.Name == "-" { // struct field ignored by `json:"-"` struct tag
			continue
		}

		if structField.Embedded() || jsonTag.Inline {
			varType, err := p.ParseNamedType(parent, structField.Type())
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

		field, err := p.parseStructField(parent, structField, jsonTag)
		if err != nil {
			return nil, fmt.Errorf("parsing struct field %v: %w", i, err)
		}

		if field != nil {
			structType.Fields = appendOrOverrideExistingField(structType.Fields, field)
		}
	}

	p.Schema.Types = append(p.Schema.Types, structType)

	return &schema.VarType{
		Expr: webrpcTypeName,
		Type: schema.T_Struct,
		Struct: &schema.VarStructType{
			Name: webrpcTypeName,
			Type: structType,
		},
	}, nil
}

// parses single Go struct field
// if the field is embedded, ie. `json:",inline"`, parse recursively
func (p *Parser) parseStructField(parent *types.Named, field *types.Var, jsonTag JsonTag) (*schema.TypeField, error) {
	fieldName := field.Name()
	fieldType := field.Type()

	jsonFieldName := fieldName
	goFieldType := p.GoTypeName(fieldType)
	optional := false

	goFieldImport := p.GoTypeImport(fieldType)

	if jsonTag.Name != "" {
		if jsonTag.Name == "-" { // struct field ignored by `json:"-"` struct tag
			return nil, nil
		}
		jsonFieldName = jsonTag.Name
	}

	if jsonTag.Omitempty {
		optional = jsonTag.Omitempty
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

	if _, ok := fieldType.Underlying().(*types.Pointer); ok {
		optional = true
		goFieldType = "*" + goFieldType
	}

	if _, ok := fieldType.Underlying().(*types.Struct); ok {
		// Anonymous struct fields.
		// Example:
		//   type Something struct {
		// 	   AnonymousField struct { // no explicit struct type name
		//       Name string
		//     }
		//   }
		/* "Anonymous" + field.Name() */
	}

	// TODO: Can we ever see type aliases here? If so, how do you trigger this?
	if named, ok := fieldType.(*types.Named); ok {
		if named.Obj().IsAlias() {
			panic(fmt.Sprintf("alias: %v", fieldType))
		}
	}

	// TODO: Can we ever see type aliases here? If so, how do you trigger this?
	if named, ok := fieldType.Underlying().(*types.Named); ok {
		if named.Obj().IsAlias() {
			panic(fmt.Sprintf("alias: %v", fieldType))
		}
	}

	varType, err := p.ParseNamedType(parent, fieldType)
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
