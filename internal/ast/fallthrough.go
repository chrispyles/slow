package ast

import (
	"github.com/chrispyles/slow/internal/execute"
)

type fallthroughError struct{}

func (*fallthroughError) Error() string { return "" }

type FallthroughNode struct{}

func (*FallthroughNode) Execute(e *execute.Environment) (execute.Value, error) {
	return nil, &fallthroughError{}
}
