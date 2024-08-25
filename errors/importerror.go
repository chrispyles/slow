package errors

import "fmt"

func NewImportError(name string) error {
	return newError("ImportError", fmt.Sprintf("no such module %q", name))
}
