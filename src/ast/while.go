package ast

import (
	"github.com/chrispyles/slow/src/execute"
)

type WhileNode struct {
	Cond execute.Expression
	Body execute.Block
}

func (n *WhileNode) Execute(e *execute.Environment) (execute.Value, error) {
	frame := e.NewFrame()
	var val execute.Value
	for {
		expr, err := n.Cond.Execute(e)
		if err != nil {
			return nil, err
		}
		if expr.ToBool() {
			val, err = n.Body.Execute(frame)
			if _, ok := err.(*breakError); ok {
				break
			} else if _, ok := err.(*continueError); ok {
				continue
			} else if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	return val, nil
}
