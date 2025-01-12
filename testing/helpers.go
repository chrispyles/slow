package testing

import (
	"testing"

	"github.com/chrispyles/slow/execute"
)

func MustMakeEnv(t *testing.T, vars map[string]execute.Value) *execute.Environment {
	e := execute.NewEnvironment()
	for k, v := range vars {
		if err := e.Declare(k); err != nil {
			t.Fatalf("failed to declare variable %q: %v", k, err)
		}
		if _, err := e.Set(k, v); err != nil {
			t.Fatalf("failed to set variable %q: %v", k, err)
		}
	}
	return e
}
