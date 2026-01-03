package errors

import "fmt"

func NewValueError(msg string) error {
	return newError("ValueError", msg)
}

func WrapValueError(val string, toType Type, err error) error {
	return wrapError("ValueError", fmt.Sprintf("error converting %q to type %q: %+v", val, toType.String(), err), err)
}
