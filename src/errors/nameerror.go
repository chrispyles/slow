package errors

import "fmt"

func NewNameError(name string) error {
	return newError("NameError", fmt.Sprintf("no variable %q has been declared", name))
}
