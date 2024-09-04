package errors_test

import (
	"testing"

	"github.com/chrispyles/slow/errors"
)

func TestNewImportError(t *testing.T) {
	e := errors.NewImportError("my_mod")
	want := "ImportError: no such module \"my_mod\""
	if got := e.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}
