package types

import (
	"fmt"
	"math"

	"github.com/chrispyles/slow/internal/errors"
	"github.com/chrispyles/slow/internal/execute"
)

// -------------------------------------------------------------------------------------------------
// Type definition
// -------------------------------------------------------------------------------------------------

type floatType struct{}

func (t *floatType) IsNumeric() bool {
	return true
}

func (t *floatType) New(v execute.Value) (execute.Value, error) {
	vc, err := v.ToFloat()
	if err != nil {
		return nil, err
	}
	return NewFloat(vc), nil
}

func (t *floatType) String() string {
	return "float"
}

var FloatType = &floatType{}

// -------------------------------------------------------------------------------------------------
// Type implementation
// -------------------------------------------------------------------------------------------------

type Float struct {
	value float64
}

func NewFloat(v float64) *Float {
	return &Float{v}
}

func (v *Float) CloneIfPrimitive() execute.Value {
	return NewFloat(v.value)
}

func (v *Float) CompareTo(o execute.Value) (int, bool) {
	switch o.Type() {
	case BoolType:
		fallthrough
	case FloatType:
		fallthrough
	case IntType:
		fallthrough
	case UintType:
		return compareNumbers(v.value, must(o.ToFloat())), true
	}
	return 0, false
}

func (v *Float) Equals(o execute.Value) bool {
	of, ok := o.(*Float)
	if !ok {
		return false
	}
	return v.value == of.value
}

func (v *Float) GetAttribute(a string) (execute.Value, error) {
	return nil, errors.NewAttributeError(v.Type(), a)
}

func (v *Float) GetIndex(execute.Value) (execute.Value, error) {
	return nil, errors.IndexingNotSupported(v.Type())
}

func (v *Float) SetIndex(execute.Value, execute.Value) error {
	return errors.IndexingNotSupported(v.Type())
}

func (v *Float) HasAttribute(a string) bool {
	return false
}

func (v *Float) HashBytes() ([]byte, error) {
	return numToBytes(v.value), nil
}

func (v *Float) Length() (uint64, error) {
	return 0, errors.NoLengthError(v.Type())
}

// HasRemainder returns true if the decimal part of this float is nonzero (e.g. this is false for
// 1.0 but true for 1.1).
func (v *Float) HasRemainder() bool {
	return v.value != math.Trunc(v.value)
}

func (v *Float) SetAttribute(a string, _ execute.Value) error {
	return errors.NewAttributeError(v.Type(), a)
}

func (v *Float) String() string {
	out := fmt.Sprint(v.value)
	if v.value == math.Trunc(v.value) {
		out += ".0"
	}
	return out
}

func (v *Float) ToBool() bool {
	if v.value == 0 {
		return false
	}
	return true
}

func (v *Float) ToBytes() ([]byte, error) {
	return numToBytes(v.value), nil
}

func (v *Float) ToCallable() (execute.Callable, error) {
	return nil, errors.NewTypeError(v.Type(), FuncType)
}

func (v *Float) ToFloat() (float64, error) {
	return v.value, nil
}

func (v *Float) ToInt() (int64, error) {
	return int64(v.value), nil
}

func (v *Float) ToIterator() (execute.Iterator, error) {
	return nil, errors.NewTypeError(v.Type(), IteratorType)
}

func (v *Float) ToStr() (string, error) {
	return fmt.Sprintf("%f", v.value), nil
}

func (v *Float) ToUint() (uint64, error) {
	return uint64(v.value), nil
}

func (v *Float) Type() execute.Type {
	return FloatType
}
