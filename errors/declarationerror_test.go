package errors_test

import (
	"testing"

	"github.com/chrispyles/slow/errors"
)

func TestDeclarationError(t *testing.T) {
	e := errors.NewDeclarationError("foo")

	got, want := e.Error(), "DeclarationError: variable \"foo\" has already been declared"
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}
