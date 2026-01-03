package ast

import (
	"github.com/chrispyles/slow/src/execute"
	"github.com/chrispyles/slow/src/operators"
)

type BinaryOpNode struct {
	Op    *operators.BinaryOperator
	Left  execute.Expression
	Right execute.Expression
}

func (n *BinaryOpNode) Execute(e *execute.Environment) (execute.Value, error) {
	le, err := n.Left.Execute(e)
	if err != nil {
		return nil, err
	}
	re, err := n.Right.Execute(e)
	if err != nil {
		return nil, err
	}
	val, err := n.Op.Value(le, re)
	if err != nil {
		return nil, err
	}
	if n.Op.IsReassignmentOperator() {
		switch left := n.Left.(type) {
		case *VariableNode:
			return e.Set(left.Name, val)
		case *AttributeNode:
			expr, err := left.Left.Execute(e)
			if err != nil {
				return nil, err
			}
			if err := expr.SetAttribute(left.Right, val); err != nil {
				return nil, err
			}
			return val, nil
		case *IndexNode:
			expr, err := left.Container.Execute(e)
			if err != nil {
				return nil, err
			}
			idx, err := left.Index.Execute(e)
			if err != nil {
				return nil, err
			}
			if err := expr.SetIndex(idx, val); err != nil {
				return nil, err
			}
			return val, nil
		default:
			panic("unexpected node type in reassignment operator")
		}
	}
	return val, nil
}
