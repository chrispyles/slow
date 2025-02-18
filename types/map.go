package types

import (
	"fmt"
	"hash/maphash"
	"strings"

	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
)

// -------------------------------------------------------------------------------------------------
// Type definition
// -------------------------------------------------------------------------------------------------

type mapType struct{}

func (t *mapType) IsNumeric() bool {
	return false
}

func (t *mapType) New(v execute.Value) (execute.Value, error) {
	if v == nil {
		return NewMap(), nil
	}
	panic("mapType.New() is not supported with a non-nil argument")
}

func (t *mapType) String() string {
	return "map"
}

var MapType = &mapType{}

// -------------------------------------------------------------------------------------------------
// Type implementation
// -------------------------------------------------------------------------------------------------

var mapMethods = map[string]func(*Map) execute.Value{
	"get": func(v *Map) execute.Value {
		return NewGoFunc("map.get", func(vs ...execute.Value) (execute.Value, error) {
			if got, want := len(vs), 2; got > want {
				return nil, errors.CallError("map.get", got, want)
			}
			if got := len(vs); got == 0 {
				return nil, errors.CallError("map.get", got, 1)
			}
			var defaultValue execute.Value
			if len(vs) == 2 {
				defaultValue = vs[1]
			}
			return v.Get(vs[0], defaultValue)
		})
	},
	"set": func(v *Map) execute.Value {
		return NewGoFunc("map.set", func(vs ...execute.Value) (execute.Value, error) {
			if got, want := len(vs), 2; got != want {
				return nil, errors.CallError("map.set", got, want)
			}
			if v.immutable {
				return nil, errors.NewValueError("map is immutable")
			}
			return v.Set(vs[0], vs[1])
		})
	},
	"to_immutable": func(v *Map) execute.Value {
		name := "map.to_immutable"
		return NewGoFunc(name, func(vs ...execute.Value) (execute.Value, error) {
			if got, want := len(vs), 0; got != want {
				return nil, errors.CallError(name, got, want)
			}
			return &Map{v.seed, v.cloneEntries(), v.size, true}, nil
		})
	},
	"to_mutable": func(v *Map) execute.Value {
		name := "map.to_mutable"
		return NewGoFunc(name, func(vs ...execute.Value) (execute.Value, error) {
			if got, want := len(vs), 0; got != want {
				return nil, errors.CallError(name, got, want)
			}
			return &Map{v.seed, v.cloneEntries(), v.size, false}, nil
		})
	},
}

type mapEntries map[uint64][]*mapEntry

type Map struct {
	seed      maphash.Seed
	entries   mapEntries
	size      uint64
	immutable bool
}

type mapEntry struct {
	key   execute.Value
	value execute.Value
}

func NewMap() *Map {
	return &Map{seed: maphash.MakeSeed(), entries: make(map[uint64][]*mapEntry)}
}

func (v *Map) cloneEntries() mapEntries {
	entries := make(map[uint64][]*mapEntry)
	for k, vs := range v.entries {
		for _, v := range vs {
			vCopy := *v
			entries[k] = append(entries[k], &vCopy)
		}
	}
	return entries
}

func (v *Map) hash(val execute.Value) (uint64, error) {
	hb, err := val.HashBytes()
	if err != nil {
		return 0, err
	}
	return maphash.Bytes(v.seed, hb), nil
}

func (v *Map) Get(key execute.Value, defaultValue execute.Value) (execute.Value, error) {
	h, err := v.hash(key)
	if err != nil {
		return nil, err
	}
	if es, ok := v.entries[h]; ok {
		for _, e := range es {
			if key.Equals(e.key) {
				return e.value, nil
			}
		}
	}
	if defaultValue != nil {
		return defaultValue, nil
	}
	return nil, errors.NewKeyError(key.String())
}

func (v *Map) Has(key execute.Value) (bool, error) {
	h, err := v.hash(key)
	if err != nil {
		return false, err
	}
	if es, ok := v.entries[h]; ok {
		for _, e := range es {
			if key.Equals(e.key) {
				return true, nil
			}
		}
	}
	return false, nil
}

func (v *Map) Set(key execute.Value, value execute.Value) (execute.Value, error) {
	if v.immutable {
		return nil, errors.NewValueError("map is immutable")
	}
	h, err := v.hash(key)
	if err != nil {
		return nil, err
	}
	if _, ok := v.entries[h]; !ok {
		v.entries[h] = []*mapEntry{}
	}
	var found bool
	for _, e := range v.entries[h] {
		if e.key.Equals(key) {
			e.value = value
			found = true
		}
	}
	if !found {
		v.entries[h] = append(v.entries[h], &mapEntry{key, value})
		v.size++
	}
	return NewBool(found), nil
}

func (v *Map) CloneIfPrimitive() execute.Value {
	return v
}

func (v *Map) CompareTo(o execute.Value) (int, bool) {
	return 0, false
}

func (v *Map) Equals(o execute.Value) bool {
	if l2, ok := o.(*Map); ok {
		// TODO: is it possible to get two *Map pointers that represent the same map in the Slow environment?
		return v == l2
	}
	return false
}

func (v *Map) GetAttribute(a string) (execute.Value, error) {
	if methodFactory, ok := mapMethods[a]; ok {
		return methodFactory(v), nil
	}
	return nil, errors.NewAttributeError(v.Type(), a)
}

func (v *Map) GetIndex(i execute.Value) (execute.Value, error) {
	return v.Get(i, nil)
}

func (v *Map) HasAttribute(a string) bool {
	_, ok := mapMethods[a]
	return ok
}

func (v *Map) HashBytes() ([]byte, error) {
	return nil, errors.UnhashableTypeError(v.Type())
}

func (v *Map) Length() (uint64, error) {
	return v.size, nil
}

func (v *Map) SetAttribute(a string, _ execute.Value) error {
	if v.immutable {
		return errors.NewValueError("map is immutable")
	}
	if v.HasAttribute(a) {
		return errors.AssignmentError(v.Type(), a)
	}
	return errors.NewAttributeError(v.Type(), a)
}

func (v *Map) SetIndex(i execute.Value, val execute.Value) error {
	if v.immutable {
		return errors.NewValueError("map is immutable")
	}
	_, err := v.Set(i, val)
	return err
}

func (v *Map) String() string {
	var items []string
	for _, es := range v.entries {
		for _, e := range es {
			items = append(items, fmt.Sprintf("%s: %s", e.key.String(), e.value.String()))
		}
	}
	return fmt.Sprintf("{%s}", strings.Join(items, ", "))
}

func (v *Map) ToBool() bool {
	return true
}

func (v *Map) ToBytes() ([]byte, error) {
	return nil, errors.NewTypeError(v.Type(), BytesType)
}

func (v *Map) ToCallable() (execute.Callable, error) {
	return nil, errors.NewTypeError(v.Type(), FuncType)
}

func (v *Map) ToFloat() (float64, error) {
	return 0, errors.NewTypeError(v.Type(), FloatType)
}

func (v *Map) ToInt() (int64, error) {
	return 0, errors.NewTypeError(v.Type(), IntType)
}

func (v *Map) ToIterator() (execute.Iterator, error) {
	return newMapIterator(v), nil
}

func (v *Map) ToStr() (string, error) {
	return v.String(), nil
}

func (v *Map) ToUint() (uint64, error) {
	return 0, errors.NewTypeError(v.Type(), UintType)
}

func (v *Map) Type() execute.Type {
	return MapType
}

type mapIterator struct {
	map_ *Map
	keys []execute.Value
	idx  int
}

func newMapIterator(m *Map) *mapIterator {
	var keys []execute.Value
	for _, es := range m.entries {
		for _, e := range es {
			keys = append(keys, e.key)
		}
	}
	return &mapIterator{m, keys, 0}
}

func (i *mapIterator) HasNext() bool {
	return i.idx < len(i.keys)
}

func (i *mapIterator) Next() (execute.Value, error) {
	if i.idx >= len(i.keys) {
		return nil, errors.NewIndexError(fmt.Sprintf("%d", i.idx))
	}
	key := i.keys[i.idx]
	i.idx++
	return key, nil
}
