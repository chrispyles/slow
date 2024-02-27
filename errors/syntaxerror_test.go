package errors_test

import (
	"testing"

	"github.com/chrispyles/slow/errors"
	slowtesting "github.com/chrispyles/slow/testing"
)

func TestSyntaxError(t *testing.T) {
	mb := &slowtesting.MockBuffer{LineNumberRet: 3}

	e := errors.NewSyntaxError(mb, "foo", "bar")

	got, want := e.Error(), "SyntaxError on line 3: foo: \"bar\""
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}

	e = errors.NewSyntaxError(mb, "foo", "")

	got, want = e.Error(), "SyntaxError on line 3: foo"
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}

func TestUnexpectedSymbolError(t *testing.T) {
	mb := &slowtesting.MockBuffer{LineNumberRet: 3}

	e := errors.UnexpectedSymbolError(mb, "foo", "bar")

	got, want := e.Error(), "SyntaxError on line 3: unexpected symbol, expected \"bar\": \"foo\""
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}

	e = errors.UnexpectedSymbolError(mb, "foo", "")

	got, want = e.Error(), "SyntaxError on line 3: unexpected symbol: \"foo\""
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}
