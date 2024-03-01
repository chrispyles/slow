package builtins

import (
	"fmt"
	"os"

	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/types"
)

var builtins = []struct {
	name string
	f    types.FuncImpl
}{
	{
		name: "exit",
		f: func(vs ...execute.Value) (execute.Value, error) {
			var code int64
			if len(vs) > 0 {
				code, _ = vs[0].ToInt()
			}
			fmt.Printf("Exiting with code %d\n", code)
			os.Exit(int(code))
			return types.Null, nil
		},
	},
	{
		name: "len",
		f: func(vs ...execute.Value) (execute.Value, error) {
			if len(vs) != 1 {
				return nil, errors.CallError("len", len(vs), 1)
			}
			l, err := vs[0].Length()
			if err != nil {
				return nil, err
			}
			return types.NewUint(l), nil
		},
	},
	{
		name: "print",
		f: func(vs ...execute.Value) (execute.Value, error) {
			for _, v := range vs {
				fmt.Print(v.String())
			}
			fmt.Print("\n")
			return types.Null, nil
		},
	},
	{
		name: "range",
		f: func(vs ...execute.Value) (execute.Value, error) {
			// TODO: make this into a generator once available so all values aren't coalesced
			// TODO: support step argument
			lower := int64(0)
			upper, err := vs[0].ToInt()
			if err != nil {
				return nil, err
			}
			if len(vs) > 1 {
				lower = upper
				upper, err = vs[1].ToInt()
				if err != nil {
					return nil, err
				}
			}
			values := []execute.Value{}
			for i := lower; i < upper; i++ {
				values = append(values, types.NewInt(i))
			}
			return types.NewList(values), nil
		},
	},
	{
		name: "type",
		f: func(vs ...execute.Value) (execute.Value, error) {
			if len(vs) != 1 {
				return nil, errors.CallError("type", len(vs), 1)
			}
			return types.NewStr(vs[0].Type().String()), nil
		},
	},
}

// NewRootEnvironment creates a new root environment populated with every builtin.
func NewRootEnvironment() *execute.Environment {
	e := execute.NewEnvironment()
	for _, b := range builtins {
		f := types.NewGoFunc(b.name, b.f)
		e.Declare(b.name)
		e.Set(b.name, f)
	}
	return e
}
