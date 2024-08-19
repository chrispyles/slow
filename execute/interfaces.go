package execute

type Callable interface {
	Call(*Environment, ...Value) (Value, error)
}

type Expression interface {
	Execute(e *Environment) (Value, error)
}

type Iterator interface {
	HasNext() bool
	Next() (Value, error)
}

type Type interface {
	IsNumeric() bool
	Matches(Type) bool
	String() string
}

type ComparisonErrorBuilder struct {
	Source Type
	Other  Type
}

func ComparisonError(s, o Type) *ComparisonErrorBuilder {
	return &ComparisonErrorBuilder{s, o}
}

type Value interface {
	// CloneIfPrimitive returns a copy of this value object if it is a primitive (i.e. pass-by-value).
	// Values that are pass-by-reference should return a reference to the same instance.
	CloneIfPrimitive() Value
	// A negative return value means this value is less than the other one, 0 means equal to, and
	// positive means greather than. For values that are pass-by-reference, or if the two types are
	// incomparable, this second return value should always be false. The logic for comparing the
	// references for ==/!= will be handled in the caller. The second return value should be true iff
	// the types are pass-by-value (primitive) and comparable.
	CompareTo(Value) (int, bool)
	// N.B. Equals() should return false if the value is not of the same type (i.e.
	// int(1) != float(1) for the purposes of this method).
	Equals(Value) bool
	GetAttribute(string) (Value, error)
	HashBytes() ([]byte, error)
	Length() (uint64, error)
	// String returns the formatted representation of the value, like __repr__ in Python.
	String() string
	ToBool() bool
	ToCallable() (Callable, error)
	ToFloat() (float64, error)
	ToInt() (int64, error)
	ToIterator() (Iterator, error)
	ToStr() (string, error)
	ToUint() (uint64, error)
	Type() Type
}
