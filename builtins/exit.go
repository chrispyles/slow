package builtins

import (
	"os"

	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/types"
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
	printlnf("Exiting with code %d", code)
	osExit(int(code))
	return types.Null, nil
}
