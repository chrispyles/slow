package builtins

import (
	"os"
	"strings"

	"github.com/chrispyles/slow/builtins/modules"
	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/eval"
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
	if strings.HasSuffix(name, ".slo") {
		return importFile(name)
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

func importFile(path string) (execute.Value, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err // TODO: wrap error
	}
	env := RootEnvironment.NewFrame()
	eval.Eval(string(bytes), env, false)
	return types.NewModule(path, env), nil
}
