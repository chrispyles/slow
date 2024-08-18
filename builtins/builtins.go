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
				var out string
				if s, ok := v.(*types.Str); ok {
					// The Str.ToStr method returns the value without the delimiting quotes, so we use it here
					// so as not to print the quotes when printing string values.
					out, _ = s.ToStr()
				} else {
					out = v.String()
				}
				fmt.Print(out)
			}
			fmt.Print("\n")
			return types.Null, nil
		},
	},
	{
		name: "range",
		f: func(vs ...execute.Value) (execute.Value, error) {
			// TODO: make this into a generator once available so all values aren't coalesced
			if len(vs) == 0 {
				return nil, errors.CallError("range", len(vs), 1)
			}
			var lower execute.Value
			var step execute.Value
			upper := vs[0]
			if len(vs) > 1 {
				lower = upper
				upper = vs[1]
			} else {
				lower = types.NewUint(0)
			}
			if len(vs) > 2 {
				step = vs[2]
			} else {
				step = types.NewUint(1)
			}
			if len(vs) > 3 {
				return nil, errors.CallError("range", len(vs), 3)
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
		f: func(vs ...execute.Value) (execute.Value, error) {
			if len(vs) != 1 {
				return nil, errors.CallError("type", len(vs), 1)
			}
			return types.NewStr(vs[0].Type().String()), nil
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
