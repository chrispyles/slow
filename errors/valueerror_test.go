package errors_test

import (
	"testing"

	"github.com/chrispyles/slow/errors"
)

func TestValueError(t *testing.T) {
	e := errors.NewValueError("foo")

	got, want := e.Error(), "ValueError: foo"
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}
