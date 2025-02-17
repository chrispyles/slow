package ast

import (
	"github.com/chrispyles/slow/execute"
)

type ConstantNode struct {
	Value execute.Value
}

func (n *ConstantNode) Execute(e *execute.Environment) (execute.Value, error) {
	return n.Value, nil
}
