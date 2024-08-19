package errors

import "fmt"

type KeyError struct {
	key string
	msg string
}

func NewKeyError(key string) *KeyError {
	return &KeyError{key: key}
}

func MapModifiedDuringIterationError(key string) *KeyError {
	return &KeyError{msg: fmt.Sprintf("KeyError: map was modified during iteration and no longer has key %s", key)}
}

func (e *KeyError) Error() string {
	if e.msg != "" {
		return e.msg
	}
	return fmt.Sprintf("KeyError: map has no key %s", e.key)
}
