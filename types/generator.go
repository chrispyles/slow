package types

import (
	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
)

// -------------------------------------------------------------------------------------------------
// Type definition
// -------------------------------------------------------------------------------------------------

type generatorType struct{}

func (t *generatorType) IsNumeric() bool {
	return false
}

func (t *generatorType) New(v execute.Value) (execute.Value, error) {
	panic("generatorType.New() is not supported")
}

func (t *generatorType) String() string {
	return "generator"
}

var GeneratorType = &generatorType{}

// -------------------------------------------------------------------------------------------------
// Type implementation
// -------------------------------------------------------------------------------------------------

type generatorImpl interface {
	execute.Iterator
	WithContainerLen(uint64) *Generator
}

// Generator is a type that wraps an execute.Iterator and satisfies both the execute.Iterator and
// execute.Value interfaces.
type Generator struct {
	impl generatorImpl
}

func NewGenerator(gi generatorImpl) *Generator {
	return &Generator{gi}
}

// execute.Iterator methos

func (v *Generator) HasNext() bool {
	return v.impl.HasNext()
}

func (v *Generator) Next() (execute.Value, error) {
	return v.impl.Next()
}

func (v *Generator) WithContainerLen(l uint64) *Generator {
	return v.impl.WithContainerLen(l)
}

// execute.Value methods

func (v *Generator) CloneIfPrimitive() execute.Value {
	return v
}

func (v *Generator) CompareTo(o execute.Value) (int, bool) {
	return 0, false
}

func (v *Generator) Equals(o execute.Value) bool {
	og, ok := o.(*Generator)
	return ok && v == og
}

func (v *Generator) GetAttribute(a string) (execute.Value, error) {
	return nil, errors.NewAttributeError(v.Type(), a)
}

func (v *Generator) GetIndex(execute.Value) (execute.Value, error) {
	return nil, errors.IndexingNotSupported(v.Type())
}

func (v *Generator) HasAttribute(a string) bool {
	return false
}

func (v *Generator) HashBytes() ([]byte, error) {
	return nil, errors.UnhashableTypeError(v.Type())
}

func (v *Generator) Length() (uint64, error) {
	return 0, errors.NoLengthError(GeneratorType)
}

func (v *Generator) SetAttribute(a string, _ execute.Value) error {
	return errors.NewAttributeError(v.Type(), a)
}

func (v *Generator) SetIndex(execute.Value, execute.Value) error {
	return errors.IndexingNotSupported(v.Type())
}

func (v *Generator) String() string {
	return "<generator>"
}

func (v *Generator) ToBool() bool {
	return true
}

func (v *Generator) ToBytes() ([]byte, error) {
	return nil, errors.NewTypeError(v.Type(), BytesType)
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
