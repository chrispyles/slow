package interpreter

import (
	"bufio"
	"io"

	"github.com/chrispyles/slow/src/builtins"
	evallib "github.com/chrispyles/slow/src/eval"
	"github.com/chrispyles/slow/src/execute"
	"github.com/chrispyles/slow/src/printer"
	"github.com/chrispyles/slow/src/reader"
)

var (
	read = reader.Read
	eval = evallib.Eval
)

func Run(code string, interactiveReader io.Reader) *execute.Environment {
	env := builtins.RootEnvironment.NewFrame()
	if code != "" {
		eval(code, env, false)
	}

	if interactiveReader == nil {
		return env
	}

	rdr := bufio.NewReader(interactiveReader)
	for {
		stmt, err := read(rdr)
		if err != nil {
			printError(err)
			continue
		}
		if stmt == "\n" {
			// Don't attempt to execute an empty line
			continue
		}
		eval(stmt, env, true)
	}
}

func printError(err error) {
	printer.Printlnf("%+v", err)
}
