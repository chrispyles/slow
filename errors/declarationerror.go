package errors

import "fmt"

type DeclarationError struct {
	name string
}

func NewDeclarationError(n string) *DeclarationError {
	return &DeclarationError{name: n}
}

func (e *DeclarationError) Error() string {
	return fmt.Sprintf("DeclarationError: variable %q has already been declared", e.name)
}
