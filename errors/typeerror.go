package errors

import "fmt"

type TypeError struct {
	from    string
	to      string
	message string
}

func NewTypeError(f Type, t Type) *TypeError {
	return &TypeError{from: f.String(), to: t.String()}
}

func IncompatibleType(t Type, op string) *TypeError {
	return TypeErrorFromMessage(fmt.Sprintf("TypeError: type %q cannot be used with operator %q", t.String(), op))
}

func IncompatibleTypes(l, r Type, op string) *TypeError {
	return TypeErrorFromMessage(fmt.Sprintf("TypeError: types %q and %q cannot be used together with operator %q", l.String(), r.String(), op))
}

func IncomparableType(t Type, op string) *TypeError {
	return TypeErrorFromMessage(fmt.Sprintf("TypeError: type %q cannot be used with comparison operator %q", t.String(), op))
}

func CallError(name string, got, want int) *TypeError {
	return TypeErrorFromMessage(fmt.Sprintf("TypeError: function %s accepts %d arguments but %d were given", name, want, got))
}

func NoLengthError(t Type) *TypeError {
	return TypeErrorFromMessage(fmt.Sprintf("TypeError: type %q does not have a length", t.String()))
}

func UnhashableTypeError(t Type) *TypeError {
	return TypeErrorFromMessage(fmt.Sprintf("TypeError: type %q is not hashable", t.String()))
}

func TypeErrorFromMessage(msg string) *TypeError {
	return &TypeError{message: msg}
}

func (e *TypeError) Error() string {
	if e.message != "" {
		return e.message
	}
	return fmt.Sprintf("TypeError: type %q cannot be used as type %q", e.from, e.to)
}
