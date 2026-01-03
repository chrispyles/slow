package testing

import (
	"reflect"

	"github.com/chrispyles/slow/internal/execute"
)

type MockExpression struct {
	ExecuteRet execute.Value
	ExecuteErr error
	Calls      []uintptr
}

func (m *MockExpression) Execute(env *execute.Environment) (execute.Value, error) {
	m.Calls = append(m.Calls, reflect.ValueOf(env).Pointer())
	return m.ExecuteRet, m.ExecuteErr
}
