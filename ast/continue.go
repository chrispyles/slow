package ast

import (
	"github.com/chrispyles/slow/execute"
)

type continueError struct{}

func (*continueError) Error() string { return "" }

type ContinueNode struct{}

func (*ContinueNode) Execute(e *execute.Environment) (execute.Value, error) {
	return nil, &continueError{}
}
