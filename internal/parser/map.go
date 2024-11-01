package parser

import (
	"fmt"
	"go/types"

	"github.com/webrpc/webrpc/schema"
)

func (p *Parser) ParseMap(typeName string, m *types.Map) (*schema.VarType, error) {
	key, err := p.ParseNamedType(typeName, m.Key())
	if err != nil {
		return nil, fmt.Errorf("failed to parse map key type: %w", err)
	}

	value, err := p.ParseNamedType(typeName, m.Elem())
	if err != nil {
		return nil, fmt.Errorf("failed to parse map value type: %w", err)
	}

	varType := &schema.VarType{
		Expr: fmt.Sprintf("map<%v,%v>", key, value),
		Type: schema.T_Map,
		Map: &schema.VarMapType{
			Key:   key,
			Value: value,
		},
	}

	return varType, nil
}
