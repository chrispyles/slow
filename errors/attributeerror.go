package errors

import "fmt"

type AttributeError struct {
	type_ string
	name  string
}

func NewAttributeError(t Type, n string) *AttributeError {
	return &AttributeError{type_: t.String(), name: n}
}

func (e *AttributeError) Error() string {
	return fmt.Sprintf("AttributeError: type %q has no attribute %q", e.type_, e.name)
}
