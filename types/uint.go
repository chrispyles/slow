package types

import (
	"fmt"

	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
)

type Uint struct {
	value uint64
}

func NewUint(v uint64) *Uint {
	return &Uint{v}
}

func (v *Uint) Equals(o execute.Value) bool {
	of, ok := o.(*Uint)
	if !ok {
		return false
	}
	return v.value == of.value
}

func (v *Uint) String() string {
	return fmt.Sprintf("%du", v.value)
}

func (v *Uint) GetAttribute(a string) (execute.Value, error) {
	return nil, errors.NewAttributeError(v.Type(), a)
}

func (v *Uint) Length() (uint64, error) {
	return 0, errors.NoLengthError(v.Type())
}

func (v *Uint) ToBool() bool {
	if v.value == 0 {
		return false
	}
	return true
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
