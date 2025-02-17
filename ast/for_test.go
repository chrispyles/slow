package ast

import (
	"testing"

	asttesting "github.com/chrispyles/slow/ast/internal/testing"
	"github.com/chrispyles/slow/execute"
	slowtesting "github.com/chrispyles/slow/testing"
	"github.com/chrispyles/slow/types"
)

func TestForNode(t *testing.T) {
	for _, tc := range []asttesting.TestCase{
		{
			Name: "simple",
			Node: &ForNode{
				IterName: "i",
				Iter: &ConstantNode{
					Value: types.NewList([]execute.Value{types.NewInt(1), types.NewInt(2), types.NewInt(3)}),
				},
				Body: execute.Block{
					&CallNode{
						Func: &AttributeNode{
							Left:  &VariableNode{Name: "l"},
							Right: "append",
						},
						Args: []execute.Expression{&VariableNode{Name: "i"}},
					},
				},
			},
			Env: slowtesting.MustMakeEnv(t, map[string]execute.Value{
				"l": types.NewList(nil),
			}),
			Want: types.Null,
			WantEnv: slowtesting.MustMakeEnv(t, map[string]execute.Value{
				"l": types.NewList([]execute.Value{types.NewInt(1), types.NewInt(2), types.NewInt(3)}),
			}),
			// TODO: test break, continue, and error in loop body
		},
	} {
		asttesting.RunTestCase(t, tc)
	}
}
