package ast

import (
	"github.com/chrispyles/slow/src/execute"
)

type ForNode struct {
	IterName string
	Iter     execute.Expression
	Body     execute.Block
}

func (n *ForNode) Execute(e *execute.Environment) (execute.Value, error) {
	var val execute.Value
	expr, err := n.Iter.Execute(e)
	if err != nil {
		return nil, err
	}
	iter, err := expr.ToIterator()
	if err != nil {
		return nil, err
	}
	for iter.HasNext() {
		frame := e.NewFrame()
		if err := frame.Declare(n.IterName); err != nil {
			return nil, err
		}
		expr, err := iter.Next()
		if err != nil {
			return nil, err
		}
		if _, err := frame.Set(n.IterName, expr); err != nil {
			return nil, err
		}
		val, err = n.Body.Execute(frame)
		if _, ok := err.(*breakError); ok {
			break
		} else if _, ok := err.(*continueError); ok {
			continue
		} else if err != nil {
			return nil, err
		}
	}
	return val, nil
}
