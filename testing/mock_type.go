package testing

import "github.com/chrispyles/slow/execute"

type MockType struct {
	IsNumericRet  bool
	MatchingTypes map[execute.Type]bool
	StringRet     string
}

func (m *MockType) IsNumeric() bool {
	return m.IsNumericRet
}

func (m *MockType) Matches(t execute.Type) bool {
	if m.MatchingTypes == nil {
		return false
	}
	return m.MatchingTypes[t]
}

func NewMockType() *MockType {
	return &MockType{StringRet: "MockType"}
}
