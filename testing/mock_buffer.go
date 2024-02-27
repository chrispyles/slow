package testing

type MockBuffer struct {
	LineNumberRet int
}

func (m *MockBuffer) LineNumber() int {
	return m.LineNumberRet
}
