package ast

import (
	"testing"

	asttesting "github.com/chrispyles/slow/src/ast/internal/testing"
)

func TestContinueNode(t *testing.T) {
	asttesting.RunTestCase(t, asttesting.TestCase{
		Node:    &ContinueNode{},
		WantErr: &continueError{},
	})
}
