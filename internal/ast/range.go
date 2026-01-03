package ast

import (
	"github.com/chrispyles/slow/internal/execute"
	"github.com/chrispyles/slow/internal/types"
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
	}
	if n.Stop != nil {
		stop, err = n.Stop.Execute(e)
		if err != nil {
			return nil, err
		}
	}
	if n.Step != nil {
		step, err = n.Step.Execute(e)
		if err != nil {
			return nil, err
		}
	}
	return types.NewRangeGenerator(start, stop, step)
}
