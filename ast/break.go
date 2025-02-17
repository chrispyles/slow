package ast

import (
	"github.com/chrispyles/slow/execute"
)

type breakError struct{}

func (*breakError) Error() string { return "" }

type BreakNode struct{}

func (*BreakNode) Execute(e *execute.Environment) (execute.Value, error) {
	return nil, &breakError{}
}
