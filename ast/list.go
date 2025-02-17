package ast

import (
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/types"
)

type ListNode struct {
	Values []execute.Expression
}

func (n *ListNode) Execute(e *execute.Environment) (execute.Value, error) {
	vs := make([]execute.Value, len(n.Values))
	for i, expr := range n.Values {
		v, err := expr.Execute(e)
		if err != nil {
			return nil, err
		}
		vs[i] = v
	}
	return types.NewList(vs), nil
}
