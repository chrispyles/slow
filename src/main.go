package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/chrispyles/slow/src/interpreter"
)

var (
	interpreterFlag = flag.Bool("i", false, "start the interpreter after running the file")
)

func main() {
	flag.Parse()
	if flag.NArg() > 1 {
		panic(fmt.Errorf("slow accepts 0 or 1 arguments, not %d", flag.NArg()))
	}

	var code []byte
	if flag.NArg() == 1 {
		var err error
		code, err = os.ReadFile(flag.Arg(0))
		if err != nil {
			panic(err)
		}
	}

	interactive := flag.NArg() == 0 || *interpreterFlag
	var rdr io.Reader
	if interactive {
		rdr = os.Stdin
	}

	interpreter.Run(string(code), rdr)
}
