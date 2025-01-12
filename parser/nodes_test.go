package parser

import (
	"testing"

	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/operators"
	slowtesting "github.com/chrispyles/slow/testing"
	"github.com/chrispyles/slow/types"
	"github.com/google/go-cmp/cmp"
)

type testCase struct {
	name        string
	node        execute.Expression
	env         *execute.Environment
	want        execute.Value
	wantEnv     *execute.Environment
	wantSameEnv bool
	wantErr     error
}

func runTestCase(t *testing.T, tc testCase) {
	t.Run(tc.name, func(t *testing.T) {
		wantEnv := tc.wantEnv
		if tc.wantSameEnv {
			wantEnv = tc.env.Copy()
		}
		got, err := tc.node.Execute(tc.env)
		if diff := cmp.Diff(tc.wantErr, err, slowtesting.AllowUnexported()); diff != "" {
			t.Errorf("Execute() returned incorrect error (-want +got):\n%s", diff)
		}
		if diff := cmp.Diff(tc.want, got, slowtesting.AllowUnexported()); diff != "" {
			t.Errorf("Execute() returned unexpected diff (-want +got):\n%s", diff)
		}
		if diff := cmp.Diff(wantEnv, tc.env, slowtesting.AllowUnexported(), slowtesting.EquateFuncs()); diff != "" {
			t.Errorf("env after Execute() has unexpected diff (-want +got):\n%s", diff)
		}
	})
}

func TestAssignmentNode(t *testing.T) {
	for _, tc := range []testCase{
		{
			name: "variable_target_not_yet_assigned",
			node: &AssignmentNode{
				Left:  assignmentTarget{variable: "foo"},
				Right: &ConstantNode{Value: types.NewUint(2)},
			},
			env:     slowtesting.MustMakeEnv(t, map[string]execute.Value{"foo": nil}),
			want:    types.NewUint(2),
			wantEnv: slowtesting.MustMakeEnv(t, map[string]execute.Value{"foo": types.NewUint(2)}),
		},
		{
			name: "variable_target_reassign",
			node: &AssignmentNode{
				Left:  assignmentTarget{variable: "foo"},
				Right: &ConstantNode{Value: types.NewUint(2)},
			},
			env:     slowtesting.MustMakeEnv(t, map[string]execute.Value{"foo": types.NewInt(1)}),
			want:    types.NewUint(2),
			wantEnv: slowtesting.MustMakeEnv(t, map[string]execute.Value{"foo": types.NewUint(2)}),
		},
		// TODO: once classes are available, add an attribute target test that doens't return an error
		{
			name: "attribute_target_failure",
			node: &AssignmentNode{
				Left:  assignmentTarget{attribute: &AttributeNode{Left: &VariableNode{Name: "foo"}, Right: "bar"}},
				Right: &ConstantNode{Value: types.NewUint(2)},
			},
			env:         slowtesting.MustMakeEnv(t, map[string]execute.Value{"foo": types.NewInt(1)}),
			wantErr:     errors.NewAttributeError(types.IntType, "bar"),
			wantSameEnv: true,
		},
		{
			name: "index_target",
			node: &AssignmentNode{
				Left: assignmentTarget{index: &IndexNode{
					Container: &VariableNode{Name: "foo"},
					Index:     &ConstantNode{Value: types.NewInt(1)},
				}},
				Right: &ConstantNode{Value: types.NewUint(0)},
			},
			env: slowtesting.MustMakeEnv(t, map[string]execute.Value{
				"foo": types.NewList([]execute.Value{types.NewInt(2), types.NewInt(3)}),
			}),
			want: types.NewList([]execute.Value{types.NewInt(2), types.NewUint(0)}),
			wantEnv: slowtesting.MustMakeEnv(t, map[string]execute.Value{
				"foo": types.NewList([]execute.Value{types.NewInt(2), types.NewUint(0)}),
			}),
		},
		{
			name: "index_target_failure",
			node: &AssignmentNode{
				Left: assignmentTarget{index: &IndexNode{
					Container: &VariableNode{Name: "foo"},
					Index:     &ConstantNode{Value: types.NewStr("bar")},
				}},
				Right: &ConstantNode{Value: types.NewUint(0)},
			},
			env: slowtesting.MustMakeEnv(t, map[string]execute.Value{
				"foo": types.NewList([]execute.Value{types.NewInt(2), types.NewInt(3)}),
			}),
			wantErr:     errors.NonNumericIndexError(types.StrType, types.ListType),
			wantSameEnv: true,
		},
	} {
		runTestCase(t, tc)
	}
}

func TestAttributeNode(t *testing.T) {
	mv := &slowtesting.MockValue{Attributes: map[string]execute.Value{"foo": types.NewInt(1)}}
	runTestCase(t, testCase{
		name:        "success",
		node:        &AttributeNode{Left: &VariableNode{Name: "mockValue"}, Right: "foo"},
		env:         slowtesting.MustMakeEnv(t, map[string]execute.Value{"mockValue": mv}),
		want:        types.NewInt(1),
		wantSameEnv: true,
	})
	_, wantErr := mv.GetAttribute("bar")
	if wantErr == nil {
		t.Fatalf("GetAttribute() for nonexistent attribute returned nil error")
	}
	runTestCase(t, testCase{
		name:        "error",
		node:        &AttributeNode{Left: &VariableNode{Name: "mockValue"}, Right: "bar"},
		env:         slowtesting.MustMakeEnv(t, map[string]execute.Value{"mockValue": mv}),
		wantErr:     wantErr,
		wantSameEnv: true,
	})
}

func TestBinaryOpNode(t *testing.T) {
	for _, tc := range []testCase{
		{
			name: "simple",
			node: &BinaryOpNode{
				Op:    operators.BinOp_PLUS,
				Left:  &ConstantNode{Value: types.NewInt(2)},
				Right: &ConstantNode{Value: types.NewInt(3)},
			},
			want: types.NewInt(5),
		},
		{
			name: "reassignment_operator_variable",
			node: &BinaryOpNode{
				Op:    operators.BinOp_RPLUS,
				Left:  &VariableNode{Name: "foo"},
				Right: &ConstantNode{Value: types.NewInt(3)},
			},
			env:     slowtesting.MustMakeEnv(t, map[string]execute.Value{"foo": types.NewInt(2)}),
			want:    types.NewInt(5),
			wantEnv: slowtesting.MustMakeEnv(t, map[string]execute.Value{"foo": types.NewInt(5)}),
		},
		// TODO: once classes are available, add an attribute target test that doens't return an error
		{
			name: "reassignment_operator_index",
			node: &BinaryOpNode{
				Op: operators.BinOp_RPLUS,
				Left: &IndexNode{
					Container: &VariableNode{Name: "foo"},
					Index:     &ConstantNode{types.NewInt(0)},
				},
				Right: &ConstantNode{Value: types.NewInt(3)},
			},
			env: slowtesting.MustMakeEnv(t, map[string]execute.Value{
				"foo": types.NewList([]execute.Value{types.NewInt(2)}),
			}),
			want: types.NewInt(5),
			wantEnv: slowtesting.MustMakeEnv(t, map[string]execute.Value{
				"foo": types.NewList([]execute.Value{types.NewInt(5)}),
			}),
		},
		{
			name: "reassignment_operator_non_declared_variable",
			node: &BinaryOpNode{
				Op:    operators.BinOp_RPLUS,
				Left:  &VariableNode{Name: "bar"},
				Right: &ConstantNode{Value: types.NewInt(3)},
			},
			env:         slowtesting.MustMakeEnv(t, map[string]execute.Value{"foo": types.NewInt(2)}),
			wantErr:     errors.NewNameError("bar"),
			wantSameEnv: true,
		},
	} {
		runTestCase(t, tc)
	}
}

func TestBreakNode(t *testing.T) {
	runTestCase(t, testCase{
		node:    &BreakNode{},
		wantErr: &breakError{},
	})
}

func TestCallNode(t *testing.T) {
	for _, tc := range []testCase{
		{
			name: "simple",
			node: &CallNode{
				Func: &VariableNode{Name: "foo"},
				Args: []execute.Expression{&VariableNode{Name: "x"}},
			},
			env: slowtesting.MustMakeEnv(t, map[string]execute.Value{
				"foo": types.NewGoFunc("foo", func(vs ...execute.Value) (execute.Value, error) {
					return vs[0], nil
				}),
				"x": types.NewInt(5),
			}),
			want:        types.NewInt(5),
			wantSameEnv: true,
		},
		{
			name: "object_method",
			node: &CallNode{
				Func: &AttributeNode{Left: &VariableNode{Name: "foo"}, Right: "append"},
				Args: []execute.Expression{&VariableNode{Name: "x"}},
			},
			env: slowtesting.MustMakeEnv(t, map[string]execute.Value{
				"foo": types.NewList(nil),
				"x":   types.NewInt(1),
			}),
			want: types.Null,
			wantEnv: slowtesting.MustMakeEnv(t, map[string]execute.Value{
				"foo": types.NewList([]execute.Value{types.NewInt(1)}),
				"x":   types.NewInt(1),
			}),
		},
	} {
		runTestCase(t, tc)
	}
}

func TestConstantNode(t *testing.T) {
	mv := &slowtesting.MockValue{}
	runTestCase(t, testCase{
		node:        &ConstantNode{Value: mv},
		want:        mv,
		wantSameEnv: true,
	})
}

func TestContinueNode(t *testing.T) {
	runTestCase(t, testCase{
		node:    &ContinueNode{},
		wantErr: &continueError{},
	})
}

func TestFallthroughNode(t *testing.T) {
	runTestCase(t, testCase{
		node:    &FallthroughNode{},
		wantErr: &fallthroughError{},
	})

}

func TestForNode(t *testing.T) {
	for _, tc := range []testCase{
		{
			name: "simple",
			node: &ForNode{
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
			env: slowtesting.MustMakeEnv(t, map[string]execute.Value{
				"l": types.NewList(nil),
			}),
			want: types.Null,
			wantEnv: slowtesting.MustMakeEnv(t, map[string]execute.Value{
				"l": types.NewList([]execute.Value{types.NewInt(1), types.NewInt(2), types.NewInt(3)}),
			}),
			// TODO: test break, continue, and error in loop body
		},
	} {
		runTestCase(t, tc)
	}
}

func TestFuncNode(t *testing.T) {

}

func TestIfNode(t *testing.T) {

}

func TestIndexNode(t *testing.T) {

}

func TestListNode(t *testing.T) {

}

func TestMapNode(t *testing.T) {

}

func TestReturnNode(t *testing.T) {

}

func TestSwitchNode(t *testing.T) {

}

func TestUnaryOpNode(t *testing.T) {

}

func TestVarNode(t *testing.T) {

}

func TestVariableNode(t *testing.T) {

}

func TestWhileNode(t *testing.T) {

}
