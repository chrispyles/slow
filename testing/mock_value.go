package testing

import (
	"fmt"

	"github.com/chrispyles/slow/execute"
)

func (m *MockType) String() string {
	return m.StringRet
}

type MockValue struct {
	EqualsRet bool

	Attributes map[string]execute.Value

	LengthRet uint64
	LengthErr error

	StringRet string

	ToBoolRet bool

	ToBytesRet []byte
	ToBytesErr error

	ToCallableRet execute.Callable
	ToCallableErr error

	ToFloatRet float64
	ToFloatErr error

	ToIntRet int64
	ToIntErr error

	ToIteratorRet execute.Iterator
	ToIteratorErr error

	ToStrRet string
	ToStrErr error

	ToUintRet uint64
	ToUintErr error

	TypeRet *MockType
}

func (m *MockValue) CloneIfPrimitive() execute.Value {
	return m
}

func (m *MockValue) CompareTo(execute.Value) (int, bool) {
	return 0, false
}

func (m *MockValue) Equals(execute.Value) bool {
	return m.EqualsRet
}

func (m *MockValue) GetAttribute(a string) (execute.Value, error) {
	if m.Attributes == nil {
		return nil, fmt.Errorf("attribute error: type %q, attribute %q", m.Type(), a)
	}
	ret, ok := m.Attributes[a]
	if !ok {
		return nil, fmt.Errorf("attribute error: type %q, attribute %q", m.Type(), a)
	}
	return ret, nil
}

func (m *MockValue) GetIndex(execute.Value) (execute.Value, error) {
	return nil, nil
}

func (m *MockValue) HasAttribute(a string) bool {
	return false
}

func (m *MockValue) HashBytes() ([]byte, error) {
	return nil, nil
}

func (m *MockValue) Length() (uint64, error) {
	return m.LengthRet, m.LengthErr
}

func (m *MockValue) SetAttribute(string, execute.Value) error {
	return nil
}

func (m *MockValue) SetIndex(execute.Value, execute.Value) error {
	return nil
}

func (m *MockValue) String() string {
	return m.StringRet
}

func (m *MockValue) ToBool() bool {
	return m.ToBoolRet
}

func (m *MockValue) ToBytes() ([]byte, error) {
	return m.ToBytesRet, m.ToBytesErr
}

func (m *MockValue) ToCallable() (execute.Callable, error) {
	return m.ToCallableRet, m.ToCallableErr
}

func (m *MockValue) ToFloat() (float64, error) {
	return m.ToFloatRet, m.ToFloatErr
}

func (m *MockValue) ToInt() (int64, error) {
	return m.ToIntRet, m.ToIntErr
}

func (m *MockValue) ToIterator() (execute.Iterator, error) {
	return m.ToIteratorRet, m.ToIteratorErr
}

func (m *MockValue) ToStr() (string, error) {
	return m.ToStrRet, m.ToStrErr
}

func (m *MockValue) ToUint() (uint64, error) {
	return m.ToUintRet, m.ToUintErr
}

func (m *MockValue) Type() execute.Type {
	return m.TypeRet
}
