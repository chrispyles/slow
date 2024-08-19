package printer

import "fmt"

const (
	prefix = ""
)

func Print(s string) {
	fmt.Print(prefix + s)
}

func Println(s string) {
	Print(s + "\n")
}

func Printf(s string, args ...any) {
	fmt.Printf(prefix+s, args...)
}

func Printlnf(s string, args ...any) {
	Printf(s+"\n", args...)
}
