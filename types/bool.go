package types

import (
	"fmt"

	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
)

type Bool struct {
	value bool
}

func NewBool(v bool) *Bool {
	return &Bool{v}
}

func (v *Bool) CloneIfPrimitive() execute.Value {
	return NewBool(v.value)
}

func (v *Bool) CompareTo(o execute.Value) (int, bool) {
	switch o.Type() {
	case BoolType:
		ob := o.(*Bool)
		if v.value == ob.value {
			return 0, true
		} else if ob.value {
			return -1, true
		}
		return 1, true
	case FloatType:
		ou := o.(*Float)
		return compareNumbers(must(v.ToFloat()), ou.value), true
	case IntType:
		oi := o.(*Int)
		return compareNumbers(must(v.ToInt()), oi.value), true
	case UintType:
		ou := o.(*Uint)
		return compareNumbers(must(v.ToUint()), ou.value), true
	}
	return 0, false
}

func (v *Bool) Equals(o execute.Value) bool {
	ob, ok := o.(*Bool)
	if !ok {
		return false
	}
	return v.value == ob.value
}

func (v *Bool) GetAttribute(a string) (execute.Value, error) {
	return nil, errors.NewAttributeError(v.Type(), a)
}

func (v *Bool) GetIndex(execute.Value) (execute.Value, error) {
	return nil, errors.IndexingNotSupported(v.Type())
}

func (v *Bool) HasAttribute(a string) bool {
	return false
}

func (v *Bool) HashBytes() ([]byte, error) {
	if v.value {
		return []byte{0x01}, nil
	}
	return []byte{0x00}, nil
}

func (v *Bool) Length() (uint64, error) {
	return 0, errors.NoLengthError(v.Type())
}

func (v *Bool) SetAttribute(a string, _ execute.Value) error {
	return errors.NewAttributeError(v.Type(), a)
}

func (v *Bool) SetIndex(execute.Value, execute.Value) error {
	return errors.IndexingNotSupported(v.Type())
}

func (v *Bool) String() string {
	return fmt.Sprint(v.value)
}

func (v *Bool) ToBool() bool {
	return v.value
}

func (v *Bool) ToCallable() (execute.Callable, error) {
	return nil, errors.NewTypeError(v.Type(), FuncType)
}

func (v *Bool) ToFloat() (float64, error) {
	if v.value {
		return 1, nil
	} else {
		return 0, nil
	}
}

func (v *Bool) ToInt() (int64, error) {
	if v.value {
		return 1, nil
	} else {
		return 0, nil
	}
}

func (v *Bool) ToIterator() (execute.Iterator, error) {
	return nil, errors.NewTypeError(v.Type(), IteratorType)
}

func (v *Bool) ToStr() (string, error) {
	if v.value {
		return "true", nil
	} else {
		return "false", nil
	}
}

func (v *Bool) ToUint() (uint64, error) {
	if v.value {
		return 1, nil
	} else {
		return 0, nil
	}
}

func (v *Bool) Type() execute.Type {
	return BoolType
}
