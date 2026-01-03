package errors_test

import (
	"testing"

	"github.com/chrispyles/slow/src/errors"
	slowtesting "github.com/chrispyles/slow/src/testing"
)

func TestSyntaxError(t *testing.T) {
	mb := &slowtesting.MockBuffer{LineNumberRet: 3}

	e := errors.NewSyntaxError(mb, "foo", "bar")

	got, want := e.Error(), "SyntaxError: foo on line 3: \"bar\""
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}

	e = errors.NewSyntaxError(mb, "foo", "")

	got, want = e.Error(), "SyntaxError: foo on line 3"
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}

func TestUnexpectedSymbolError(t *testing.T) {
	mb := &slowtesting.MockBuffer{LineNumberRet: 3}

	e := errors.UnexpectedSymbolError(mb, "foo", "bar")

	got, want := e.Error(), "SyntaxError: unexpected symbol, expected \"bar\" on line 3: \"foo\""
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}

	e = errors.UnexpectedSymbolError(mb, "foo", "")

	got, want = e.Error(), "SyntaxError: unexpected symbol on line 3: \"foo\""
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}

func TestInterpreterSyntaxError(t *testing.T) {
	e := errors.InterpreterSyntaxError("foo")

	got, want := e.Error(), "SyntaxError: foo"
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}
