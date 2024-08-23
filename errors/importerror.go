package errors

import "fmt"

type ImportError struct {
	name string
}

func NewImportError(name string) *ImportError {
	return &ImportError{name: name}
}

func (e *ImportError) Error() string {
	return fmt.Sprintf("ImportError: no such module %q", e.name)
}
