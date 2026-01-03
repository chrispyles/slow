package ast

import (
	"github.com/chrispyles/slow/src/execute"
	"github.com/chrispyles/slow/src/types"
)

type MapNode struct {
	Values [][]execute.Expression
}

func (n *MapNode) Execute(e *execute.Environment) (execute.Value, error) {
	m := types.NewMap()
	for _, kv := range n.Values {
		k, err := kv[0].Execute(e)
		if err != nil {
			return nil, err
		}
		v, err := kv[1].Execute(e)
		if err != nil {
			return nil, err
		}
		if _, err := m.Set(k, v); err != nil {
			return nil, err
		}
	}
	return m, nil
}
