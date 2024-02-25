package errors

import (
	"fmt"
)

type EOFError struct {
	lineNumber int
}

func NewEOFError(buf Buffer) *EOFError {
	return &EOFError{buf.LineNumber()}
}

func (e *EOFError) Error() string {
	return fmt.Sprintf("EOFError: ran out of input on line %d", e.lineNumber)
}
