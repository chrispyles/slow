package errors

import "fmt"

func NewIndexError(key string) error {
	return newError("IndexError", fmt.Sprintf("invalid index value: %s", key))
}

func IndexingNotSupported(t Type) error {
	return newError("IndexError", fmt.Sprintf("type %q does not support indexing", t.String()))
}

func NonNumericIndexError(indexType, containerType Type) error {
	return newError("IndexError", fmt.Sprintf("type %q can't be used as an index in type %q", indexType.String(), containerType.String()))
}

func SetIndexNotSupported(t Type) error {
	return newError("IndexError", fmt.Sprintf("type %q doesn't support index assignment", t.String()))
}
