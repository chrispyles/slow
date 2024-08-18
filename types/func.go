package types

import (
	"fmt"

	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
)

// ReturnError is an error that indicates that a return statement has been executed.
type ReturnError struct {
	Value execute.Value
}

func (*ReturnError) Error() string { return "" }

// FuncImpl is a function whose logic is implemented in Go, for builtins.
type FuncImpl func(...execute.Value) (execute.Value, error)

type Func struct {
	name string
	args []string
	body execute.Block
	impl FuncImpl
}

// NewFunc creates a new types.Func for a user-defined function.
func NewFunc(name string, args []string, body execute.Block) *Func {
	return &Func{name: name, args: args, body: body}
}

// NewGoFunc creates a new types.Func for a builtin funtion, whose logic is implemented in Go.
func NewGoFunc(name string, impl FuncImpl) *Func {
	return &Func{name: name, impl: impl}
}

func (v *Func) Call(env *execute.Environment, args ...execute.Value) (execute.Value, error) {
	if v.impl != nil {
		return v.impl(args...)
	}
	if got, want := len(args), len(v.args); got != want {
		return nil, errors.CallError(v.name, got, want)
	}
	for i := range v.args {
		name, val := v.args[i], args[i]
		if err := env.Declare(name); err != nil {
			return nil, err
		}
		if _, err := env.Set(name, val); err != nil {
			return nil, err
		}
	}
	for _, expr := range v.body {
		_, err := expr.Execute(env)
		if err != nil {
			if re, ok := err.(*ReturnError); ok {
				return re.Value, nil
			}
			return nil, err
		}
	}
	// A function with no return statement returns null.
	return Null, nil
}

func (v *Func) CloneIfPrimitive() execute.Value {
	return v
}

func (v *Func) CompareTo(o execute.Value) (int, bool) {
	return 0, false
}

func (v *Func) Equals(o execute.Value) bool {
	oc, err := o.ToCallable()
	if err != nil {
		return false
	}
	return v == oc
}

func (v *Func) Length() (uint64, error) {
	return 0, errors.NoLengthError(v.Type())
}

func (v *Func) String() string {
	return fmt.Sprintf("<function %s>", v.name)
}

func (v *Func) GetAttribute(a string) (execute.Value, error) {
	return nil, errors.NewAttributeError(v.Type(), a)
}

func (v *Func) ToBool() bool {
	return true
}

func (v *Func) ToCallable() (execute.Callable, error) {
	return v, nil
}

func (v *Func) ToFloat() (float64, error) {
	return 0, errors.NewTypeError(v.Type(), FloatType)
}

func (v *Func) ToInt() (int64, error) {
	return 0, errors.NewTypeError(v.Type(), IntType)
}

func (v *Func) ToIterator() (execute.Iterator, error) {
	return nil, errors.NewTypeError(v.Type(), IteratorType)
}

func (v *Func) ToStr() (string, error) {
	return "", errors.NewTypeError(v.Type(), StrType)
}

func (v *Func) ToUint() (uint64, error) {
	return 0, errors.NewTypeError(v.Type(), UintType)
}

func (v *Func) Type() execute.Type {
	return FuncType
}
