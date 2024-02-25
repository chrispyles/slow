package errors

import "fmt"

type NameError struct {
	name string
}

func NewNameError(n string) *NameError {
	return &NameError{name: n}
}

func (e *NameError) Error() string {
	return fmt.Sprintf("NameError: no such variable %q", e.name)
}
