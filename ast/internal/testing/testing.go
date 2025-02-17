package testing

import (
	"testing"

	"github.com/chrispyles/slow/execute"
	slowtesting "github.com/chrispyles/slow/testing"
	"github.com/google/go-cmp/cmp"
)

type TestCase struct {
	Name        string
	Node        execute.Expression
	Env         *execute.Environment
	Want        execute.Value
	WantEnv     *execute.Environment
	WantSameEnv bool
	WantErr     error
}

func RunTestCase(t *testing.T, tc TestCase) {
	t.Run(tc.Name, func(t *testing.T) {
		wantEnv := tc.WantEnv
		if tc.WantSameEnv {
			wantEnv = tc.Env.Copy()
		}
		got, err := tc.Node.Execute(tc.Env)
		if diff := cmp.Diff(tc.WantErr, err, slowtesting.AllowUnexported()); diff != "" {
			t.Errorf("Execute() returned incorrect error (-want +got):\n%s", diff)
		}
		if diff := cmp.Diff(tc.Want, got, slowtesting.AllowUnexported()); diff != "" {
			t.Errorf("Execute() returned unexpected diff (-want +got):\n%s", diff)
		}
		if diff := cmp.Diff(wantEnv, tc.Env, slowtesting.AllowUnexported(), slowtesting.EquateFuncs()); diff != "" {
			t.Errorf("env after Execute() has unexpected diff (-want +got):\n%s", diff)
		}
	})
}
