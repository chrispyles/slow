package errors_test

import (
	"testing"

	"github.com/chrispyles/slow/internal/errors"
	slowtesting "github.com/chrispyles/slow/internal/testing"
)

func TestEOFError(t *testing.T) {
	mb := &slowtesting.MockBuffer{LineNumberRet: 3}

	e := errors.NewEOFError(mb)

	got, want := e.Error(), "EOFError: ran out of input on line 3"
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}
