package builtins

import (
	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/types"
)

func typeImpl(args ...execute.Value) (execute.Value, error) {
	if len(args) != 1 {
		return nil, errors.CallError("type", len(args), 1)
	}
	return types.NewStr(args[0].Type().String()), nil
}
