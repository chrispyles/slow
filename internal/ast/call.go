package ast

import (
	"github.com/chrispyles/slow/internal/execute"
)

type CallNode struct {
	Func execute.Expression
	Args []execute.Expression
}

func (n *CallNode) Execute(e *execute.Environment) (execute.Value, error) {
	expr, err := n.Func.Execute(e)
	if err != nil {
		return nil, err
	}
	callable, err := expr.ToCallable()
	if err != nil {
		return nil, err
	}
	var args []execute.Value
	for _, a := range n.Args {
		v, err := a.Execute(e)
		if err != nil {
			return nil, err
		}
		args = append(args, v)
	}
	return callable.Call(e.NewFrame(), args...)
}
