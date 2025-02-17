package ast

import (
	"testing"

	asttesting "github.com/chrispyles/slow/ast/internal/testing"
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/types"
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
