package errors_test

import (
	"testing"

	"github.com/chrispyles/slow/errors"
)

func TestNameError(t *testing.T) {
	e := errors.NewNameError("foo")

	got, want := e.Error(), "NameError: no such variable \"foo\""
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}
