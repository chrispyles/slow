package errors_test

import (
	"testing"

	"github.com/chrispyles/slow/internal/errors"
	slowtesting "github.com/chrispyles/slow/internal/testing"
)

func TestIndexError(t *testing.T) {
	e := errors.NewIndexError("foo")

	got, want := e.Error(), "IndexError: invalid index value: foo"
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}

func TestIndexingNotSupported(t *testing.T) {
	e := errors.IndexingNotSupported(slowtesting.NewMockType())
	want := "IndexError: type \"MockType\" does not support indexing"
	if got := e.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestNonNumericIndexError(t *testing.T) {
	e := errors.NonNumericIndexError(&slowtesting.MockType{StringRet: "t1"}, &slowtesting.MockType{StringRet: "t2"})
	want := "IndexError: type \"t1\" can't be used as an index in type \"t2\""
	if got := e.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestSetIndexNotSupported(t *testing.T) {
	e := errors.SetIndexNotSupported(slowtesting.NewMockType())
	want := "IndexError: type \"MockType\" doesn't support index assignment"
	if got := e.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}
