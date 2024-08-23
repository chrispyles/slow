package builtins

import (
	"github.com/chrispyles/slow/builtins/modules"
	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/types"
)

func importImpl(args ...execute.Value) (execute.Value, error) {
	if len(args) != 1 {
		return nil, errors.CallError("import", len(args), 1)
	}
	if _, ok := args[0].(*types.Str); !ok {
		return nil, errors.NewTypeError(args[0].Type(), types.StrType)
	}
	// TODO: relative imports
	name, err := args[0].ToStr()
	if err != nil {
		return nil, err
	}
	m, ok := modules.Get(name)
	if !ok {
		return nil, errors.NewImportError(name)
	}
	env, err := m.Import()
	if err != nil {
		return nil, err
	}
	return types.NewModule(name, env), nil
}
