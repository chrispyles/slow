package errors

func NewZeroDivisionError() error {
	return newError("ZeroDivisionError", "attempted to divide by zero")
}
