package types

import (
	"fmt"
	"strings"

	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
)

// -------------------------------------------------------------------------------------------------
// Type definition
// -------------------------------------------------------------------------------------------------

type listType struct{}

func (t *listType) IsNumeric() bool {
	return false
}

func (t *listType) New(v execute.Value) (execute.Value, error) {
	if v == nil {
		return NewList(nil), nil
	}
	panic("listType.New() is not supported with a non-nil argument")
}

func (t *listType) String() string {
	return "list"
}

var ListType = &listType{}

// -------------------------------------------------------------------------------------------------
// Type implementation
// -------------------------------------------------------------------------------------------------

var listMethods = map[string]func(*List) execute.Value{
	"append": func(v *List) execute.Value {
		name := "list.append"
		return NewGoFunc(name, func(vs ...execute.Value) (execute.Value, error) {
			if got, want := len(vs), 1; got != want {
				return nil, errors.CallError(name, got, want)
			}
			if v.immutable {
				return nil, errors.NewValueError("list is immutable")
			}
			v.values = append(v.values, vs[0])
			return Null, nil
		})
	},
	"to_immutable": func(v *List) execute.Value {
		name := "list.to_immutable"
		return NewGoFunc(name, func(vs ...execute.Value) (execute.Value, error) {
			if got, want := len(vs), 0; got != want {
				return nil, errors.CallError(name, got, want)
			}
			return &List{v.values, true}, nil
		})
	},
	"to_mutable": func(v *List) execute.Value {
		name := "list.to_mutable"
		return NewGoFunc(name, func(vs ...execute.Value) (execute.Value, error) {
			if got, want := len(vs), 0; got != want {
				return nil, errors.CallError(name, got, want)
			}
			return &List{v.values, false}, nil
		})
	},
}

type List struct {
	values    []execute.Value
	immutable bool
}

func NewList(vs []execute.Value) *List {
	return &List{values: vs}
}

func (v *List) CloneIfPrimitive() execute.Value {
	return v
}

func (v *List) CompareTo(o execute.Value) (int, bool) {
	return 0, false
}

func (v *List) Equals(o execute.Value) bool {
	if l2, ok := o.(*List); ok {
		return v == l2
	}
	return false
}

func (v *List) GetAttribute(a string) (execute.Value, error) {
	if methodFactory, ok := listMethods[a]; ok {
		return methodFactory(v), nil
	}
	return nil, errors.NewAttributeError(v.Type(), a)
}

func (v *List) GetIndex(i execute.Value) (execute.Value, error) {
	if g, ok := i.(*Generator); ok {
		var indices []int
		g = g.WithContainerLen(uint64(len(v.values)))
		for g.HasNext() {
			iv, err := g.Next()
			if err != nil {
				return nil, err
			}
			i, err := numericIndex(iv, v.Type())
			if err != nil {
				return nil, err
			}
			indices = append(indices, i)
		}
		var sublist []execute.Value
		for _, i := range indices {
			v, err := v.getIndex(i)
			if err != nil {
				return nil, err
			}
			sublist = append(sublist, v)
		}
		return NewList(sublist), nil
	}
	idx, err := numericIndex(i, v.Type())
	if err != nil {
		return nil, err
	}
	return v.getIndex(idx)
}

func (v *List) getIndex(idx int) (execute.Value, error) {
	idx, ok := normalizeIndex(idx, len(v.values))
	if !ok {
		return nil, errors.NewIndexError(fmt.Sprintf("%d", idx))
	}
	return v.values[idx], nil
}

func (v *List) HasAttribute(a string) bool {
	_, ok := listMethods[a]
	return ok
}

func (v *List) HashBytes() ([]byte, error) {
	return nil, errors.UnhashableTypeError(v.Type())
}

func (v *List) Length() (uint64, error) {
	return uint64(len(v.values)), nil
}

func (v *List) SetAttribute(a string, _ execute.Value) error {
	if v.immutable {
		return errors.NewValueError("list is immutable")
	}
	if v.HasAttribute(a) {
		return errors.AssignmentError(v.Type(), a)
	}
	return errors.NewAttributeError(v.Type(), a)
}

func (v *List) SetIndex(i execute.Value, val execute.Value) error {
	if v.immutable {
		return errors.NewValueError("list is immutable")
	}
	idx, err := numericIndex(i, v.Type())
	if err != nil {
		return err
	}
	idx, ok := normalizeIndex(idx, len(v.values))
	if !ok {
		return errors.NewIndexError(fmt.Sprintf("%d", idx))
	}
	v.values[idx] = val
	return nil
}

func (v *List) String() string {
	items := make([]string, len(v.values))
	for i, v := range v.values {
		items[i] = v.String()
	}
	return fmt.Sprintf("[%s]", strings.Join(items, ", "))
}

func (v *List) ToBool() bool {
	return true
}

func (v *List) ToBytes() ([]byte, error) {
	return nil, errors.NewTypeError(v.Type(), BytesType)
}

func (v *List) ToCallable() (execute.Callable, error) {
	return nil, errors.NewTypeError(v.Type(), FuncType)
}

func (v *List) ToFloat() (float64, error) {
	return 0, errors.NewTypeError(v.Type(), FloatType)
}

func (v *List) ToInt() (int64, error) {
	return 0, errors.NewTypeError(v.Type(), IntType)
}

func (v *List) ToIterator() (execute.Iterator, error) {
	return &listIterator{v, 0}, nil
}

func (v *List) ToStr() (string, error) {
	return v.String(), nil
}

func (v *List) ToUint() (uint64, error) {
	return 0, errors.NewTypeError(v.Type(), UintType)
}

func (v *List) Type() execute.Type {
	return ListType
}

type listIterator struct {
	list *List
	idx  int
}

func (i *listIterator) HasNext() bool {
	return i.idx <= len(i.list.values)-1
}

func (i *listIterator) Next() (execute.Value, error) {
	if i.idx >= len(i.list.values) {
		return nil, errors.NewIndexError(fmt.Sprintf("%d", i.idx))
	}
	v := i.list.values[i.idx]
	i.idx++
	return v, nil
}
