package errors

import "fmt"

func NewTypeError(from, to Type) error {
	return newError("TypeError", fmt.Sprintf("type %q cannot be used as type %q", from, to))
}

func IncompatibleType(t Type, op string) error {
	return newError("TypeError", fmt.Sprintf("type %q cannot be used with operator %q", t.String(), op))
}

func IncompatibleTypes(l, r Type, op string) error {
	return newError("TypeError", fmt.Sprintf("types %q and %q cannot be used together with operator %q", l.String(), r.String(), op))
}

func IncomparableType(t Type, op string) error {
	return newError("TypeError", fmt.Sprintf("type %q cannot be used with comparison operator %q", t.String(), op))
}

func CallError(name string, got, want int) error {
	return newError("TypeError", fmt.Sprintf("function %s accepts %d arguments but %d were given", name, want, got))
}

func NoLengthError(t Type) error {
	return newError("TypeError", fmt.Sprintf("type %q does not have a length", t.String()))
}

func UnhashableTypeError(t Type) error {
	return newError("TypeError", fmt.Sprintf("type %q is not hashable", t.String()))
}

func InvalidTypeCastTarget(dst Type) error {
	return newError("TypeError", fmt.Sprintf("type %q does not support type casting", dst.String()))
}

func TypeErrorFromMessage(msg string) error {
	return newError("TypeError", msg)
}
