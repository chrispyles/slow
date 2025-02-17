package types

import (
	"fmt"

	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
)

// -------------------------------------------------------------------------------------------------
// Type definition
// -------------------------------------------------------------------------------------------------

type intType struct{}

func (t *intType) IsNumeric() bool {
	return true
}

func (t *intType) New(v execute.Value) (execute.Value, error) {
	vc, err := v.ToInt()
	if err != nil {
		return nil, err
	}
	return NewInt(vc), nil
}

func (t *intType) String() string {
	return "int"
}

var IntType = &intType{}

// -------------------------------------------------------------------------------------------------
// Type implementation
// -------------------------------------------------------------------------------------------------

type Int struct {
	value int64
}

func NewInt(v int64) *Int {
	return &Int{v}
}

func (v *Int) CloneIfPrimitive() execute.Value {
	return NewInt(v.value)
}

func (v *Int) CompareTo(o execute.Value) (int, bool) {
	switch o.Type() {
	case FloatType:
		return compareNumbers(float64(v.value), must(o.ToFloat())), true
	case BoolType:
		fallthrough
	case IntType:
		fallthrough
	case UintType:
		return compareNumbers(v.value, must(o.ToInt())), true
	}
	return 0, false
}

func (v *Int) Equals(o execute.Value) bool {
	of, ok := o.(*Int)
	if !ok {
		return false
	}
	return v.value == of.value
}

func (v *Int) GetAttribute(a string) (execute.Value, error) {
	return nil, errors.NewAttributeError(v.Type(), a)
}

func (v *Int) GetIndex(execute.Value) (execute.Value, error) {
	return nil, errors.IndexingNotSupported(v.Type())
}

func (v *Int) HasAttribute(a string) bool {
	return false
}

func (v *Int) HashBytes() ([]byte, error) {
	return numToBytes(v.value), nil
}

func (v *Int) Length() (uint64, error) {
	return 0, errors.NoLengthError(v.Type())
}

func (v *Int) SetAttribute(a string, _ execute.Value) error {
	return errors.NewAttributeError(v.Type(), a)
}

func (v *Int) SetIndex(execute.Value, execute.Value) error {
	return errors.IndexingNotSupported(v.Type())
}

func (v *Int) String() string {
	return fmt.Sprint(v.value)
}

func (v *Int) ToBool() bool {
	if v.value == 0 {
		return false
	}
	return true
}

func (v *Int) ToBytes() ([]byte, error) {
	return numToBytes(v.value), nil
}

func (v *Int) ToCallable() (execute.Callable, error) {
	return nil, errors.NewTypeError(v.Type(), FuncType)
}

func (v *Int) ToFloat() (float64, error) {
	return float64(v.value), nil
}

func (v *Int) ToInt() (int64, error) {
	return v.value, nil
}

func (v *Int) ToIterator() (execute.Iterator, error) {
	return nil, errors.NewTypeError(v.Type(), IteratorType)
}

func (v *Int) ToStr() (string, error) {
	return fmt.Sprintf("%d", v.value), nil
}

func (v *Int) ToUint() (uint64, error) {
	return uint64(v.value), nil
}

func (v *Int) Type() execute.Type {
	return IntType
}
