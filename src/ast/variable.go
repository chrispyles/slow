package ast

import (
	"github.com/chrispyles/slow/src/execute"
)

// VariableNode represents a variable access, not a declaration.
type VariableNode struct {
	Name string
}

func (n *VariableNode) Execute(e *execute.Environment) (execute.Value, error) {
	return e.Get(n.Name)
}
