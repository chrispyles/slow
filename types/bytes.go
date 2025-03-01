package types

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"unicode/utf8"

	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
)

// -------------------------------------------------------------------------------------------------
// Type definition
// -------------------------------------------------------------------------------------------------

type bytesType struct{}

func (t *bytesType) IsNumeric() bool {
	return true
}

func (t *bytesType) New(v execute.Value) (execute.Value, error) {
	vc, err := v.ToBytes()
	if err != nil {
		return nil, err
	}
	return NewBytes(vc), nil
}

func (t *bytesType) String() string {
	return "bytes"
}

var BytesType = &bytesType{}

// -------------------------------------------------------------------------------------------------
// Type implementation
// -------------------------------------------------------------------------------------------------

type Bytes struct {
	value []byte
}

func NewBytes(value []byte) *Bytes {
	return &Bytes{value}
}

func (v *Bytes) CloneIfPrimitive() execute.Value {
	return NewBytes(v.value[:])
}

func (v *Bytes) CompareTo(o execute.Value) (int, bool) {
	if ob, ok := o.(*Bytes); ok {
		return bytes.Compare(v.value, ob.value), true
	}
	return 0, false
}

func (v *Bytes) Equals(o execute.Value) bool {
	if ob, ok := o.(*Bytes); ok {
		return bytes.Equal(v.value, ob.value)
	}
	return false
}

func (v *Bytes) GetAttribute(a string) (execute.Value, error) {
	return nil, errors.NewAttributeError(v.Type(), a)
}

func (v *Bytes) GetIndex(i execute.Value) (execute.Value, error) {
	idx, err := numericIndex(i, v.Type())
	if err != nil {
		return nil, err
	}
	idx, ok := normalizeIndex(idx, len(v.value))
	if !ok {
		return nil, errors.NewIndexError(fmt.Sprintf("%d", idx))
	}
	return NewBytes([]byte{v.value[idx]}), nil
}

func (v *Bytes) HasAttribute(a string) bool {
	return false
}

func (v *Bytes) HashBytes() ([]byte, error) {
	return []byte(v.value), nil
}

func (v *Bytes) Length() (uint64, error) {
	return uint64(len(v.value)), nil
}

func (v *Bytes) SetAttribute(a string, _ execute.Value) error {
	return errors.NewAttributeError(v.Type(), a)
}

func (v *Bytes) SetIndex(execute.Value, execute.Value) error {
	return errors.SetIndexNotSupported(v.Type())
}

func (v *Bytes) String() string {
	return fmt.Sprintf("0x%X", v.value)
}

func (v *Bytes) ToBool() bool {
	// A bytes object with all null bytes is falsey.
	for _, b := range v.value {
		if b != 0 {
			return true
		}
	}
	return false
}

func (v *Bytes) ToBytes() ([]byte, error) {
	return append([]byte{}, v.value...), nil
}

func (v *Bytes) ToCallable() (execute.Callable, error) {
	return nil, errors.NewTypeError(v.Type(), FuncType)
}

func (v *Bytes) ToFloat() (float64, error) {
	bits := bytesTo64Bits(v.value)
	return math.Float64frombits(binary.BigEndian.Uint64(bits)), nil
}

func (v *Bytes) ToInt() (int64, error) {
	bits := bytesTo64Bits(v.value)
	var n int64
	buf := bytes.NewBuffer(bits)
	binary.Read(buf, binary.BigEndian, &n)
	return n, nil
}

func (v *Bytes) ToIterator() (execute.Iterator, error) {
	return &bytesIterator{b: v}, nil
}

func (v *Bytes) ToStr() (string, error) {
	if !utf8.Valid(v.value) {
		return "", errors.NewValueError("bytes are not valid UTF-8")
	}
	return string(v.value), nil
}

func (v *Bytes) ToUint() (uint64, error) {
	bits := bytesTo64Bits(v.value)
	return binary.BigEndian.Uint64(bits), nil
}

func (v *Bytes) Type() execute.Type {
	return BytesType
}

type bytesIterator struct {
	idx int
	b   *Bytes
}

func (bi *bytesIterator) HasNext() bool {
	return bi.idx < len(bi.b.value)
}

func (bi *bytesIterator) Next() (execute.Value, error) {
	c := bi.b.value[bi.idx]
	bi.idx++
	return NewBytes([]byte{c}), nil
}

func bytesTo64Bits(b []byte) []byte {
	bits := make([]byte, 8)
	for i := 0; i < 8; i++ {
		if i < len(b) {
			bits[i] = b[i]
		} else {
			bits[i] = 0
		}
	}
	return bits
}
