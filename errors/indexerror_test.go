package errors_test

import (
	"testing"

	"github.com/chrispyles/slow/errors"
)

func TestIndexError(t *testing.T) {
	e := errors.NewIndexError("foo")

	got, want := e.Error(), "IndexError: invalid index value: foo"
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}
