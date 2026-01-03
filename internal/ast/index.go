package ast

import (
	"github.com/chrispyles/slow/internal/execute"
)

type IndexNode struct {
	Container execute.Expression
	Index     execute.Expression
}

func (n *IndexNode) Execute(e *execute.Environment) (execute.Value, error) {
	container, err := n.Container.Execute(e)
	if err != nil {
		return nil, err
	}
	index, err := n.Index.Execute(e)
	if err != nil {
		return nil, err
	}
	return container.GetIndex(index)
}
