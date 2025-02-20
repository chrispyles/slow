package ast

import (
	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/types"
)

var castingUnsupportedTypes = map[execute.Type]bool{
	types.FuncType:      true,
	types.GeneratorType: true,
	types.IteratorType:  true,
	types.ListType:      true,
	types.MapType:       true,
	types.ModuleType:    true,
	types.NullType:      true,
}

type CastNode struct {
	Expr execute.Expression
	Type execute.Type
}

func (n *CastNode) Execute(e *execute.Environment) (execute.Value, error) {
	if castingUnsupportedTypes[n.Type] {
		return nil, errors.InvalidTypeCastTarget(n.Type)
	}
	val, err := n.Expr.Execute(e)
	if err != nil {
		return nil, err
	}
	return n.Type.New(val)
}
