package builtins

import (
	"github.com/chrispyles/slow/internal/errors"
	"github.com/chrispyles/slow/internal/execute"
	"github.com/chrispyles/slow/internal/types"
)

func typeImpl(args ...execute.Value) (execute.Value, error) {
	if len(args) != 1 {
		return nil, errors.CallError("type", len(args), 1)
	}
	return types.NewStr(args[0].Type().String()), nil
}
