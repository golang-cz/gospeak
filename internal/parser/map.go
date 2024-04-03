package parser

import (
	"fmt"
	"go/types"

	"github.com/webrpc/webrpc/schema"
)

func (p *Parser) ParseMap(parent *types.Named, m *types.Map) (*schema.VarType, error) {
	key, err := p.ParseNamedType(parent, m.Key())
	if err != nil {
		return nil, fmt.Errorf("failed to parse map key type: %w", err)
	}

	value, err := p.ParseNamedType(parent, m.Elem())
	if err != nil {
		return nil, fmt.Errorf("failed to parse map value type: %w", err)
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
