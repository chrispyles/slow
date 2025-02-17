package ast

import (
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/types"
)

type VarNode struct {
	Name    string
	IsConst bool
	// Value of the expression if it is assigned; may be nil
	Value execute.Expression
}

func (n *VarNode) Execute(e *execute.Environment) (execute.Value, error) {
	var val execute.Value
	if n.Value != nil {
		var err error
		val, err = n.Value.Execute(e)
		if err != nil {
			return nil, err
		}
	}
	if n.IsConst {
		if val == nil {
			panic("encountered const node with no value")
		}
		return e.DeclareConst(n.Name, val)
	}
	if err := e.Declare(n.Name); err != nil {
		return nil, err
	}
	if val != nil {
		return e.Set(n.Name, val)
	}
	return types.Null, nil
}
