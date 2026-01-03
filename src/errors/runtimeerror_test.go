package errors_test

import (
	"testing"

	"github.com/chrispyles/slow/src/errors"
)

func TestNewRuntimeError(t *testing.T) {
	e := errors.NewRuntimeError("foo")

	got, want := e.Error(), "RuntimeError: foo"
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}
