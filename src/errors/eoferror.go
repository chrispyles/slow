package errors

import (
	"fmt"
)

func NewEOFError(buf Buffer) error {
	return newError("EOFError", fmt.Sprintf("ran out of input on line %d", buf.LineNumber()))
}
