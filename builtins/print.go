package builtins

import (
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/printer"
	"github.com/chrispyles/slow/types"
)

var println = printer.Println

func printImpl(args ...execute.Value) (execute.Value, error) {
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
	println(fullout)
	return types.Null, nil
}
