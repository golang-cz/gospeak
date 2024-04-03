package parser

import (
	"fmt"
	"go/types"

	"github.com/webrpc/webrpc/schema"
)

func (p *Parser) ParseSlice(parent *types.Named, sliceTyp *types.Slice) (*schema.VarType, error) {
	elem, err := p.ParseNamedType(parent, sliceTyp.Elem())
	if err != nil {
		return nil, fmt.Errorf("failed to parse slice type: %w", err)
	}

	varType := &schema.VarType{
		Expr: fmt.Sprintf("[]%v", elem.String()),
		Type: schema.T_List,
		List: &schema.VarListType{
			Elem: elem,
		},
	}

	return varType, nil
}
