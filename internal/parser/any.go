package parser

import (
	"go/types"

	"github.com/webrpc/webrpc/schema"
)

func (p *Parser) ParseAny(parent *types.Named, iface *types.Interface) (*schema.VarType, error) {
	varType := &schema.VarType{
		Expr: "any",
		Type: schema.T_Any,
	}

	return varType, nil
}
