package ast

import (
	"github.com/chrispyles/slow/src/execute"
	"github.com/chrispyles/slow/src/operators"
)

type UnaryOpNode struct {
	Op   *operators.UnaryOperator
	Expr execute.Expression
}

func (n *UnaryOpNode) Execute(e *execute.Environment) (execute.Value, error) {
	operand, err := n.Expr.Execute(e)
	if err != nil {
		return nil, err
	}
	val, err := n.Op.Value(operand)
	if err != nil {
		return nil, err
	}
	if n.Op.IsReassignmentOperator() {
		// Reassignment operators return the value of the operand BEFORE the operation, but update its
		// value in the environment/object.
		switch expr := n.Expr.(type) {
		case *VariableNode:
			_, err := e.Set(expr.Name, val)
			if err != nil {
				return nil, err
			}
			return operand, nil
		case *AttributeNode:
			lVal, err := expr.Left.Execute(e)
			if err != nil {
				return nil, err
			}
			if err := lVal.SetAttribute(expr.Right, val); err != nil {
				return nil, err
			}
			return operand, nil
		case *IndexNode:
			lVal, err := expr.Container.Execute(e)
			if err != nil {
				return nil, err
			}
			idx, err := expr.Index.Execute(e)
			if err != nil {
				return nil, err
			}
			if err := lVal.SetIndex(idx, val); err != nil {
				return nil, err
			}
			return operand, nil
		default:
			panic("unexpected node type in reassignment operator")
		}
	}
	return val, nil
}
