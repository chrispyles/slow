package builtins

import (
	"os"

	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/printer"
	"github.com/chrispyles/slow/types"
)

var builtins = []struct {
	name string
	f    types.FuncImpl
}{
	{
		name: "exit",
		f: func(args ...execute.Value) (execute.Value, error) {
			var code int64
			if len(args) > 0 {
				code, _ = args[0].ToInt()
			}
			printer.Printlnf("Exiting with code %d", code)
			os.Exit(int(code))
			return types.Null, nil
		},
	},
	{
		name: "len",
		f: func(args ...execute.Value) (execute.Value, error) {
			if len(args) != 1 {
				return nil, errors.CallError("len", len(args), 1)
			}
			l, err := args[0].Length()
			if err != nil {
				return nil, err
			}
			return types.NewUint(l), nil
		},
	},
	{
		name: "print",
		f: func(args ...execute.Value) (execute.Value, error) {
			var fullout string
			for _, v := range args {
				var out string
				if s, ok := v.(*types.Str); ok {
					// The Str.ToStr method returns the value without the delimiting quotes, so we use it here
					// so as not to print the quotes when printing string values.
					out, _ = s.ToStr()
				} else {
					out = v.String()
				}
				fullout += out
			}
			printer.Println(fullout)
			return types.Null, nil
		},
	},
	{
		name: "range",
		f: func(args ...execute.Value) (execute.Value, error) {
			if len(args) == 0 {
				return nil, errors.CallError("range", len(args), 1)
			}
			var lower execute.Value
			var step execute.Value
			upper := args[0]
			if len(args) > 1 {
				lower = upper
				upper = args[1]
			} else {
				lower = types.NewUint(0)
			}
			if len(args) > 2 {
				step = args[2]
			} else {
				step = types.NewUint(1)
			}
			if len(args) > 3 {
				return nil, errors.CallError("range", len(args), 3)
			}
			rg, err := newRangeGenerator(lower, upper, step)
			if err != nil {
				return nil, err
			}
			return types.NewGenerator(rg), nil
		},
	},
	{
		name: "type",
		f: func(args ...execute.Value) (execute.Value, error) {
			if len(args) != 1 {
				return nil, errors.CallError("type", len(args), 1)
			}
			return types.NewStr(args[0].Type().String()), nil
		},
	},
}

// NewRootEnvironment creates a new root environment populated with every builtin. The returned
// environment is frozen, so a child frame must be created from it before declaring any variables.
func NewRootEnvironment() *execute.Environment {
	e := execute.NewEnvironment()
	for _, b := range builtins {
		f := types.NewGoFunc(b.name, b.f)
		e.Declare(b.name)
		e.Set(b.name, f)
	}
	// freeze the root environment so nothing is inadvertently overwritten
	e.Freeze()
	return e
}
