package errors

import "fmt"

func NewSyntaxError(buf Buffer, message string, symbol string) error {
	var end string
	if symbol != "" {
		end = fmt.Sprintf(": %q", symbol)
	}
	return newError("SyntaxError", fmt.Sprintf("%s on line %d%s", message, buf.LineNumber(), end))
}

func UnexpectedSymbolError(buf Buffer, got, want string) error {
	var msg string
	if want != "" {
		msg = fmt.Sprintf("unexpected symbol, expected %q on line %d", want, buf.LineNumber())
	} else {
		msg = fmt.Sprintf("unexpected symbol on line %d", buf.LineNumber())
	}
	if got != "" {
		msg += fmt.Sprintf(": %q", got)
	}
	return newError("SyntaxError", msg)
}

func InterpreterSyntaxError(msg string) error {
	return newError("SyntaxError", msg)
}
