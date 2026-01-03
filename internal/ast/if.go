package ast

import (
	"github.com/chrispyles/slow/internal/execute"
)

type IfNode struct {
	Cond     execute.Expression
	Body     execute.Block
	ElseBody execute.Block
}

func (n *IfNode) Execute(e *execute.Environment) (execute.Value, error) {
	expr, err := n.Cond.Execute(e)
	if err != nil {
		return nil, err
	}
	frame := e.NewFrame()
	if expr.ToBool() {
		return n.Body.Execute(frame)
	} else {
		return n.ElseBody.Execute(frame)
	}
}
