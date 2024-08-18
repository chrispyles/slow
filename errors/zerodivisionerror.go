package errors

type ZeroDivisionError struct{}

func NewZeroDivisionError() *ZeroDivisionError {
	return &ZeroDivisionError{}
}

func (e *ZeroDivisionError) Error() string {
	return "ZeroDivisionError: attempted to divide by zero"
}
