package types

import (
	"fmt"

	"github.com/chrispyles/slow/internal/errors"
	"github.com/chrispyles/slow/internal/execute"
)

// -------------------------------------------------------------------------------------------------
// Type definition
// -------------------------------------------------------------------------------------------------

type uintType struct{}

func (t *uintType) IsNumeric() bool {
	return true
}

func (t *uintType) New(v execute.Value) (execute.Value, error) {
	vc, err := v.ToUint()
	if err != nil {
		return nil, err
	}
	return NewUint(vc), nil
}

func (t *uintType) String() string {
	return "uint"
}

var UintType = &uintType{}

// -------------------------------------------------------------------------------------------------
// Type implementation
// -------------------------------------------------------------------------------------------------

type Uint struct {
	value uint64
}

func NewUint(v uint64) *Uint {
	return &Uint{v}
}

func (v *Uint) CloneIfPrimitive() execute.Value {
	return NewUint(v.value)
}

func (v *Uint) CompareTo(o execute.Value) (int, bool) {
	switch o.Type() {
	case FloatType:
		return compareNumbers(float64(v.value), must(o.ToFloat())), true
	case IntType:
		return compareNumbers(int64(v.value), must(o.ToInt())), true
	case BoolType:
		fallthrough
	case UintType:
		return compareNumbers(v.value, must(o.ToUint())), true
	}
	return 0, false
}

func (v *Uint) Equals(o execute.Value) bool {
	of, ok := o.(*Uint)
	if !ok {
		return false
	}
	return v.value == of.value
}

func (v *Uint) GetAttribute(a string) (execute.Value, error) {
	return nil, errors.NewAttributeError(v.Type(), a)
}

func (v *Uint) GetIndex(execute.Value) (execute.Value, error) {
	return nil, errors.IndexingNotSupported(v.Type())
}

func (v *Uint) HasAttribute(a string) bool {
	return false
}

func (v *Uint) HashBytes() ([]byte, error) {
	return numToBytes(v.value), nil
}

func (v *Uint) Length() (uint64, error) {
	return 0, errors.NoLengthError(v.Type())
}

func (v *Uint) SetAttribute(a string, _ execute.Value) error {
	return errors.NewAttributeError(v.Type(), a)
}

func (v *Uint) SetIndex(execute.Value, execute.Value) error {
	return errors.IndexingNotSupported(v.Type())
}

func (v *Uint) String() string {
	return fmt.Sprintf("%du", v.value)
}

func (v *Uint) ToBool() bool {
	if v.value == 0 {
		return false
	}
	return true
}

func (v *Uint) ToBytes() ([]byte, error) {
	return numToBytes(v.value), nil
}

func (v *Uint) ToCallable() (execute.Callable, error) {
	return nil, errors.NewTypeError(v.Type(), FuncType)
}

func (v *Uint) ToFloat() (float64, error) {
	return float64(v.value), nil
}

func (v *Uint) ToInt() (int64, error) {
	return int64(v.value), nil
}

func (v *Uint) ToIterator() (execute.Iterator, error) {
	return nil, errors.NewTypeError(v.Type(), IteratorType)
}

func (v *Uint) ToStr() (string, error) {
	return fmt.Sprintf("%d", v.value), nil
}

func (v *Uint) ToUint() (uint64, error) {
	return v.value, nil
}

func (v *Uint) Type() execute.Type {
	return UintType
}
