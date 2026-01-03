package ast

import (
	"github.com/chrispyles/slow/src/execute"
	"github.com/chrispyles/slow/src/types"
)

type DeferNode struct {
	Expr execute.Expression
}

func (n *DeferNode) Execute(e *execute.Environment) (execute.Value, error) {
	return nil, &types.DeferError{Expr: n.Expr}
}
