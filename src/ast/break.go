package ast

import (
	"github.com/chrispyles/slow/src/execute"
)

type breakError struct{}

func (*breakError) Error() string { return "" }

type BreakNode struct{}

func (*BreakNode) Execute(e *execute.Environment) (execute.Value, error) {
	return nil, &breakError{}
}
