package builtins

import (
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/printer"
	"github.com/chrispyles/slow/types"
)

var (
	println  = printer.Println
	printlnf = printer.Printlnf
)

// RootEnvironment is a frozen environment containing all builtins. All execution environments
// should be child frames of this environment.
var RootEnvironment *execute.Environment

var builtins = []struct {
	name string
	f    types.FuncImpl
}{
	{
		name: "exit",
		f:    exitImpl,
	},
	{
		name: "import",
		f:    importImpl,
	},
	{
		name: "len",
		f:    lenImpl,
	},
	{
		name: "print",
		f:    printImpl,
	},
	{
		name: "range",
		f:    rangeImpl,
	},
	{
		name: "type",
		f:    typeImpl,
	},
}

// newRootEnvironment creates a new root environment populated with every builtin. The returned
// environment is frozen, so a child frame must be created from it before declaring any variables.
func newRootEnvironment() *execute.Environment {
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

func init() {
	RootEnvironment = newRootEnvironment()
}
