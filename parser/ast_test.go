package parser

import (
	"testing"

	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/operators"
	"github.com/chrispyles/slow/types"
	"github.com/google/go-cmp/cmp"
)

var allowTypesUnexported = cmp.AllowUnexported(
	operators.BinaryOperator{},
	operators.UnaryOperator{},
	types.Bool{},
	types.Float{},
	types.Func{},
	types.Int{},
	types.Iterator{},
	types.List{},
	types.Str{},
	types.Uint{},
)

// TODO: these tests also test the operator logic; is this OK?
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

func TestSmallASTs(t *testing.T) {
	tests := []struct {
		name string
		code string
		want *AST
	}{
		{
			name: "bin_op_combination",
			code: "x % 2 == 0u",
			want: &AST{
				Nodes: execute.Block{
					&BinaryOpNode{
						Op: operators.BinOp_EQ,
						Left: &BinaryOpNode{
							Op:    operators.BinOp_MOD,
							Left:  &VariableNode{Name: "x"},
							Right: &ConstantNode{Value: types.NewInt(2)},
						},
						Right: &ConstantNode{Value: types.NewUint(0)},
					},
				},
			},
		},
		{
			name: "arithmetic_with_method_calls",
			code: "m.get(i) + m.get(i + 1)",
			want: &AST{
				Nodes: execute.Block{
					&BinaryOpNode{
						Op: operators.BinOp_PLUS,
						Left: &CallNode{
							Func: &AttributeNode{
								Left:  &VariableNode{Name: "m"},
								Right: "get",
							},
							Args: []execute.Expression{
								&VariableNode{Name: "i"},
							},
						},
						Right: &CallNode{
							Func: &AttributeNode{
								Left:  &VariableNode{Name: "m"},
								Right: "get",
							},
							Args: []execute.Expression{
								&BinaryOpNode{
									Op:    operators.BinOp_PLUS,
									Left:  &VariableNode{Name: "i"},
									Right: &ConstantNode{Value: types.NewInt(1)},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewAST(tc.code)
			if err != nil {
				t.Fatalf("NewAST returned unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.want, got, allowTypesUnexported); diff != "" {
				t.Errorf("Incorrect AST (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCreateASTHailstone(t *testing.T) {
	code := `
func hailstone(x) {
	var l = [x]
	while x != 1 {
		if x % 2 == 0 {
			x //= 2
		}
		else {
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
								&VariableNode{Name: "x"},
							},
						},
					},
					&WhileNode{
						Cond: &BinaryOpNode{
							Op:    operators.BinOp_NEQ,
							Left:  &VariableNode{Name: "x"},
							Right: &ConstantNode{Value: types.NewInt(1)},
						},
						Body: execute.Block{
							&IfNode{
								Cond: &BinaryOpNode{
									Op: operators.BinOp_EQ,
									Left: &BinaryOpNode{
										Op:    operators.BinOp_MOD,
										Left:  &VariableNode{Name: "x"},
										Right: &ConstantNode{Value: types.NewInt(2)},
									},
									Right: &ConstantNode{Value: types.NewInt(0)},
								},
								Body: execute.Block{
									&BinaryOpNode{
										Op:    operators.BinOp_RFDIV,
										Left:  &VariableNode{Name: "x"},
										Right: &ConstantNode{Value: types.NewInt(2)},
									},
								},
								ElseBody: execute.Block{
									&AssignmentNode{
										Left: "x",
										Right: &BinaryOpNode{
											Op: operators.BinOp_PLUS,
											Left: &BinaryOpNode{
												Op:    operators.BinOp_TIMES,
												Left:  &ConstantNode{Value: types.NewInt(3)},
												Right: &VariableNode{Name: "x"},
											},
											Right: &ConstantNode{Value: types.NewInt(1)},
										},
									},
								},
							},
							&CallNode{
								Func: &AttributeNode{
									Left:  &VariableNode{Name: "l"},
									Right: "append",
								},
								Args: []execute.Expression{
									&VariableNode{Name: "x"},
								},
							},
						},
					},
					&ReturnNode{
						Value: &VariableNode{Name: "l"},
					},
				},
			},
			&ForNode{
				IterName: "x",
				Iter: &CallNode{
					Func: &VariableNode{Name: "range"},
					Args: []execute.Expression{
						&ConstantNode{Value: types.NewInt(1)},
						&ConstantNode{Value: types.NewInt(20)},
					},
				},
				Body: execute.Block{
					&CallNode{
						Func: &VariableNode{Name: "print"},
						Args: []execute.Expression{
							&CallNode{
								Func: &VariableNode{Name: "hailstone"},
								Args: []execute.Expression{
									&VariableNode{Name: "x"},
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

func TestCreateASTFibonacciMap(t *testing.T) {
	code := `
func fib(n) {
	var m = {0: 0, 1: 1}
	for i in range(2, n + 1) {
		var mi = m.get(i - 1) + m.get(i - 2)
		m.set(i, mi)
	}
	return m.get(n)
}

for x in range(20) {
	print(fib(x))
}	
`

	want := &AST{
		Nodes: execute.Block{
			&FuncNode{
				Name:     "fib",
				ArgNames: []string{"n"},
				Body: execute.Block{
					&VarNode{
						Name: "m",
						Value: &MapNode{
							Values: [][]execute.Expression{
								{
									&ConstantNode{Value: types.NewInt(0)},
									&ConstantNode{Value: types.NewInt(0)},
								},
								{
									&ConstantNode{Value: types.NewInt(1)},
									&ConstantNode{Value: types.NewInt(1)},
								},
							},
						},
					},
					&ForNode{
						IterName: "i",
						Iter: &CallNode{
							Func: &VariableNode{Name: "range"},
							Args: []execute.Expression{
								&ConstantNode{Value: types.NewInt(2)},
								&BinaryOpNode{
									Op:    operators.BinOp_PLUS,
									Left:  &VariableNode{Name: "n"},
									Right: &ConstantNode{Value: types.NewInt(1)},
								},
							},
						},
						Body: execute.Block{
							&VarNode{
								Name: "mi",
								Value: &BinaryOpNode{
									Op: operators.BinOp_PLUS,
									Left: &CallNode{
										Func: &AttributeNode{
											Left:  &VariableNode{Name: "m"},
											Right: "get",
										},
										Args: []execute.Expression{
											&BinaryOpNode{
												Op:    operators.BinOp_MINUS,
												Left:  &VariableNode{Name: "i"},
												Right: &ConstantNode{Value: types.NewInt(1)},
											},
										},
									},
									Right: &CallNode{
										Func: &AttributeNode{
											Left:  &VariableNode{Name: "m"},
											Right: "get",
										},
										Args: []execute.Expression{
											&BinaryOpNode{
												Op:    operators.BinOp_MINUS,
												Left:  &VariableNode{Name: "i"},
												Right: &ConstantNode{Value: types.NewInt(2)},
											},
										},
									},
								},
							},
							&CallNode{
								Func: &AttributeNode{
									Left:  &VariableNode{Name: "m"},
									Right: "set",
								},
								Args: []execute.Expression{
									&VariableNode{Name: "i"},
									&VariableNode{Name: "mi"},
								},
							},
						},
					},
					&ReturnNode{
						Value: &CallNode{
							Func: &AttributeNode{
								Left:  &VariableNode{Name: "m"},
								Right: "get",
							},
							Args: []execute.Expression{
								&VariableNode{Name: "n"},
							},
						},
					},
				},
			},
			&ForNode{
				IterName: "x",
				Iter: &CallNode{
					Func: &VariableNode{Name: "range"},
					Args: []execute.Expression{
						&ConstantNode{Value: types.NewInt(20)},
					},
				},
				Body: execute.Block{
					&CallNode{
						Func: &VariableNode{Name: "print"},
						Args: []execute.Expression{
							&CallNode{
								Func: &VariableNode{Name: "fib"},
								Args: []execute.Expression{
									&VariableNode{Name: "x"},
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
