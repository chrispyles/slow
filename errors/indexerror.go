package errors

import "fmt"

type IndexError struct {
	key string
}

func NewIndexError(key string) *IndexError {
	return &IndexError{key: key}
}

func (e *IndexError) Error() string {
	return fmt.Sprintf("IndexError: invalid index value: %s", e.key)
}
