package builtins

import (
	"github.com/chrispyles/slow/internal/errors"
	"github.com/chrispyles/slow/internal/execute"
	"github.com/chrispyles/slow/internal/types"
)

func lenImpl(args ...execute.Value) (execute.Value, error) {
	if len(args) != 1 {
		return nil, errors.CallError("len", len(args), 1)
	}
	l, err := args[0].Length()
	if err != nil {
		return nil, err
	}
	return types.NewUint(l), nil
}
