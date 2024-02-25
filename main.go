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
		panic(fmt.Errorf("ast.Execute returned nil"))
	}

	if printOut && val != types.Null {
		print(val.String(), "\n")
	}

	if *config.Debug {
		fmt.Print("\nAST evaluated to: ")
		litter.Dump(val)
	}
}

func main() {
	flag.Parse()
	if flag.NArg() > 1 {
		// TODO: error bc there can only be 1 file run at a time
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
		fmt.Println("this is where I would put an interpreter... IF I HAD ONE")

		rdr := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("whisper to me> ")
			line, err := rdr.ReadString('\n')
			if err != nil {
				// TODO: handle more gracefully
				panic(err)
			}

			eval(line, frame, true)
		}
	}
}

func printError(err error) {
	fmt.Printf("%+v\n", err)
}
