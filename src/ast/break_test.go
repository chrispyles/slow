package ast

import (
	"testing"

	asttesting "github.com/chrispyles/slow/src/ast/internal/testing"
)

func TestBreakNode(t *testing.T) {
	asttesting.RunTestCase(t, asttesting.TestCase{
		Node:    &BreakNode{},
		WantErr: &breakError{},
	})
}
