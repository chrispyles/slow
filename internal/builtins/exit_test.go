package builtins

import (
	"os"
	"testing"

	"github.com/chrispyles/slow/internal/execute"
	"github.com/chrispyles/slow/internal/types"
)

func TestBuiltins_exit(t *testing.T) {
	makeExitMock := func() []any {
		calls := make([]any, 1)
		var i int
		osExit = func(c int) {
			calls[i] = c
			i++
		}
		return calls
	}
	cleanupExitMock := func() {
		osExit = os.Exit
	}
	doBuiltinTest(t, []builtinTest{
		{
			name:        "no_args",
			fn:          "exit",
			args:        []execute.Value{},
			makeMock:    makeExitMock,
			cleanupMock: cleanupExitMock,
			want:        types.Null,
			wantPrints:  []string{"Exiting with code 0\n"},
			wantCalls:   []any{0},
		},
		{
			name:        "one_arg",
			fn:          "exit",
			args:        []execute.Value{types.NewInt(1)},
			makeMock:    makeExitMock,
			cleanupMock: cleanupExitMock,
			want:        types.Null,
			wantPrints:  []string{"Exiting with code 1\n"},
			wantCalls:   []any{1},
		},
		{
			name:        "one_float",
			fn:          "exit",
			args:        []execute.Value{types.NewFloat(1.2)},
			makeMock:    makeExitMock,
			cleanupMock: cleanupExitMock,
			want:        types.Null,
			wantPrints:  []string{"Exiting with code 1\n"},
			wantCalls:   []any{1},
		},
		{
			name:        "arg_cant_be_converted_to_int",
			fn:          "exit",
			args:        []execute.Value{types.NewList(nil)},
			makeMock:    makeExitMock,
			cleanupMock: cleanupExitMock,
			want:        types.Null,
			wantPrints:  []string{"Exiting with code 1\n"},
			wantCalls:   []any{1},
		},
		{
			name:        "falsey_non_numeric_arg",
			fn:          "exit",
			args:        []execute.Value{types.NewBytes(nil)},
			makeMock:    makeExitMock,
			cleanupMock: cleanupExitMock,
			want:        types.Null,
			wantPrints:  []string{"Exiting with code 0\n"},
			wantCalls:   []any{0},
		},
		{
			name:        "multiple_args",
			fn:          "exit",
			args:        []execute.Value{types.NewInt(2), types.NewInt(3)},
			makeMock:    makeExitMock,
			cleanupMock: cleanupExitMock,
			want:        types.Null,
			wantPrints:  []string{"Exiting with code 2\n"},
			wantCalls:   []any{2},
		},
	})
}
