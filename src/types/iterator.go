package types

import "github.com/chrispyles/slow/src/execute"

// -------------------------------------------------------------------------------------------------
// Type definition
// -------------------------------------------------------------------------------------------------

type iteratorType struct{}

func (t *iteratorType) IsNumeric() bool {
	return false
}

func (t *iteratorType) New(v execute.Value) (execute.Value, error) {
	panic("iteratorType.New() is not supported")
}

func (t *iteratorType) String() string {
	return "iterator"
}

var IteratorType = &iteratorType{}

// -------------------------------------------------------------------------------------------------
// Type implementation
// -------------------------------------------------------------------------------------------------

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
