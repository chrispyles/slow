package ast

import (
	"testing"

	asttesting "github.com/chrispyles/slow/internal/ast/internal/testing"
)

func TestFallthroughNode(t *testing.T) {
	asttesting.RunTestCase(t, asttesting.TestCase{
		Node:    &FallthroughNode{},
		WantErr: &fallthroughError{},
	})
}
