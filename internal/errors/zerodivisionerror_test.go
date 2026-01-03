package errors_test

import (
	"testing"

	"github.com/chrispyles/slow/internal/errors"
)

func TestZeroDivisionError(t *testing.T) {
	e := errors.NewZeroDivisionError()

	got, want := e.Error(), "ZeroDivisionError: attempted to divide by zero"
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}
