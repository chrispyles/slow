package errors

import (
	goerrors "errors"
	"fmt"
	"os"
)

func NewFileNotFoundError(path string) error {
	return newError("FileNotFoundError", fmt.Sprintf("file %q does not exist", path))
}

func WrapFileError(err error, path string) error {
	if goerrors.Is(err, os.ErrNotExist) {
		return NewFileNotFoundError(path)
	}
	return wrapError("FileError", fmt.Sprintf("file %q", path), err)
}
