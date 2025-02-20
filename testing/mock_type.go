package testing

import (
	"github.com/chrispyles/slow/execute"
)

type MockType struct {
	IsNumericRet bool
	NewRet       execute.Value
	NewErr       error
	StringRet    string
}

func (m *MockType) IsNumeric() bool {
	return m.IsNumericRet
}

func (m *MockType) New(t execute.Value) (execute.Value, error) {
	if m.NewErr != nil {
		return nil, m.NewErr
	}
	if m.NewRet != nil {
		return m.NewRet, nil
	}
	return &MockValue{}, nil
}

func (m *MockType) String() string {
	if m == nil {
		return "MockType"
	}
	return m.StringRet
}

func NewMockType() *MockType {
	return &MockType{StringRet: "MockType"}
}
