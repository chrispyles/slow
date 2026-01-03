package errors

type Buffer interface {
	LineNumber() int
}

type Type interface {
	String() string
}
