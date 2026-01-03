package ast

import (
	"testing"

	asttesting "github.com/chrispyles/slow/internal/ast/internal/testing"
	"github.com/chrispyles/slow/internal/execute"
	slowtesting "github.com/chrispyles/slow/internal/testing"
	"github.com/chrispyles/slow/internal/types"
)

func TestCallNode(t *testing.T) {
	for _, tc := range []asttesting.TestCase{
		{
			Name: "simple",
			Node: &CallNode{
				Func: &VariableNode{Name: "foo"},
				Args: []execute.Expression{&VariableNode{Name: "x"}},
			},
			Env: slowtesting.MustMakeEnv(t, map[string]execute.Value{
				"foo": types.NewGoFunc("foo", func(vs ...execute.Value) (execute.Value, error) {
					return vs[0], nil
				}),
				"x": types.NewInt(5),
			}),
			Want:        types.NewInt(5),
			WantSameEnv: true,
		},
		{
			Name: "object_method",
			Node: &CallNode{
				Func: &AttributeNode{Left: &VariableNode{Name: "foo"}, Right: "append"},
				Args: []execute.Expression{&VariableNode{Name: "x"}},
			},
			Env: slowtesting.MustMakeEnv(t, map[string]execute.Value{
				"foo": types.NewList(nil),
				"x":   types.NewInt(1),
			}),
			Want: types.Null,
			WantEnv: slowtesting.MustMakeEnv(t, map[string]execute.Value{
				"foo": types.NewList([]execute.Value{types.NewInt(1)}),
				"x":   types.NewInt(1),
			}),
		},
	} {
		asttesting.RunTestCase(t, tc)
	}
}
