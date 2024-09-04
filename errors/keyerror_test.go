package errors_test

import (
	"testing"

	"github.com/chrispyles/slow/errors"
)

func TestNewKeyError(t *testing.T) {
	e := errors.NewKeyError("foo")

	got, want := e.Error(), "KeyError: map has no key \"foo\""
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}

func TestMapModifiedDuringIterationError(t *testing.T) {
	e := errors.MapModifiedDuringIterationError("foo")

	got, want := e.Error(), "KeyError: map was modified during iteration and no longer has key \"foo\""
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}
