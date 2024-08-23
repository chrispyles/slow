package interpreter

import (
	"bufio"
	"os"

	"github.com/chrispyles/slow/builtins"
	"github.com/chrispyles/slow/eval"
	"github.com/chrispyles/slow/printer"
	"github.com/chrispyles/slow/reader"
)

func Run(code string, interactive bool) {
	env := builtins.RootEnvironment.NewFrame()
	if code != "" {
		eval.Eval(code, env, false)
	}

	if !interactive {
		return
	}

	rdr := bufio.NewReader(os.Stdin)
	for {
		stmt, err := reader.Read(rdr)
		if err != nil {
			printError(err)
			continue
		}
		if stmt == "\n" {
			// Don't attempt to execute an empty line
			continue
		}
		eval.Eval(stmt, env, true)
	}
}

func printError(err error) {
	printer.Printlnf("%+v", err)
}
