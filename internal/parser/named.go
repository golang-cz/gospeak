package parser

import (
	"fmt"
	"go/types"

	"github.com/webrpc/webrpc/schema"
)

func (p *Parser) ParseNamedType(typeName string, typ types.Type) (varType *schema.VarType, err error) {
	// Return a parsedType from cache, if exists.
	if parsedType, ok := p.ParsedTypes[typ]; ok {
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
	p.ParsedTypes[typ] = cacheDoNotReturn

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
		typeName := p.GoTypeName(typ)

		if pkg != nil {
			if typeName == "time.Time" {
				return &schema.VarType{
					Expr: "timestamp",
					Type: schema.T_Timestamp,
				}, nil
			}
		}

		// If the type implements encoding.TextMarshaler, it's a string.
		if isTextMarshaler(v, pkg) {
			return &schema.VarType{
				Expr: "string",
				Type: schema.T_String,
			}, nil
		}

		switch u := underlying.(type) {

		case *types.Pointer:
			// Named pointer. Webrpc can't handle that.
			// Example:
			//   type NamedPtr *Obj

			// Go for the underlying element type name (ie. `Obj`).
			return p.ParseNamedType(p.GoTypeName(underlying), u.Underlying())

		case *types.Slice, *types.Array:
			// Named slice/array. Webrpc can't handle that.
			// Example:
			//  type NamedSlice []int
			//  type NamedSlice []Obj

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

			var elem types.Type
			//                  = u.Elem().Underlying()
			// NOTE: Calling the above assignment fails to build with this error:
			//       "u.Elem undefined (type types.Type has no field or method Elem)"
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
				basicType, err := p.ParseBasic(basic)
				if err != nil {
					return nil, fmt.Errorf("failed to parse []namedBasicType: %w", err)
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
			return p.ParseNamedType(p.GoTypeName(underlying), u.Underlying())

		default:
			if isJsonMarshaller(v, pkg) {
				return &schema.VarType{
					Expr: "any",
					Type: schema.T_Any,
				}, nil
			}

			return p.ParseNamedType(typeName, underlying)
		}

	case *types.Basic:
		return p.ParseBasic(v)

	case *types.Struct:
		return p.ParseStruct(typeName, v)

	case *types.Slice:
		return p.ParseSlice(typeName, v)

	case *types.Interface:
		return p.ParseAny(typeName, v)

	case *types.Map:
		return p.ParseMap(typeName, v)

	case *types.Pointer:
		if typeName == "" {
			return p.ParseNamedType(p.GoTypeName(v), v.Elem())
		}
		return p.ParseNamedType(typeName, v.Elem())

	default:
		return nil, fmt.Errorf("unsupported argument type %T", typ)
	}
}
