package gospeak

import (
	"fmt"
	"go/types"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/golang-cz/textcase"
	"github.com/pkg/errors"
	"github.com/webrpc/webrpc/schema"
)

// Given `db:"id,omitempty,pk" json:"id,string"` struct tag,
// this regex will return the following three submatches:
// [0]: json:"id,string"
// [1]: id
// [2]: ,string
var jsonTagRegex, _ = regexp.Compile(`\s?json:\"([^,\"]*)(,[^\"]*)?\"`)

func (p *parser) parseType(typ types.Type) (*schema.VarType, error) {
	return p.parseNamedType("", typ)
}

func (p *parser) parseNamedType(typeName string, typ types.Type) (varType *schema.VarType, err error) {
	// Return a parsedType from cache, if exists.
	if parsedType, ok := p.parsedTypes[typ]; ok {
		return parsedType, nil
	}

	// Otherwise, create a new parsedType record and warm up the cache up-front.
	// Claim the cache key and fill in the value later in defer(). Meanwhile, any
	// following recursive call(s) to this function (ie. on recursive types like
	// self-referencing structs, linked lists, graphs, circular dependencies etc.)
	// will return early with the same pointer.
	//
	// Note: We're parsing sequentially, no need for sync.Map.
	cacheDoNotReturn := &schema.VarType{
		Expr: typeName,
	}
	p.parsedTypes[typ] = cacheDoNotReturn

	defer func() {
		if varType != nil {
			*cacheDoNotReturn = *varType // Update the cache value via pointer dereference.
			varType = cacheDoNotReturn
		}
	}()

	switch v := typ.(type) {
	case *types.Named:
		pkg := v.Obj().Pkg()
		underlying := v.Underlying()
		typeName := p.goTypeName(typ)

		switch u := underlying.(type) {

		case *types.Pointer:
			// Named pointer. Webrpc can't handle that.
			// Example:
			//   type NamedPtr *Obj

			// Go for the underlying element type name (ie. `Obj`).
			return p.parseNamedType(p.goTypeName(underlying), u.Underlying())

		case *types.Slice, *types.Array:
			// Named slice/array. Webrpc can't handle that.
			// Example:
			//  type NamedSlice []int
			//  type NamedSlice []Obj

			// If the named type is a slice/array and implements encoding.TextMarshaler,
			// we assume it's []string.
			if isTextMarshaler(v, pkg) {
				return &schema.VarType{
					Expr: "[]string",
					Type: schema.T_List,
					List: &schema.VarListType{
						Elem: &schema.VarType{
							Expr: "string",
							Type: schema.T_String,
						},
					},
				}, nil
			}

			// If the named type is a slice/array and implements json.Marshaler,
			// we assume it's []any.
			if isJsonMarshaller(v, pkg) {
				return &schema.VarType{
					Expr: "[]any",
					Type: schema.T_List,
					List: &schema.VarListType{
						Elem: &schema.VarType{
							Expr: "any",
							Type: schema.T_Any,
						},
					},
				}, nil
			}

			var elem types.Type // = u.Elem().Underlying()
			// NOTE: Calling the above u.Elem().Underlying() directly fails to build with
			//       "u.Elem undefined (type types.Type has no field or method Elem)" error
			//       even though both *types.Slice and *types.Array have the .Elem() method.
			switch underlyingElem := u.(type) {
			case *types.Slice:
				elem = underlyingElem.Elem().Underlying()
			case *types.Array:
				elem = underlyingElem.Elem().Underlying()
			}

			// If the named type is a slice/array and its underlying element
			// type is basic type (ie. `int`), we go for it directly.
			if basic, ok := elem.(*types.Basic); ok {
				basicType, err := p.parseBasic(basic)
				if err != nil {
					return nil, errors.Wrap(err, "failed to parse []namedBasicType")
				}
				return &schema.VarType{
					Expr: fmt.Sprintf("[]%v", basicType.String()),
					Type: schema.T_List,
					List: &schema.VarListType{
						Elem: basicType,
					},
				}, nil
			}

			// Otherwise, go for the underlying element type name (ie. `Obj`).
			return p.parseNamedType(p.goTypeName(underlying), u.Underlying())

		default:
			if isTextMarshaler(v, pkg) {
				return &schema.VarType{
					Expr: "string",
					Type: schema.T_String,
				}, nil
			}

			if isJsonMarshaller(v, pkg) {
				return &schema.VarType{
					Expr: "any",
					Type: schema.T_Any,
				}, nil
			}

			if pkg == nil {
				return p.parseNamedType(typeName, underlying)
			}

			if pkg.Path() == "time" && v.Obj().Id() == "Time" {
				return &schema.VarType{
					Expr: "timestamp",
					Type: schema.T_Timestamp,
				}, nil
			}

			return p.parseNamedType(typeName, underlying)
		}

	case *types.Basic:
		return p.parseBasic(v)

	case *types.Struct:
		return p.parseStruct(typeName, v)

	case *types.Slice:
		return p.parseSlice(typeName, v)

	case *types.Interface:
		return p.parseInterface(typeName, v)

	case *types.Map:
		return p.parseMap(typeName, v)

	case *types.Pointer:
		if typeName == "" {
			return p.parseNamedType(p.goTypeName(v), v.Elem())
		}
		return p.parseNamedType(typeName, v.Elem())

	default:
		return nil, errors.Errorf("unsupported argument type %T", typ)
	}
}

func findFirstAlphaCharPosition(s string) int {
	for i, char := range s {
		if unicode.IsLetter(char) {
			return i
		}
	}
	return 0
}

func (p *parser) goTypeName(typ types.Type) string {
	name := typ.String() // []github.com/golang-cz/gospeak/pkg.Typ

	// switch name {
	// case "context.Context":
	// case "*github.com/golang-cz/gospeak/_examples/petStore/proto.Pet":
	// case "github.com/golang-cz/gospeak/_examples/petStore/proto.Pet":
	// case "int64":
	// case "string":
	// case "bool":

	// default:
	// 	panic(name)
	// }

	pos := findFirstAlphaCharPosition(name)
	prefix := name[:pos]

	name = filepath.Base(name)                                 // pkg.Typ
	name = strings.TrimPrefix(name, p.schemaPkgName+".")       // Typ (ignore root pkg).
	name = strings.TrimPrefix(name, "command-line-arguments.") // Ignore "command-line-arguments" pkg autogenerated by Go tool chain.

	return prefix + name // []Typ
}

func (p *parser) ridlTypeName(typ types.Type) string {
	goTypeName := p.goTypeName(typ)

	if unicode.IsUpper(rune(goTypeName[0])) {
		return textcase.PascalCase(strings.ReplaceAll(goTypeName, ".", ""))
	}

	return textcase.CamelCase(strings.ReplaceAll(goTypeName, ".", ""))
}

func (p *parser) parseBasic(typ *types.Basic) (*schema.VarType, error) {
	var varType schema.VarType
	if err := schema.ParseVarTypeExpr(p.schema, typ.Name(), &varType); err != nil {
		return nil, errors.Wrapf(err, "failed to parse basic type: %v", typ.Name())
	}

	return &varType, nil
}

func (p *parser) parseStruct(typeName string, structTyp *types.Struct) (varType *schema.VarType, err error) {
	msg := &schema.Type{
		Kind: "struct",
		Name: typeName,
	}

	for i := 0; i < structTyp.NumFields(); i++ {
		field := structTyp.Field(i)
		if !field.Exported() {
			continue
		}

		tag := structTyp.Tag(i)
		if field.Embedded() || strings.Contains(tag, `json:",inline"`) {
			varType, err := p.parseNamedType("", field.Type())
			if err != nil {
				return nil, errors.Wrapf(err, "failed to parse var %v", field.Name())
			}

			if varType.Type == schema.T_Struct {
				for _, embeddedField := range varType.Struct.Type.Fields {
					msg.Fields = appendTypeFieldAndDeleteExisting(msg.Fields, embeddedField)
				}
			}
			continue
		}

		optional := false

		fieldName := field.Name()
		jsonFieldName := fieldName
		goFieldName := fieldName //strings.Title(fieldName)

		fieldType := field.Type()
		ridlFieldType := fieldType.String() //p.ridlTypeName(fieldType)
		goFieldType := p.goTypeName(fieldType)

		if strings.Contains(tag, `json:"`) {
			submatches := jsonTagRegex.FindStringSubmatch(tag)
			// Submatches from the jsonTagRegex:
			// [0]: json:"deleted_by,omitempty,string"
			// [1]: deleted_by
			// [2]: ,omitempty,string
			if len(submatches) != 3 {
				return nil, errors.Errorf("unexpected number of json struct tag submatches")
			}
			if submatches[1] == "-" { // suppressed field in JSON struct tag
				continue
			}
			if submatches[1] != "" { // field name defined in JSON struct tag
				jsonFieldName = submatches[1]
			}
			optional = strings.Contains(submatches[2], ",omitempty")
			if strings.Contains(submatches[2], ",string") { // field type should be string in JSON
				msg.Fields = appendTypeFieldAndDeleteExisting(msg.Fields, &schema.TypeField{
					Name: jsonFieldName,
					Type: &schema.VarType{
						Expr: "string",
						Type: schema.T_String,
					},
					TypeExtra: schema.TypeExtra{Optional: optional},
				})
				continue
			}
		}

		if _, ok := field.Type().Underlying().(*types.Pointer); ok {
			optional = true
		}

		structTypeName := ""
		if _, ok := field.Type().Underlying().(*types.Struct); ok {
			// Anonymous struct fields.
			// Example:
			//   type Something struct {
			// 	   AnonymousField struct { // no explicit struct type name
			//       Name string
			//     }
			//   }
			structTypeName = typeName + field.Name()
		}

		varType, err := p.parseNamedType(structTypeName, fieldType)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse var %v", field.Name())
		}
		varType.Expr = ridlFieldType

		structField := &schema.TypeField{
			Name: jsonFieldName,
			Type: varType,
			TypeExtra: schema.TypeExtra{
				Meta: []schema.TypeFieldMeta{
					{"go.field.name": goFieldName},
					{"go.field.type": goFieldType},
				},
				Optional: optional,
			},
		}

		msg.Fields = appendTypeFieldAndDeleteExisting(msg.Fields, structField)
	}

	p.schema.Types = append(p.schema.Types, msg)

	return &schema.VarType{
		Expr: typeName,
		Type: schema.T_Struct,
		Struct: &schema.VarStructType{
			Name: typeName,
			Type: msg,
		},
	}, nil
}

func (p *parser) parseSlice(typeName string, sliceTyp *types.Slice) (*schema.VarType, error) {
	elem, err := p.parseNamedType(typeName, sliceTyp.Elem())
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse slice type")
	}

	actualTypeName := typeName
	if actualTypeName == "" {
		actualTypeName = elem.String()
	}

	varType := &schema.VarType{
		Expr: fmt.Sprintf("[]%v", actualTypeName),
		Type: schema.T_List,
		List: &schema.VarListType{
			Elem: elem,
		},
	}

	return varType, nil
}

func (p *parser) parseInterface(typeName string, iface *types.Interface) (*schema.VarType, error) {
	varType := &schema.VarType{
		Expr: "any",
		Type: schema.T_Any,
	}

	return varType, nil
}

func (p *parser) parseMap(typeName string, m *types.Map) (*schema.VarType, error) {
	key, err := p.parseNamedType(typeName, m.Key())
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse map key type")
	}

	value, err := p.parseNamedType(typeName, m.Elem())
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse map value type")
	}

	varType := &schema.VarType{
		Expr: fmt.Sprintf("map<%v,%v>", key, value),
		Type: schema.T_Map,
		Map: &schema.VarMapType{
			Key:   key.Type,
			Value: value,
		},
	}

	return varType, nil
}

// Returns true if the given type implements encoding.TextMarshaler
// and encoding.TextUnmarshaler interfaces.
func isTextMarshaler(typ types.Type, pkg *types.Package) bool {
	marshalTextMethod, _, _ := types.LookupFieldOrMethod(typ, true, pkg, "MarshalText")
	unmarshalTextMethod, _, _ := types.LookupFieldOrMethod(typ, true, pkg, "UnmarshalText")
	if marshalTextMethod != nil &&
		unmarshalTextMethod != nil &&
		strings.HasSuffix(marshalTextMethod.String(), ".MarshalText() ([]byte, error)") &&
		strings.HasSuffix(unmarshalTextMethod.String(), ".UnmarshalText(text []byte) error") {
		return true
	}
	return false
}

// Returns true if the given type implements json.Marshaler and
// json.Unmarshaler interfaces.
func isJsonMarshaller(typ types.Type, pkg *types.Package) bool {
	marshalJSONMethod, _, _ := types.LookupFieldOrMethod(typ, true, pkg, "MarshalJSON")
	unmarshalJSONMethod, _, _ := types.LookupFieldOrMethod(typ, true, pkg, "UnmarshalJSON")
	if marshalJSONMethod != nil &&
		unmarshalJSONMethod != nil &&
		strings.HasSuffix(marshalJSONMethod.String(), ".MarshalJSON() ([]byte, error)") &&
		strings.HasSuffix(unmarshalJSONMethod.String(), ".UnmarshalJSON(text []byte) error") {
		return true
	}
	return false
}

// Appends message field to the given slice, while also removing any previously defined field of the same name.
// This lets us overwrite embedded fields, exactly how Go does it behind the scenes in the JSON marshaller.
func appendTypeFieldAndDeleteExisting(slice []*schema.TypeField, newItem *schema.TypeField) []*schema.TypeField {
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
