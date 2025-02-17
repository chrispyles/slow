package types

import (
	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
)

// -------------------------------------------------------------------------------------------------
// Type definition
// -------------------------------------------------------------------------------------------------

type nullType struct{}

func (t *nullType) IsNumeric() bool {
	return false
}

func (t *nullType) New(v execute.Value) (execute.Value, error) {
	panic("nullType.New() is not supported")
}

func (t *nullType) String() string {
	return "null"
}

var NullType = &nullType{}

// -------------------------------------------------------------------------------------------------
// Type implementation
// -------------------------------------------------------------------------------------------------

// null is a struct meant to be used as a singleton to represent a null value.
type null struct{}

var Null = &null{}

func (v *null) CloneIfPrimitive() execute.Value {
	return Null
}

func (v *null) CompareTo(o execute.Value) (int, bool) {
	return 0, false
}

func (v *null) Equals(o execute.Value) bool {
	return o == Null
}

func (v *null) GetAttribute(a string) (execute.Value, error) {
	return nil, errors.NewAttributeError(v.Type(), a)
}

func (v *null) GetIndex(execute.Value) (execute.Value, error) {
	return nil, errors.IndexingNotSupported(v.Type())
}

func (v *null) HasAttribute(a string) bool {
	return false
}

func (v *null) HashBytes() ([]byte, error) {
	return []byte{0x00}, nil
}

func (v *null) Length() (uint64, error) {
	return 0, errors.NoLengthError(v.Type())
}

func (v *null) SetAttribute(a string, _ execute.Value) error {
	return errors.NewAttributeError(v.Type(), a)
}

func (v *null) SetIndex(execute.Value, execute.Value) error {
	return errors.IndexingNotSupported(v.Type())
}

func (v *null) String() string {
	return "null"
}

func (v *null) ToBool() bool {
	return false
}

func (v *null) ToBytes() ([]byte, error) {
	return nil, errors.NewTypeError(v.Type(), BytesType)
}

func (v *null) ToCallable() (execute.Callable, error) {
	return nil, errors.NewTypeError(v.Type(), FuncType)
}

func (v *null) ToFloat() (float64, error) {
	return 0, errors.NewTypeError(v.Type(), FloatType)
}

func (v *null) ToInt() (int64, error) {
	return 0, errors.NewTypeError(v.Type(), IntType)
}

func (v *null) ToIterator() (execute.Iterator, error) {
	return nil, errors.NewTypeError(v.Type(), IteratorType)
}

func (v *null) ToStr() (string, error) {
	return v.String(), nil
}

func (v *null) ToUint() (uint64, error) {
	return 0, errors.NewTypeError(v.Type(), UintType)
}

func (v *null) Type() execute.Type {
	return NullType
}
