package builtins

import (
	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/types"
)

func rangeImpl(args ...execute.Value) (execute.Value, error) {
	if len(args) == 0 {
		return nil, errors.CallError("range", len(args), 1)
	}
	var lower execute.Value
	var step execute.Value
	upper := args[0]
	if len(args) > 1 {
		lower = upper
		upper = args[1]
	}
	if len(args) > 2 {
		step = args[2]
	}
	if len(args) > 3 {
		return nil, errors.CallError("range", len(args), 3)
	}
	return types.NewRangeGenerator(lower, upper, step)
}
