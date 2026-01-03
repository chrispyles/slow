package eval

import (
	"errors"
	"reflect"
	"testing"

	"github.com/chrispyles/slow/src/execute"
	slowtesting "github.com/chrispyles/slow/src/testing"
	"github.com/chrispyles/slow/src/types"
	"github.com/google/go-cmp/cmp"
)

var mockValueStringRet = "mockASTMockValue"

func TestEval(t *testing.T) {
	origMakeAST := makeAST
	origPrintln := println
	makeMakeAST := func(err error) ([]string, *mockAST) {
		calls := make([]string, 1)
		mast := &mockAST{}
		i := 0
		makeAST = func(s string) (execute.AST, error) {
			calls[i] = s
			i++
			return mast, err
		}
		return calls, mast
	}
	makeMockPrintln := func() []*string {
		calls := make([]*string, 1)
		i := 0
		println = func(s string) {
			calls[i] = &s
			i++
		}
		return calls
	}
	tests := []struct {
		name             string
		in               string
		env              *execute.Environment
		printExprValue   bool
		astExecRet       execute.Value
		astExecErr       error
		wantPrintlnCalls []*string
	}{
		{
			name:             "success",
			in:               "some code",
			env:              execute.NewEnvironment(),
			wantPrintlnCalls: make([]*string, 1),
		},
		{
			name:             "success_with_print",
			in:               "some code",
			env:              execute.NewEnvironment(),
			printExprValue:   true,
			wantPrintlnCalls: []*string{&mockValueStringRet},
		},
		{
			name:             "no_print_null",
			in:               "some code",
			env:              execute.NewEnvironment(),
			printExprValue:   true,
			astExecRet:       types.Null,
			wantPrintlnCalls: make([]*string, 1),
		},
		{
			name:             "ast_exec_error",
			in:               "some code",
			env:              execute.NewEnvironment(),
			printExprValue:   true,
			astExecErr:       errors.New("nuh-uh"),
			wantPrintlnCalls: make([]*string, 1),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				makeAST = origMakeAST
				println = origPrintln
			})
			makeASTCalls, mast := makeMakeAST(nil)
			mast.ret = tc.astExecRet
			mast.err = tc.astExecErr
			printlnCalls := makeMockPrintln()
			Eval(tc.in, tc.env, tc.printExprValue)
			if diff := cmp.Diff([]string{tc.in}, makeASTCalls); diff != "" {
				t.Errorf("Eval() called makeAST incorrectly (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff([]uintptr{reflect.ValueOf(tc.env).Pointer()}, mast.calls); diff != "" {
				t.Errorf("Eval() called ast.Execute incorrectly (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.wantPrintlnCalls, printlnCalls); diff != "" {
				t.Errorf("Eval() called println incorrectly (-want +got):\n%s", diff)
			}
		})
	}
}

type mockAST struct {
	calls []uintptr
	ret   execute.Value
	err   error
}

func (m *mockAST) Execute(env *execute.Environment) (execute.Value, error) {
	m.calls = append(m.calls, reflect.ValueOf(env).Pointer())
	if m.err != nil {
		return nil, m.err
	}
	if m.ret != nil {
		return m.ret, nil
	}
	return &slowtesting.MockValue{StringRet: mockValueStringRet}, nil
}

func (m *mockAST) String() string {
	return "mockAST"
}
