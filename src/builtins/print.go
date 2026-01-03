package builtins

import (
	"github.com/chrispyles/slow/src/execute"
	"github.com/chrispyles/slow/src/printer"
	"github.com/chrispyles/slow/src/types"
)

func printImpl(args ...execute.Value) (execute.Value, error) {
	var fullout string
	for _, v := range args {
		var out string
		if s, ok := v.(*types.Str); ok {
			// The Str.ToStr method returns the value without the delimiting quotes, so we use it here
			// so as not to print the quotes when printing string values.
			out = s.Value()
		} else {
			out = v.String()
		}
		fullout += out
	}
	printer.Println(fullout)
	return types.Null, nil
}
