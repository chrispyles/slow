package printer

import "fmt"

var doPrint func(s string) = defaultPrint

func defaultPrint(s string) {
	fmt.Print(s)
}

func Set(dst func(s string)) {
	if dst == nil {
		dst = defaultPrint
	}
	doPrint = dst
}

func Print(s string) {
	doPrint(s)
}

func Println(s string) {
	Print(s + "\n")
}

func Printf(s string, args ...any) {
	doPrint(fmt.Sprintf(s, args...))
}

func Printlnf(s string, args ...any) {
	Printf(s+"\n", args...)
}
