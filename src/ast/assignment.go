package ast

import (
	"github.com/chrispyles/slow/src/execute"
)

type AssignmentTarget struct {
	Variable  string
	Attribute *AttributeNode
	Index     *IndexNode
}

type AssignmentNode struct {
	Left  AssignmentTarget
	Right execute.Expression
}

func (n *AssignmentNode) Execute(e *execute.Environment) (execute.Value, error) {
	expr, err := n.Right.Execute(e)
	if err != nil {
		return nil, err
	}
	if n := n.Left.Variable; n != "" {
		return e.Set(n, expr)
	}
	if an := n.Left.Attribute; an != nil {
		val, err := an.Left.Execute(e)
		if err != nil {
			return nil, err
		}
		err = val.SetAttribute(an.Right, expr)
		if err != nil {
			return nil, err
		}
		return expr, nil
	}
	if in := n.Left.Index; in != nil {
		val, err := in.Container.Execute(e)
		if err != nil {
			return nil, err
		}
		idx, err := in.Index.Execute(e)
		if err != nil {
			return nil, err
		}
		if err := val.SetIndex(idx, expr); err != nil {
			return nil, err
		}
		return val, nil
	}
	panic("unhandled target case in AssignmentNode.Execute")
}
