package testing

import (
	"github.com/chrispyles/slow/execute"
)

type MockType struct {
	IsNumericRet  bool
	MatchingTypes map[execute.Type]bool
	StringRet     string
}

func (m *MockType) IsNumeric() bool {
	return m.IsNumericRet
}

func (m *MockType) New(t execute.Value) (execute.Value, error) {
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
