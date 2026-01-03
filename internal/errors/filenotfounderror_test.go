package errors_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/chrispyles/slow/internal/errors"
)

func TestNewFileNotFoundError(t *testing.T) {
	e := errors.NewFileNotFoundError("foo.txt")
	want := "FileNotFoundError: file \"foo.txt\" does not exist"
	if got := e.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestWrapFileError(t *testing.T) {
	err := fmt.Errorf("foo: %w", os.ErrNotExist)
	e := errors.WrapFileError(err, "foo.txt")
	want := "FileNotFoundError: file \"foo.txt\" does not exist"
	if got := e.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}

	err = fmt.Errorf("foo: %w", os.ErrPermission)
	e = errors.WrapFileError(err, "foo.txt")
	want = "FileError: file \"foo.txt\": foo: permission denied"
	if got := e.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}
