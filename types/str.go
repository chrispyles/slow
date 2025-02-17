package types

import (
	"fmt"
	"strconv"

	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
)

// -------------------------------------------------------------------------------------------------
// Type definition
// -------------------------------------------------------------------------------------------------

type strType struct{}

func (t *strType) IsNumeric() bool {
	return false
}

func (t *strType) New(v execute.Value) (execute.Value, error) {
	vc, err := v.ToStr()
	if err != nil {
		return nil, err
	}
	return NewStr(vc), nil
}

func (t *strType) String() string {
	return "str"
}

var StrType = &strType{}

// -------------------------------------------------------------------------------------------------
// Type implementation
// -------------------------------------------------------------------------------------------------

type Str struct {
	value string
}

func NewStr(v string) *Str {
	return &Str{value: v}
}

func (v *Str) Value() string {
	return v.value
}

func (v *Str) CloneIfPrimitive() execute.Value {
	return NewStr(v.value)
}

func (v *Str) CompareTo(o execute.Value) (int, bool) {
	if o.Type() == StrType {
		os := must(o.ToStr())
		if v.value == os {
			return 0, true
		} else if v.value < os {
			return -1, true
		}
		return 1, true
	}
	return 0, false
}

func (v *Str) Equals(o execute.Value) bool {
	os, ok := o.(*Str)
	if !ok {
		return false
	}
	return v.value == os.value
}

func (v *Str) GetAttribute(a string) (execute.Value, error) {
	return nil, errors.NewAttributeError(v.Type(), a)
}

func (v *Str) GetIndex(i execute.Value) (execute.Value, error) {
	idx, err := numericIndex(i, v.Type())
	if err != nil {
		return nil, err
	}
	idx, ok := normalizeIndex(idx, len(v.value))
	if !ok {
		return nil, errors.NewIndexError(fmt.Sprintf("%d", idx))
	}
	return NewStr(string(v.value[idx])), nil
}

func (v *Str) HasAttribute(a string) bool {
	return false
}

func (v *Str) HashBytes() ([]byte, error) {
	return []byte(v.value), nil
}

func (v *Str) Length() (uint64, error) {
	return uint64(len(v.value)), nil
}

func (v *Str) SetAttribute(a string, _ execute.Value) error {
	return errors.NewAttributeError(v.Type(), a)
}

func (v *Str) SetIndex(execute.Value, execute.Value) error {
	return errors.SetIndexNotSupported(v.Type())
}

func (v *Str) String() string {
	return fmt.Sprintf("%q", v.value)
}

func (v *Str) ToBool() bool {
	return v.value != ""
}

func (v *Str) ToBytes() ([]byte, error) {
	return []byte(v.value), nil
}

func (v *Str) ToCallable() (execute.Callable, error) {
	return nil, errors.NewTypeError(v.Type(), FuncType)
}

func (v *Str) ToFloat() (float64, error) {
	vf, err := strconv.ParseFloat(v.value, 64)
	return vf, errors.WrapValueError(v.value, FloatType, err)
}

func (v *Str) ToInt() (int64, error) {
	vi, err := strconv.ParseInt(v.value, 10, 64)
	return vi, errors.WrapValueError(v.value, IntType, err)
}

func (v *Str) ToIterator() (execute.Iterator, error) {
	return &stringIterator{s: v}, nil
}

func (v *Str) ToStr() (string, error) {
	return v.value, nil
}

func (v *Str) ToUint() (uint64, error) {
	vu, err := strconv.ParseUint(v.value, 10, 64)
	return vu, errors.WrapValueError(v.value, UintType, err)
}

func (v *Str) Type() execute.Type {
	return StrType
}

type stringIterator struct {
	idx int
	s   *Str
}

func (si *stringIterator) HasNext() bool {
	return si.idx < len(si.s.value)
}

func (si *stringIterator) Next() (execute.Value, error) {
	c := string(si.s.value[si.idx])
	si.idx++
	return NewStr(c), nil
}
