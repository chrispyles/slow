package errors

import (
	"errors"
	"testing"
)

func TestSlowError(t *testing.T) {
	e := newError("FooError", "a message")
	want := "FooError: a message"
	if got := e.Error(); got != want {
		t.Errorf("e.Error() returned %q, want %q", got, want)
	}
}

func TestWrapError(t *testing.T) {
	err := errors.New("foobar")
	e := wrapError("FooError", "a message", err)
	want := "FooError: a message: foobar"
	if got := e.Error(); got != want {
		t.Errorf("e.Error() returned %q, want %q", got, want)
	}
}

func TestWrapNilError(t *testing.T) {
	got := wrapError("FooError", "msg", nil)
	if got != nil {
		t.Errorf("wrapError() returned non-nil error: %v", got)
	}
}
