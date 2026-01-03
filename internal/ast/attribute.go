package ast

import (
	"github.com/chrispyles/slow/internal/execute"
)

type AttributeNode struct {
	Left  execute.Expression
	Right string
}

func (n *AttributeNode) Execute(e *execute.Environment) (execute.Value, error) {
	expr, err := n.Left.Execute(e)
	if err != nil {
		return nil, err
	}
	return expr.GetAttribute(n.Right)
}
