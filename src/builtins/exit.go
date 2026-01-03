package builtins

import (
	"os"

	"github.com/chrispyles/slow/src/execute"
	"github.com/chrispyles/slow/src/printer"
	"github.com/chrispyles/slow/src/types"
)

var (
	osExit = os.Exit
)

func exitImpl(args ...execute.Value) (execute.Value, error) {
	var code int64
	if len(args) > 0 {
		var err error
		code, err = args[0].ToInt()
		if err != nil {
			code, _ = types.NewBool(args[0].ToBool()).ToInt()
		}
	}
	printer.Printlnf("Exiting with code %d", code)
	osExit(int(code))
	return types.Null, nil
}
