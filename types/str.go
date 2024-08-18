package types

import (
	"fmt"
	"strconv"

	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
)

type Str struct {
	value string
}

func NewStr(v string) *Str {
	return &Str{value: v}
}

func (v *Str) CloneIfPrimitive() execute.Value {
	return NewStr(v.value)
}

func (v *Str) Equals(o execute.Value) bool {
	of, ok := o.(*Str)
	if !ok {
		return false
	}
	return v.value == of.value
}

func (v *Str) GetAttribute(a string) (execute.Value, error) {
	return nil, errors.NewAttributeError(v.Type(), a)
}

func (v *Str) Length() (uint64, error) {
	return uint64(len(v.value)), nil
}

func (v *Str) String() string {
	return fmt.Sprintf("%q", v.value)
}

func (v *Str) ToBool() bool {
	return v.value != ""
}

func (v *Str) ToCallable() (execute.Callable, error) {
	return nil, errors.NewTypeError(v.Type(), FuncType)
}

func (v *Str) ToFloat() (float64, error) {
	return strconv.ParseFloat(v.value, 64) // TODO: wrap error
}

func (v *Str) ToInt() (int64, error) {
	return strconv.ParseInt(v.value, 10, 64) // TODO: wrap error
}

func (v *Str) ToIterator() (execute.Iterator, error) {
	// TODO: string iterator
	return nil, errors.NewTypeError(v.Type(), IteratorType)
}

func (v *Str) ToStr() (string, error) {
	return v.value, nil
}

func (v *Str) ToUint() (uint64, error) {
	return strconv.ParseUint(v.value, 10, 64) // TODO: wrap error
}

func (v *Str) Type() execute.Type {
	return StrType
}
