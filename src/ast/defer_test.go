package ast

import (
	"testing"

	asttesting "github.com/chrispyles/slow/src/ast/internal/testing"
	"github.com/chrispyles/slow/src/execute"
	"github.com/chrispyles/slow/src/types"
)

func TestDeferNode(t *testing.T) {
	asttesting.RunTestCase(t, asttesting.TestCase{
		Node: &DeferNode{
			Expr: &CallNode{
				Func: &VariableNode{Name: "print"},
				Args: []execute.Expression{
					&VariableNode{Name: "x"},
				},
			},
		},
		WantErr: &types.DeferError{
			Expr: &CallNode{
				Func: &VariableNode{Name: "print"},
				Args: []execute.Expression{
					&VariableNode{Name: "x"},
				},
			},
		},
	})
}
