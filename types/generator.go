package types

import (
	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
)

// Generator is a type that wraps an execute.Iterator and satisfies both the execute.Iterator and
// execute.Value interfaces.
type Generator struct {
	impl execute.Iterator
}

func NewGenerator(gi execute.Iterator) *Generator {
	return &Generator{gi}
}

// execute.Iterator methos

func (v *Generator) HasNext() bool {
	return v.impl.HasNext()
}

func (v *Generator) Next() (execute.Value, error) {
	return v.impl.Next()
}

// execute.Value methods

func (v *Generator) CloneIfPrimitive() execute.Value {
	return v
}

func (v *Generator) Equals(o execute.Value) bool {
	og, ok := o.(*Generator)
	return ok && v == og
}

func (v *Generator) GetAttribute(string) (execute.Value, error) {
	// TODO
	return nil, nil
}

func (v *Generator) Length() (uint64, error) {
	return 0, errors.NoLengthError(GeneratorType)
}

func (v *Generator) String() string {
	return "<generator>"
}

func (v *Generator) ToBool() bool {
	return true
}

func (v *Generator) ToCallable() (execute.Callable, error) {
	return nil, errors.NewTypeError(v.Type(), FuncType)
}

func (v *Generator) ToFloat() (float64, error) {
	return 0, errors.NewTypeError(v.Type(), FloatType)
}

func (v *Generator) ToInt() (int64, error) {
	return 0, errors.NewTypeError(v.Type(), IntType)
}

func (v *Generator) ToIterator() (execute.Iterator, error) {
	return v, nil
}

func (v *Generator) ToStr() (string, error) {
	return v.String(), nil
}

func (v *Generator) ToUint() (uint64, error) {
	return 0, errors.NewTypeError(v.Type(), UintType)
}

func (v *Generator) Type() execute.Type {
	return GeneratorType
}
