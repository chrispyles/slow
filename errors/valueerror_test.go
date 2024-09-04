package errors_test

import (
	goerrors "errors"
	"testing"

	"github.com/chrispyles/slow/errors"
	slowtesting "github.com/chrispyles/slow/testing"
)

func TestValueError(t *testing.T) {
	e := errors.NewValueError("foo")

	got, want := e.Error(), "ValueError: foo"
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}

func TestWrapValueError(t *testing.T) {
	e := errors.WrapValueError("foo", slowtesting.NewMockType(), goerrors.New("doh"))

	got, want := e.Error(), "ValueError: error converting \"foo\" to type \"MockType\": doh: doh"
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}
