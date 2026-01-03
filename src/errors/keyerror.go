package errors

import "fmt"

func NewKeyError(key string) error {
	return newError("KeyError", fmt.Sprintf("map has no key %q", key))
}

func MapModifiedDuringIterationError(key string) error {
	return newError("KeyError", fmt.Sprintf("map was modified during iteration and no longer has key %q", key))
}
