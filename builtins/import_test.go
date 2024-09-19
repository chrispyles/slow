package builtins

import (
	"fmt"
	"os"
	"testing"

	"github.com/chrispyles/slow/builtins/modules"
	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
	slowtesting "github.com/chrispyles/slow/testing"
	"github.com/chrispyles/slow/types"
)

func TestBuiltins_import(t *testing.T) {
	tests := []builtinTest{}
	for _, name := range modules.AllModules {
		mod, ok := modules.Get(name)
		if !ok {
			t.Fatalf("failed to get builtin module: %s", name)
		}
		modEnv, err := mod.Import()
		if err != nil {
			t.Fatalf("failed to import builtin module %q: %v", name, err)
		}
		tests = append(tests, builtinTest{
			name: fmt.Sprintf("builtin_%s", name),
			fn:   "import",
			args: []execute.Value{types.NewStr(name)},
			want: types.NewModule(name, modEnv),
		})
	}
	mockOsReadFile := func() []any {
		calls := make([]any, 2)
		i := 0
		osReadFile = func(name string) ([]byte, error) {
			calls[i] = name
			i++
			return []byte("this is foobar.slo"), nil
		}
		evalEval = func(c string, _ *execute.Environment, _ bool) {
			calls[i] = c
			i++
		}
		return calls
	}
	cleanupMock := func() {
		osReadFile = os.ReadFile
	}
	tests = append(tests, []builtinTest{
		{
			name:    "no_args",
			fn:      "import",
			args:    []execute.Value{},
			wantErr: errors.CallError("import", 0, 1),
		},
		{
			name:    "too_many_args",
			fn:      "import",
			args:    []execute.Value{&slowtesting.MockValue{}, &slowtesting.MockValue{}},
			wantErr: errors.CallError("import", 2, 1),
		},
		{
			name:    "non_string_arg",
			fn:      "import",
			args:    []execute.Value{types.NewInt(1)},
			wantErr: errors.NewTypeError(types.IntType, types.StrType),
		},
		{
			name:    "nonexistent_module",
			fn:      "import",
			args:    []execute.Value{types.NewStr("foobar")},
			wantErr: errors.NewImportError("foobar"),
		},
		{
			name:        "relative_import",
			fn:          "import",
			args:        []execute.Value{types.NewStr("foobar.slo")},
			want:        types.NewModule("foobar.slo", RootEnvironment.NewFrame()),
			makeMock:    mockOsReadFile,
			cleanupMock: cleanupMock,
			wantCalls:   []any{"foobar.slo", "this is foobar.slo"},
		},
	}...)
	doBuiltinTest(t, tests)
}
