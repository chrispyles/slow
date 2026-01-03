package parser

import (
	"testing"

	"github.com/chrispyles/slow/ast"
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/operators"
	"github.com/chrispyles/slow/types"
	"github.com/google/go-cmp/cmp"
)

var allowTypesUnexported = cmp.AllowUnexported(
	ast.AssignmentTarget{},
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
		// {
		// 	code: "3 + 2",
		// 	want: types.NewInt(5),
		// },
		// {
		// 	code: "3+ 2",
		// 	want: types.NewInt(5),
		// },
		// {
		// 	code: "3 +2",
		// 	want: types.NewInt(5),
		// },
		// {
		// 	code: "3+2",
		// 	want: types.NewInt(5),
		// },
		// {
		// 	code: "3 - 2",
		// 	want: types.NewInt(1),
		// },
		// {
		// 	code: "3- 2",
		// 	want: types.NewInt(1),
		// },
		// {
		// 	code: "3 -2",
		// 	want: types.NewInt(1),
		// },
		// {
		// 	code: "3-2",
		// 	want: types.NewInt(1),
		// },
		// {
		// 	code: "2 * 10",
		// 	want: types.NewInt(20),
		// },
		// {
		// 	code: "2* 10",
		// 	want: types.NewInt(20),
		// },
		// {
		// 	code: "2 *10",
		// 	want: types.NewInt(20),
		// },
		// {
		// 	code: "10 / 2",
		// 	want: types.NewFloat(5),
		// },
		// {
		// 	code: "10/ 2",
		// 	want: types.NewFloat(5),
		// },
		// {
		// 	code: "10 /2",
		// 	want: types.NewFloat(5),
		// },
		// {
		// 	code: "10/2",
		// 	want: types.NewFloat(5),
		// },
		// {
		// 	code: "2*(-3+1)",
		// 	want: types.NewInt(-4),
		// },
		{
			code: "- ( 2 * 3)",
			want: types.NewInt(-6),
		},
		// TODO: more operators
		// TODO: test errors
	}
	for _, tc := range tests {
		t.Run(tc.code, func(t *testing.T) {
			a, err := Parse(tc.code)
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
		want *ast.AST
	}{
		{
			name: "bin_op_combination",
			code: "x % 2 == 0u",
			want: &ast.AST{
				Nodes: execute.Block{
					&ast.BinaryOpNode{
						Op: operators.BinOp_EQ,
						Left: &ast.BinaryOpNode{
							Op:    operators.BinOp_MOD,
							Left:  &ast.VariableNode{Name: "x"},
							Right: &ast.ConstantNode{Value: types.NewInt(2)},
						},
						Right: &ast.ConstantNode{Value: types.NewUint(0)},
					},
				},
			},
		},
		{
			name: "arithmetic_with_method_calls",
			code: "m.get(i) + m.get(i + 1)",
			want: &ast.AST{
				Nodes: execute.Block{
					&ast.BinaryOpNode{
						Op: operators.BinOp_PLUS,
						Left: &ast.CallNode{
							Func: &ast.AttributeNode{
								Left:  &ast.VariableNode{Name: "m"},
								Right: "get",
							},
							Args: []execute.Expression{
								&ast.VariableNode{Name: "i"},
							},
						},
						Right: &ast.CallNode{
							Func: &ast.AttributeNode{
								Left:  &ast.VariableNode{Name: "m"},
								Right: "get",
							},
							Args: []execute.Expression{
								&ast.BinaryOpNode{
									Op:    operators.BinOp_PLUS,
									Left:  &ast.VariableNode{Name: "i"},
									Right: &ast.ConstantNode{Value: types.NewInt(1)},
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
			got, err := Parse(tc.code)
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

	want := &ast.AST{
		Nodes: execute.Block{
			&ast.FuncNode{
				Name:     "hailstone",
				ArgNames: []string{"x"},
				Body: execute.Block{
					&ast.VarNode{
						Name: "l",
						Value: &ast.ListNode{
							Values: []execute.Expression{
								&ast.VariableNode{Name: "x"},
							},
						},
					},
					&ast.WhileNode{
						Cond: &ast.BinaryOpNode{
							Op:    operators.BinOp_NEQ,
							Left:  &ast.VariableNode{Name: "x"},
							Right: &ast.ConstantNode{Value: types.NewInt(1)},
						},
						Body: execute.Block{
							&ast.IfNode{
								Cond: &ast.BinaryOpNode{
									Op: operators.BinOp_EQ,
									Left: &ast.BinaryOpNode{
										Op:    operators.BinOp_MOD,
										Left:  &ast.VariableNode{Name: "x"},
										Right: &ast.ConstantNode{Value: types.NewInt(2)},
									},
									Right: &ast.ConstantNode{Value: types.NewInt(0)},
								},
								Body: execute.Block{
									&ast.BinaryOpNode{
										Op:    operators.BinOp_RFDIV,
										Left:  &ast.VariableNode{Name: "x"},
										Right: &ast.ConstantNode{Value: types.NewInt(2)},
									},
								},
								ElseBody: execute.Block{
									&ast.AssignmentNode{
										Left: ast.AssignmentTarget{Variable: "x"},
										Right: &ast.BinaryOpNode{
											Op: operators.BinOp_PLUS,
											Left: &ast.BinaryOpNode{
												Op:    operators.BinOp_TIMES,
												Left:  &ast.ConstantNode{Value: types.NewInt(3)},
												Right: &ast.VariableNode{Name: "x"},
											},
											Right: &ast.ConstantNode{Value: types.NewInt(1)},
										},
									},
								},
							},
							&ast.CallNode{
								Func: &ast.AttributeNode{
									Left:  &ast.VariableNode{Name: "l"},
									Right: "append",
								},
								Args: []execute.Expression{
									&ast.VariableNode{Name: "x"},
								},
							},
						},
					},
					&ast.ReturnNode{
						Value: &ast.VariableNode{Name: "l"},
					},
				},
			},
			&ast.ForNode{
				IterName: "x",
				Iter: &ast.CallNode{
					Func: &ast.VariableNode{Name: "range"},
					Args: []execute.Expression{
						&ast.ConstantNode{Value: types.NewInt(1)},
						&ast.ConstantNode{Value: types.NewInt(20)},
					},
				},
				Body: execute.Block{
					&ast.CallNode{
						Func: &ast.VariableNode{Name: "print"},
						Args: []execute.Expression{
							&ast.CallNode{
								Func: &ast.VariableNode{Name: "hailstone"},
								Args: []execute.Expression{
									&ast.VariableNode{Name: "x"},
								},
							},
						},
					},
				},
			},
		},
	}

	got, err := Parse(code)
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

	want := &ast.AST{
		Nodes: execute.Block{
			&ast.FuncNode{
				Name:     "fib",
				ArgNames: []string{"n"},
				Body: execute.Block{
					&ast.VarNode{
						Name: "m",
						Value: &ast.MapNode{
							Values: [][]execute.Expression{
								{
									&ast.ConstantNode{Value: types.NewInt(0)},
									&ast.ConstantNode{Value: types.NewInt(0)},
								},
								{
									&ast.ConstantNode{Value: types.NewInt(1)},
									&ast.ConstantNode{Value: types.NewInt(1)},
								},
							},
						},
					},
					&ast.ForNode{
						IterName: "i",
						Iter: &ast.CallNode{
							Func: &ast.VariableNode{Name: "range"},
							Args: []execute.Expression{
								&ast.ConstantNode{Value: types.NewInt(2)},
								&ast.BinaryOpNode{
									Op:    operators.BinOp_PLUS,
									Left:  &ast.VariableNode{Name: "n"},
									Right: &ast.ConstantNode{Value: types.NewInt(1)},
								},
							},
						},
						Body: execute.Block{
							&ast.VarNode{
								Name: "mi",
								Value: &ast.BinaryOpNode{
									Op: operators.BinOp_PLUS,
									Left: &ast.CallNode{
										Func: &ast.AttributeNode{
											Left:  &ast.VariableNode{Name: "m"},
											Right: "get",
										},
										Args: []execute.Expression{
											&ast.BinaryOpNode{
												Op:    operators.BinOp_MINUS,
												Left:  &ast.VariableNode{Name: "i"},
												Right: &ast.ConstantNode{Value: types.NewInt(1)},
											},
										},
									},
									Right: &ast.CallNode{
										Func: &ast.AttributeNode{
											Left:  &ast.VariableNode{Name: "m"},
											Right: "get",
										},
										Args: []execute.Expression{
											&ast.BinaryOpNode{
												Op:    operators.BinOp_MINUS,
												Left:  &ast.VariableNode{Name: "i"},
												Right: &ast.ConstantNode{Value: types.NewInt(2)},
											},
										},
									},
								},
							},
							&ast.CallNode{
								Func: &ast.AttributeNode{
									Left:  &ast.VariableNode{Name: "m"},
									Right: "set",
								},
								Args: []execute.Expression{
									&ast.VariableNode{Name: "i"},
									&ast.VariableNode{Name: "mi"},
								},
							},
						},
					},
					&ast.ReturnNode{
						Value: &ast.CallNode{
							Func: &ast.AttributeNode{
								Left:  &ast.VariableNode{Name: "m"},
								Right: "get",
							},
							Args: []execute.Expression{
								&ast.VariableNode{Name: "n"},
							},
						},
					},
				},
			},
			&ast.ForNode{
				IterName: "x",
				Iter: &ast.CallNode{
					Func: &ast.VariableNode{Name: "range"},
					Args: []execute.Expression{
						&ast.ConstantNode{Value: types.NewInt(20)},
					},
				},
				Body: execute.Block{
					&ast.CallNode{
						Func: &ast.VariableNode{Name: "print"},
						Args: []execute.Expression{
							&ast.CallNode{
								Func: &ast.VariableNode{Name: "fib"},
								Args: []execute.Expression{
									&ast.VariableNode{Name: "x"},
								},
							},
						},
					},
				},
			},
		},
	}

	got, err := Parse(code)
	if err != nil {
		t.Fatalf("NewAST returned unexpected error: %v", err)
	}

	if diff := cmp.Diff(want, got, allowTypesUnexported); diff != "" {
		t.Errorf("Incorrect AST (-want +got):\n%s", diff)
	}
}
