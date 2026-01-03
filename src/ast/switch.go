package ast

import (
	"github.com/chrispyles/slow/src/execute"
)

type SwitchCase struct {
	CaseExpr execute.Expression
	Body     execute.Block
}

type SwitchNode struct {
	Value       execute.Expression
	Cases       []SwitchCase
	DefaultCase execute.Block
}

func (n *SwitchNode) Execute(e *execute.Environment) (execute.Value, error) {
	valueExpr, err := n.Value.Execute(e)
	if err != nil {
		return nil, err
	}
	frame := e.NewFrame()
	var fallThrough bool
	for _, c := range n.Cases {
		caseExpr, err := c.CaseExpr.Execute(e)
		if err != nil {
			return nil, err
		}
		if valueExpr.Equals(caseExpr) || fallThrough {
			fallThrough = false
			val, err := c.Body.Execute(frame)
			if _, ok := err.(*fallthroughError); ok {
				fallThrough = true
				continue
			} else {
				return val, err
			}
		}
	}
	return n.DefaultCase.Execute(frame)
}
