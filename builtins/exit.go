package builtins

import (
	"os"

	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/printer"
	"github.com/chrispyles/slow/types"
)

func exitImpl(args ...execute.Value) (execute.Value, error) {
	var code int64
	if len(args) > 0 {
		code, _ = args[0].ToInt()
	}
	printer.Printlnf("Exiting with code %d", code)
	os.Exit(int(code))
	return types.Null, nil
}
