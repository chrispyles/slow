package errors

import "fmt"

type SyntaxError struct {
	message    string
	symbol     string
	lineNumber int
}

func NewSyntaxError(buf Buffer, message string, symbol string) *SyntaxError {
	return &SyntaxError{message: message, symbol: symbol, lineNumber: buf.LineNumber()}
}

func UnexpectedSymbolError(buf Buffer, got, want string) *SyntaxError {
	msg := "unexpected symbol"
	if want != "" {
		msg = fmt.Sprintf("unexpected symbol, expected %q", want)
	}
	return NewSyntaxError(buf, msg, got)
}

func (e *SyntaxError) Error() string {
	var end string
	if e.symbol != "" {
		end = fmt.Sprintf(": %q", e.symbol)
	}
	return fmt.Sprintf("SyntaxError on line %d: %s%s", e.lineNumber, e.message, end)
}
