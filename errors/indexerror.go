package errors

import "fmt"

func NewIndexError(key string) error {
	return newError("IndexError", fmt.Sprintf("invalid index value: %s", key))
}
