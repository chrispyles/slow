package eval

import (
	"fmt"

	"github.com/chrispyles/slow/src/config"
	"github.com/chrispyles/slow/src/execute"
	"github.com/chrispyles/slow/src/parser"
	"github.com/chrispyles/slow/src/printer"
	"github.com/chrispyles/slow/src/types"
	"github.com/sanity-io/litter"
)

var (
	makeAST = parser.Parse
	println = printer.Println
)

func Eval(s string, env *execute.Environment, printExpr bool) {
	ast, err := makeAST(s)
	if err != nil {
		printError(err)
		return
	}

	if *config.Debug {
		astString := ast.String()
		fmt.Println("<AST> ", astString)
	}

	val, err := ast.Execute(env)
	if err != nil {
		printError(err)
		return
	}

	// A nil value should never be returned by evaluating an expression unless err is non-nil.
	if val == nil {
		panic("ast.Execute returned nil")
	}

	if printExpr && val != types.Null {
		println(val.String())
	}

	if *config.Debug {
		fmt.Print("\nAST evaluated to: ")
		litter.Dump(val)
	}
}

func printError(err error) {
	printer.Printlnf("%+v", err)
}
