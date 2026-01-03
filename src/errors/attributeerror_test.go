package errors_test

import (
	"testing"

	"github.com/chrispyles/slow/src/errors"
	slowtesting "github.com/chrispyles/slow/src/testing"
)

func TestAttributeError(t *testing.T) {
	mt := slowtesting.NewMockType()

	e := errors.NewAttributeError(mt, "foo")

	got, want := e.Error(), "AttributeError: type \"MockType\" has no attribute \"foo\""
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}

func TestAssignmentError(t *testing.T) {
	mt := slowtesting.NewMockType()

	e := errors.AssignmentError(mt, "foo")

	got, want := e.Error(), "AttributeError: can't reassign attribute \"foo\" in type \"MockType\""
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}
