package errors

import "fmt"

type ValueError struct {
	message string
}

func NewValueError(message string) *ValueError {
	return &ValueError{message: message}
}

func (e *ValueError) Error() string {
	return fmt.Sprintf("ValueError: %s", e.message)
}
