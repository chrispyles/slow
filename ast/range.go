package ast

import (
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/types"
)

type RangeNode struct {
	Start execute.Expression
	Stop  execute.Expression
	Step  execute.Expression
}

func (n *RangeNode) Execute(e *execute.Environment) (execute.Value, error) {
	var start, stop, step execute.Value
	var err error
	if n.Start != nil {
		start, err = n.Start.Execute(e)
		if err != nil {
			return nil, err
		}
	} else {
		start = types.NewUint(0)
	}
	stop, err = n.Stop.Execute(e)
	if err != nil {
		return nil, err
	}
	if n.Step != nil {
		step, err = n.Step.Execute(e)
		if err != nil {
			return nil, err
		}
	} else {
		step = types.NewUint(1)
	}
	return types.NewRangeGenerator(start, stop, step)
}
