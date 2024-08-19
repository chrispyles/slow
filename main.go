package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/chrispyles/slow/builtins"
	"github.com/chrispyles/slow/config"
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/parser"
	"github.com/chrispyles/slow/printer"
	"github.com/chrispyles/slow/reader"
	"github.com/chrispyles/slow/types"
	"github.com/sanity-io/litter"
)

var (
	interpreter = flag.Bool("i", false, "start the interpreter after running the file")
)

func eval(s string, env *execute.Environment, printOut bool) {
	ast, err := parser.NewAST(s)
	if err != nil {
		printError(err)
		return
	}

	astString := ast.String()
	if *config.Debug {
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

	if printOut && val != types.Null {
		printer.Println(val.String())
	}

	if *config.Debug {
		fmt.Print("\nAST evaluated to: ")
		litter.Dump(val)
	}
}

func main() {
	flag.Parse()
	if flag.NArg() > 1 {
		panic(fmt.Errorf("slow accepts 0 or 1 arguments, not %d", flag.NArg()))
	}

	env := builtins.NewRootEnvironment()
	frame := env.NewFrame()
	if flag.NArg() == 1 {
		code, err := os.ReadFile(flag.Arg(0))
		if err != nil {
			panic(err)
		}

		eval(string(code), frame, false)
	}

	if flag.NArg() == 0 || *interpreter {
		rdr := bufio.NewReader(os.Stdin)
		for {
			stmt, err := reader.Read(rdr)
			if err != nil {
				printError(err)
				continue
			}
			eval(stmt, frame, true)
		}
	}
}

func printError(err error) {
	printer.Printlnf("%+v", err)
}
