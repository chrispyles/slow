package modules

import (
	"os"

	"github.com/chrispyles/slow/internal/errors"
	"github.com/chrispyles/slow/internal/execute"
	"github.com/chrispyles/slow/internal/types"
)

type fsModule struct{}

func fs_read(args ...execute.Value) (execute.Value, error) {
	if len(args) != 1 {
		return nil, errors.CallError("fs.read", len(args), 1)
	}
	bytes, err := fs_readBytes(args...)
	if err != nil {
		return nil, err
	}
	str, err := bytes.ToStr()
	if err != nil {
		return nil, err
	}
	return types.NewStr(str), nil
}

func fs_readBytes(args ...execute.Value) (execute.Value, error) {
	if len(args) != 1 {
		return nil, errors.CallError("fs.readBytes", len(args), 1)
	}
	path, err := args[0].ToStr()
	if err != nil {
		return nil, err
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.WrapFileError(err, path)
	}
	return types.NewBytes(bytes), nil
}

var functions = map[string]types.FuncImpl{
	"read":      fs_read,
	"readBytes": fs_readBytes,
}

func (m *fsModule) Name() string {
	return "fs"
}

func (m *fsModule) Import() (*execute.Environment, error) {
	fns := make(map[string]execute.Value)
	for name, impl := range functions {
		fns[name] = types.NewGoFunc(name, impl)
	}
	return execute.FromMap(fns), nil
}
