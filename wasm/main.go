package main

import (
	"flag"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/chrispyles/slow/eval"
	"github.com/chrispyles/slow/interpreter"
	"github.com/chrispyles/slow/printer"
	"github.com/chrispyles/slow/reader"
)

var (
	interpreterFlag = flag.Bool("i", false, "start the interpreter after running the file")
)

func main() {
	env := interpreter.Run(string(""), nil)

	var out string
	printer.Set(func(s string) {
		out += s
	})

	js.Global().Set("evalSlow", js.FuncOf(func(_ js.Value, args []js.Value) any {
		in := args[0].String()
		if strings.TrimSpace(in) == "" {
			return ""
		}
		// Check that the input doesn't have any unmatched/unclosed characters.
		if _, err := reader.IsCompleteStatement(in); err != nil {
			return fmt.Sprintf("%+v", err)
		}
		eval.Eval(in, env, true)
		// Reset out after it's value is retrieved.
		defer func() { out = "" }()
		return out
	}))

	// Open a channel and block on reading from it to keep the binary running, otherwise it exits
	// and no code can be executed.
	c := make(chan int)
	<-c
}
