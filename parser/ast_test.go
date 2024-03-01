package parser

import (
	"testing"

	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/operators"
	"github.com/chrispyles/slow/types"
	"github.com/google/go-cmp/cmp"
)

var allowTypesUnexported = cmp.AllowUnexported(
	types.Bool{},
	types.Float{},
	types.Func{},
	types.Int{},
	types.Iterator{},
	types.List{},
	types.Str{},
	types.Uint{},
)

func TestArithmetic(t *testing.T) {
	tests := []struct {
		code string
		want execute.Value
	}{
		{
			code: "3 + 2",
			want: types.NewInt(5),
		},
		{
			code: "3+ 2",
			want: types.NewInt(5),
		},
		{
			code: "3 +2",
			want: types.NewInt(5),
		},
		{
			code: "3+2",
			want: types.NewInt(5),
		},
		{
			code: "3 - 2",
			want: types.NewInt(1),
		},
		{
			code: "3- 2",
			want: types.NewInt(1),
		},
		{
			code: "3 -2",
			want: types.NewInt(1),
		},
		{
			code: "3-2",
			want: types.NewInt(1),
		},
		{
			code: "2 * 10",
			want: types.NewInt(20),
		},
		{
			code: "2* 10",
			want: types.NewInt(20),
		},
		{
			code: "2 *10",
			want: types.NewInt(20),
		},
		{
			code: "10 / 2",
			want: types.NewFloat(5),
		},
		{
			code: "10/ 2",
			want: types.NewFloat(5),
		},
		{
			code: "10 /2",
			want: types.NewFloat(5),
		},
		{
			code: "10/2",
			want: types.NewFloat(5),
		},
		{
			code: "2*(-3+1)",
			want: types.NewInt(-4),
		},
		{
			code: "- ( 2 * 3)",
			want: types.NewInt(-6),
		},
		// TODO: more operators
		// TODO: test errors
	}
	for _, tc := range tests {
		t.Run(tc.code, func(t *testing.T) {
			a, err := NewAST(tc.code)
			if err != nil {
				t.Fatalf("NewAST returned an unexpected error: %v", err)
			}
			got, err := a.Execute(execute.NewEnvironment())
			if err != nil {
				t.Fatalf("a.Execute returned an unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.want, got, cmp.AllowUnexported(types.Float{}, types.Int{}, types.Uint{})); diff != "" {
				t.Errorf("a.Execute returned unexpected diff (-want +got):\n%s", diff)
			}
		})
	}
}

// TODO: update code to exercise all node types
func TestCreateAST(t *testing.T) {
	code := `func hailstone(x) {
	var l = [x]
	while x != 1 {
		if x % 2 == 0 {
			x = x / 2 # TODO: x /= 2
		} else {
			x = 3 * x + 1
		}
		l.append(x)
	}
	return l
}

for x in range(1, 20) {
	print(hailstone(x))
}
`

	want := &AST{
		Nodes: execute.Block{
			&FuncNode{
				Name:     "hailstone",
				ArgNames: []string{"x"},
				Body: execute.Block{
					&VarNode{
						Name: "l",
						Value: &ListNode{
							Values: []execute.Expression{
								&VariableNode{
									Name: "x",
								},
							},
						},
					},
					&WhileNode{
						Cond: &BinaryOpNode{
							Op: operators.BinOp_NEQ,
							Left: &VariableNode{
								Name: "x",
							},
							Right: &ConstantNode{
								Value: types.NewInt(1),
							},
						},
						Body: execute.Block{
							&IfNode{
								Cond: &BinaryOpNode{
									Op: operators.BinOp_EQ,
									Left: &BinaryOpNode{
										Op: operators.BinOp_MOD,
										Left: &VariableNode{
											Name: "x",
										},
										Right: &ConstantNode{
											Value: types.NewInt(2),
										},
									},
									Right: &ConstantNode{
										Value: types.NewInt(0),
									},
								},
								Body: execute.Block{
									&AssignmentNode{
										Left: "x",
										Right: &BinaryOpNode{
											Op: operators.BinOp_DIV,
											Left: &VariableNode{
												Name: "x",
											},
											Right: &ConstantNode{
												Value: types.NewInt(2),
											},
										},
									},
								},
								ElseBody: execute.Block{
									&AssignmentNode{
										Left: "x",
										Right: &BinaryOpNode{
											Op: operators.BinOp_PLUS,
											Left: &BinaryOpNode{
												Op: operators.BinOp_TIMES,
												Left: &ConstantNode{
													Value: types.NewInt(3),
												},
												Right: &VariableNode{
													Name: "x",
												},
											},
											Right: &ConstantNode{
												Value: types.NewInt(1),
											},
										},
									},
								},
							},
							&CallNode{
								Func: &AttributeNode{
									Left: &VariableNode{
										Name: "l",
									},
									Right: "append",
								},
								Args: []execute.Expression{
									&VariableNode{
										Name: "x",
									},
								},
							},
						},
					},
					&ReturnNode{
						Value: &VariableNode{
							Name: "l",
						},
					},
				},
			},
			&ForNode{
				IterName: "x",
				Iter: &CallNode{
					Func: &VariableNode{
						Name: "range",
					},
					Args: []execute.Expression{
						&ConstantNode{
							Value: types.NewInt(1),
						},
						&ConstantNode{
							Value: types.NewInt(20),
						},
					},
				},
				Body: execute.Block{
					&CallNode{
						Func: &VariableNode{
							Name: "print",
						},
						Args: []execute.Expression{
							&CallNode{
								Func: &VariableNode{
									Name: "hailstone",
								},
								Args: []execute.Expression{
									&VariableNode{
										Name: "x",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	got, err := NewAST(code)
	if err != nil {
		t.Fatalf("NewAST returned unexpected error: %v", err)
	}

	if diff := cmp.Diff(want, got, allowTypesUnexported); diff != "" {
		t.Errorf("Incorrect AST (-want +got):\n%s", diff)
	}
}
