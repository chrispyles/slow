package errors

import "fmt"

type SlowError struct {
	errType string
	msg     string
}

func newError(t, m string) *SlowError {
	return &SlowError{t, m}
}

func (e *SlowError) Error() string {
	return fmt.Sprintf("%s: %s", e.errType, e.msg)
}
