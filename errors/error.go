package errors

import "fmt"

type SlowError struct {
	errType string
	msg     string
	wrapped error
}

func newError(t, m string) *SlowError {
	return &SlowError{t, m, nil}
}

func wrapError(t string, m string, err error) error {
	if err == nil {
		return nil
	}
	return &SlowError{t, fmt.Sprintf("%s: %+v", m, err), err}
}

func (e *SlowError) Error() string {
	return fmt.Sprintf("%s: %s", e.errType, e.msg)
}
