package builtins

import (
	"testing"

	"github.com/chrispyles/slow/internal/execute"
	"github.com/chrispyles/slow/internal/printer"
	slowtesting "github.com/chrispyles/slow/internal/testing"
	slowcmpopts "github.com/chrispyles/slow/internal/testing/cmpopts"
	"github.com/google/go-cmp/cmp"
)

type builtinTest struct {
	name        string
	fn          string
	args        []execute.Value
	makeMock    func() []any
	cleanupMock func()
	want        execute.Value
	wantPrints  []string
	wantCalls   []any
	wantErr     error
}

func doBuiltinTest(t *testing.T, tests []builtinTest) {
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var printed []string
			printer.Set(func(s string) {
				printed = append(printed, s)
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
			if diff := cmp.Diff(tc.wantErr, err, slowcmpopts.AllowUnexported()); diff != "" {
				t.Errorf("c.Call() returned incorrect error (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.want, got, slowcmpopts.AllowUnexported(), slowcmpopts.EquateFuncs()); diff != "" {
				t.Errorf("c.Call() returned incorrect value (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.wantPrints, printed); diff != "" {
				t.Errorf("println called incorrectly (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.wantCalls, gotCalls, slowcmpopts.AllowUnexported()); diff != "" {
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
