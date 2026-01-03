package errors_test

import (
	"testing"

	"github.com/chrispyles/slow/internal/errors"
)

func TestNameError(t *testing.T) {
	e := errors.NewNameError("foo")

	got, want := e.Error(), "NameError: no variable \"foo\" has been declared"
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}
