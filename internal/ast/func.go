package ast

import (
	"github.com/chrispyles/slow/internal/execute"
	"github.com/chrispyles/slow/internal/types"
)

type FuncNode struct {
	Name     string
	ArgNames []string
	Body     execute.Block
}

func (n *FuncNode) Execute(e *execute.Environment) (execute.Value, error) {
	if err := e.Declare(n.Name); err != nil {
		return nil, err
	}
	ft := types.NewFunc(n.Name, n.ArgNames, n.Body)
	return e.Set(n.Name, ft)
}
