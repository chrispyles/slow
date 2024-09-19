package builtins

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/printer"
	slowtesting "github.com/chrispyles/slow/testing"
	"github.com/chrispyles/slow/types"
	"github.com/google/go-cmp/cmp"
)

var allowUnexported = cmp.AllowUnexported(
	errors.SlowError{},
	execute.Environment{},
	rangeGenerator{},
	types.Bool{},
	types.Float{},
	types.Func{},
	types.Generator{},
	types.Int{},
	types.Iterator{},
	types.List{},
	types.Module{},
	types.Str{},
	types.Uint{},
)

// Adapted from https://github.com/google/go-cmp/issues/162
var equateFuncs = cmp.Comparer(func(x, y types.FuncImpl) bool {
	px := *(*unsafe.Pointer)(unsafe.Pointer(&x))
	py := *(*unsafe.Pointer)(unsafe.Pointer(&y))
	return px == py
})

type builtinTest struct {
	name         string
	fn           string
	args         []execute.Value
	makeMock     func() []any
	cleanupMock  func()
	want         execute.Value
	wantPrintlns []string
	wantCalls    []any
	wantErr      error
}

func doBuiltinTest(t *testing.T, tests []builtinTest) {
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var printed []string
			println = func(s string) {
				printed = append(printed, s)
			}
			printlnf = func(s string, a ...any) {
				printed = append(printed, fmt.Sprintf(s, a...))
			}
			t.Cleanup(func() {
				println = printer.Println
				printlnf = printer.Printlnf
			})
			var gotCalls []any
			if tc.makeMock != nil {
				gotCalls = tc.makeMock()
				t.Cleanup(tc.cleanupMock)
			}
			env := RootEnvironment.NewFrame()
			fn, err := env.Get(tc.fn)
			if err != nil {
				t.Fatalf("Get() returned unexpected error: %v", err)
			}
			c, err := fn.ToCallable()
			if err != nil {
				t.Fatalf("fn.ToCallable() returned unexpected error: %v", err)
			}
			got, err := c.Call(env, tc.args...)
			if diff := cmp.Diff(tc.wantErr, err, allowUnexported); diff != "" {
				t.Errorf("c.Call() returned incorrect error (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.want, got, allowUnexported, equateFuncs); diff != "" {
				t.Errorf("c.Call() returned incorrect value (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.wantPrintlns, printed); diff != "" {
				t.Errorf("println called incorrectly (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.wantCalls, gotCalls); diff != "" {
				t.Errorf("mocked function called incorrectly (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNewRootEnvironmentIsFrozen(t *testing.T) {
	// Attempt to reassign a variable that is bound to a built-in, so we know it's already declared.
	_, err := RootEnvironment.Set("import", &slowtesting.MockValue{})
	if err == nil {
		t.Errorf("env.Set did not error")
	}
}
