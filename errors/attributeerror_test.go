package errors_test

import (
	"testing"

	"github.com/chrispyles/slow/errors"
	slowtesting "github.com/chrispyles/slow/testing"
)

func TestAttributeError(t *testing.T) {
	mt := slowtesting.NewMockType()

	e := errors.NewAttributeError(mt, "foo")

	got, want := e.Error(), "AttributeError: type \"MockType\" has no attribute \"foo\""
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}
