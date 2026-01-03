package ast

import (
	"testing"

	asttesting "github.com/chrispyles/slow/src/ast/internal/testing"
	slowtesting "github.com/chrispyles/slow/src/testing"
)

func TestConstantNode(t *testing.T) {
	mv := &slowtesting.MockValue{}
	asttesting.RunTestCase(t, asttesting.TestCase{
		Node:        &ConstantNode{Value: mv},
		Want:        mv,
		WantSameEnv: true,
	})
}
