package ast

import (
	"testing"

	asttesting "github.com/chrispyles/slow/src/ast/internal/testing"
	"github.com/chrispyles/slow/src/errors"
	"github.com/chrispyles/slow/src/execute"
	"github.com/chrispyles/slow/src/operators"
	slowtesting "github.com/chrispyles/slow/src/testing"
	"github.com/chrispyles/slow/src/types"
)

func TestBinaryOpNode(t *testing.T) {
	for _, tc := range []asttesting.TestCase{
		{
			Name: "simple",
			Node: &BinaryOpNode{
				Op:    operators.BinOp_PLUS,
				Left:  &ConstantNode{Value: types.NewInt(2)},
				Right: &ConstantNode{Value: types.NewInt(3)},
			},
			Want: types.NewInt(5),
		},
		{
			Name: "reassignment_operator_variable",
			Node: &BinaryOpNode{
				Op:    operators.BinOp_RPLUS,
				Left:  &VariableNode{Name: "foo"},
				Right: &ConstantNode{Value: types.NewInt(3)},
			},
			Env:     slowtesting.MustMakeEnv(t, map[string]execute.Value{"foo": types.NewInt(2)}),
			Want:    types.NewInt(5),
			WantEnv: slowtesting.MustMakeEnv(t, map[string]execute.Value{"foo": types.NewInt(5)}),
		},
		// TODO: once classes are available, add an attribute target test that doens't return an error
		{
			Name: "reassignment_operator_index",
			Node: &BinaryOpNode{
				Op: operators.BinOp_RPLUS,
				Left: &IndexNode{
					Container: &VariableNode{Name: "foo"},
					Index:     &ConstantNode{types.NewInt(0)},
				},
				Right: &ConstantNode{Value: types.NewInt(3)},
			},
			Env: slowtesting.MustMakeEnv(t, map[string]execute.Value{
				"foo": types.NewList([]execute.Value{types.NewInt(2)}),
			}),
			Want: types.NewInt(5),
			WantEnv: slowtesting.MustMakeEnv(t, map[string]execute.Value{
				"foo": types.NewList([]execute.Value{types.NewInt(5)}),
			}),
		},
		{
			Name: "reassignment_operator_non_declared_variable",
			Node: &BinaryOpNode{
				Op:    operators.BinOp_RPLUS,
				Left:  &VariableNode{Name: "bar"},
				Right: &ConstantNode{Value: types.NewInt(3)},
			},
			Env:         slowtesting.MustMakeEnv(t, map[string]execute.Value{"foo": types.NewInt(2)}),
			WantErr:     errors.NewNameError("bar"),
			WantSameEnv: true,
		},
	} {
		asttesting.RunTestCase(t, tc)
	}
}
