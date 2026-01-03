package ast

import (
	"testing"

	asttesting "github.com/chrispyles/slow/internal/ast/internal/testing"
	"github.com/chrispyles/slow/internal/errors"
	"github.com/chrispyles/slow/internal/execute"
	slowtesting "github.com/chrispyles/slow/internal/testing"
	"github.com/chrispyles/slow/internal/types"
)

func TestAssignmentNode(t *testing.T) {
	for _, tc := range []asttesting.TestCase{
		{
			Name: "variable_target_not_yet_assigned",
			Node: &AssignmentNode{
				Left:  AssignmentTarget{Variable: "foo"},
				Right: &ConstantNode{Value: types.NewUint(2)},
			},
			Env:     slowtesting.MustMakeEnv(t, map[string]execute.Value{"foo": nil}),
			Want:    types.NewUint(2),
			WantEnv: slowtesting.MustMakeEnv(t, map[string]execute.Value{"foo": types.NewUint(2)}),
		},
		{
			Name: "variable_target_reassign",
			Node: &AssignmentNode{
				Left:  AssignmentTarget{Variable: "foo"},
				Right: &ConstantNode{Value: types.NewUint(2)},
			},
			Env:     slowtesting.MustMakeEnv(t, map[string]execute.Value{"foo": types.NewInt(1)}),
			Want:    types.NewUint(2),
			WantEnv: slowtesting.MustMakeEnv(t, map[string]execute.Value{"foo": types.NewUint(2)}),
		},
		// TODO: once classes are available, add an attribute target test that doens't return an error
		{
			Name: "attribute_target_failure",
			Node: &AssignmentNode{
				Left:  AssignmentTarget{Attribute: &AttributeNode{Left: &VariableNode{Name: "foo"}, Right: "bar"}},
				Right: &ConstantNode{Value: types.NewUint(2)},
			},
			Env:         slowtesting.MustMakeEnv(t, map[string]execute.Value{"foo": types.NewInt(1)}),
			WantErr:     errors.NewAttributeError(types.IntType, "bar"),
			WantSameEnv: true,
		},
		{
			Name: "index_target",
			Node: &AssignmentNode{
				Left: AssignmentTarget{Index: &IndexNode{
					Container: &VariableNode{Name: "foo"},
					Index:     &ConstantNode{Value: types.NewInt(1)},
				}},
				Right: &ConstantNode{Value: types.NewUint(0)},
			},
			Env: slowtesting.MustMakeEnv(t, map[string]execute.Value{
				"foo": types.NewList([]execute.Value{types.NewInt(2), types.NewInt(3)}),
			}),
			Want: types.NewList([]execute.Value{types.NewInt(2), types.NewUint(0)}),
			WantEnv: slowtesting.MustMakeEnv(t, map[string]execute.Value{
				"foo": types.NewList([]execute.Value{types.NewInt(2), types.NewUint(0)}),
			}),
		},
		{
			Name: "index_target_failure",
			Node: &AssignmentNode{
				Left: AssignmentTarget{Index: &IndexNode{
					Container: &VariableNode{Name: "foo"},
					Index:     &ConstantNode{Value: types.NewStr("bar")},
				}},
				Right: &ConstantNode{Value: types.NewUint(0)},
			},
			Env: slowtesting.MustMakeEnv(t, map[string]execute.Value{
				"foo": types.NewList([]execute.Value{types.NewInt(2), types.NewInt(3)}),
			}),
			WantErr:     errors.NonNumericIndexError(types.StrType, types.ListType),
			WantSameEnv: true,
		},
	} {
		asttesting.RunTestCase(t, tc)
	}
}
