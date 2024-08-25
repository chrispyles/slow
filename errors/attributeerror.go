package errors

import "fmt"

func NewAttributeError(type_ Type, name string) error {
	return newError("AttributeError", fmt.Sprintf("type %q has no attribute %q", type_, name))
}
