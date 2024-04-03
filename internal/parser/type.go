package parser

import (
	"fmt"
	"go/types"

	"github.com/webrpc/webrpc/schema"
)

func (p *Parser) ParseType(typ types.Type) (*schema.VarType, error) {
	return p.ParseNamedType(nil, typ)
}

func (p *Parser) ParseBasic(typ *types.Basic) (*schema.VarType, error) {
	var varType schema.VarType
	if err := schema.ParseVarTypeExpr(p.Schema, typ.Name(), &varType); err != nil {
		return nil, fmt.Errorf("failed to parse basic type: %v: %w", typ.Name(), err)
	}

	return &varType, nil
}
