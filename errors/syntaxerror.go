package errors

import "fmt"

type SyntaxError struct {
	message     string
	symbol      string
	lineNumber  int
	interpreter bool
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

func InterpreterSyntaxError(msg string) *SyntaxError {
	return &SyntaxError{message: msg, interpreter: true}
}

func (e *SyntaxError) Error() string {
	var end string
	if e.symbol != "" {
		end = fmt.Sprintf(": %q", e.symbol)
	}
	var line string
	if !e.interpreter {
		line = fmt.Sprintf(" on line %d", e.lineNumber)
	}
	return fmt.Sprintf("SyntaxError%s: %s%s", line, e.message, end)
}
