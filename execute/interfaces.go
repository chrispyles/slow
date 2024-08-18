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

type Value interface {
	// CloneIfPrimitive returns a copy of this value object if it is a primitive (i.e. pass-by-value).
	// Values that are pass-by-reference should return a reference to the same instance.
	CloneIfPrimitive() Value
	// N.B. Equals() should return false if the value is not of the same type (i.e.
	// int(1) != float(1) for the purposes of this method).
	Equals(Value) bool
	GetAttribute(string) (Value, error)
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
