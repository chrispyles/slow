package eval

import (
	"fmt"

	"github.com/chrispyles/slow/config"
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/parser"
	"github.com/chrispyles/slow/printer"
	"github.com/chrispyles/slow/types"
	"github.com/sanity-io/litter"
)

var (
	makeAST = parser.NewAST
	println = printer.Println
)

func Eval(s string, env *execute.Environment, printExprValue bool) {
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

	if printExprValue && val != types.Null {
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
