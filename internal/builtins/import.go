package builtins

import (
	"os"
	"strings"

	"github.com/chrispyles/slow/internal/builtins/modules"
	"github.com/chrispyles/slow/internal/errors"
	"github.com/chrispyles/slow/internal/eval"
	"github.com/chrispyles/slow/internal/execute"
	"github.com/chrispyles/slow/internal/types"
)

var (
	evalEval   = eval.Eval
	osReadFile = os.ReadFile
)

func importImpl(args ...execute.Value) (execute.Value, error) {
	if len(args) != 1 {
		return nil, errors.CallError("import", len(args), 1)
	}
	argstr, ok := args[0].(*types.Str)
	if !ok {
		return nil, errors.NewTypeError(args[0].Type(), types.StrType)
	}
	name := argstr.Value()
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
	bytes, err := osReadFile(path)
	if err != nil {
		return nil, errors.WrapFileError(err, path)
	}
	env := RootEnvironment.NewFrame()
	evalEval(string(bytes), env, false)
	return types.NewModule(path, env), nil
}
