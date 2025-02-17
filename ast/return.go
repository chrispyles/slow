package ast

import (
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/types"
)

type ReturnNode struct {
	Value execute.Expression
}

func (n *ReturnNode) Execute(e *execute.Environment) (execute.Value, error) {
	value, err := n.Value.Execute(e)
	if err != nil {
		return nil, err
	}
	return nil, &types.ReturnError{Value: value}
}
