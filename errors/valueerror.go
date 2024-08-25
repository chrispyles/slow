package errors

func NewValueError(msg string) error {
	return newError("ValueError", msg)
}
