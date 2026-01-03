package errors

func NewRuntimeError(msg string) error {
	return newError("RuntimeError", msg)
}
