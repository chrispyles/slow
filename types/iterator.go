package types

import "github.com/chrispyles/slow/execute"

type Iterator struct {
	values []execute.Value
	idx    int
}

func NewIterator(vs []execute.Value) *Iterator {
	return &Iterator{values: vs}
}

func (v *Iterator) HasNext() bool {
	return v.idx < len(v.values)
}

func (v *Iterator) Next() (execute.Value, error) {
	idx := v.idx
	v.idx++
	return v.values[idx], nil
}
