package errors

import "fmt"

func NewDeclarationError(name string) error {
	return newError("DeclarationError", fmt.Sprintf("variable %q has already been declared", name))
}
