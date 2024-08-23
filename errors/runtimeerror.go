package errors

import "fmt"

type RuntimeError struct {
	msg string
}

func NewRuntimeError(msg string) *RuntimeError {
	return &RuntimeError{msg}
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("RuntimeError: %s", e.msg)
}
