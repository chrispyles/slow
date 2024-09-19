package builtins

import (
	"testing"

	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/types"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestBuiltins_range(t *testing.T) {

	doBuiltinTest(t, []builtinTest{
		{
			name: "one_arg",
			fn:   "range",
			args: []execute.Value{types.NewInt(10)},
			want: types.NewGenerator(&rangeGenerator{
				valueType: types.IntType,
				nextI:     0,
				startI:    0,
				stopI:     10,
				stepI:     1,
			}),
		},
		{
			name: "two_args",
			fn:   "range",
			args: []execute.Value{types.NewInt(10), types.NewInt(20)},
			want: types.NewGenerator(&rangeGenerator{
				valueType: types.IntType,
				nextI:     10,
				startI:    10,
				stopI:     20,
				stepI:     1,
			}),
		},
		{
			name: "three_args",
			fn:   "range",
			args: []execute.Value{types.NewInt(10), types.NewInt(20), types.NewInt(2)},
			want: types.NewGenerator(&rangeGenerator{
				valueType: types.IntType,
				nextI:     10,
				startI:    10,
				stopI:     20,
				stepI:     2,
			}),
		},
		{
			name:    "no_args",
			fn:      "range",
			args:    []execute.Value{},
			wantErr: errors.CallError("range", 0, 1),
		},
		{
			name:    "too_many_args",
			fn:      "range",
			args:    []execute.Value{types.NewInt(1), types.NewInt(1), types.NewInt(1), types.NewInt(1)},
			wantErr: errors.CallError("range", 4, 3),
		},
		{
			name: "common_float_type",
			fn:   "range",
			args: []execute.Value{types.NewUint(10), types.NewInt(20), types.NewFloat(1.2)},
			want: types.NewGenerator(&rangeGenerator{
				valueType: types.FloatType,
				nextF:     10,
				startF:    10,
				stopF:     20,
				stepF:     1.2,
			}),
		},
		{
			name: "common_int_type",
			fn:   "range",
			args: []execute.Value{types.NewUint(10), types.NewInt(20), types.NewInt(2)},
			want: types.NewGenerator(&rangeGenerator{
				valueType: types.IntType,
				nextI:     10,
				startI:    10,
				stopI:     20,
				stepI:     2,
			}),
		},
		{
			name: "common_uint_type",
			fn:   "range",
			args: []execute.Value{types.NewUint(10), types.NewUint(20), types.NewUint(2)},
			want: types.NewGenerator(&rangeGenerator{
				valueType: types.UintType,
				nextU:     10,
				startU:    10,
				stopU:     20,
				stepU:     2,
			}),
		},
		{
			name: "common_bool_type",
			fn:   "range",
			args: []execute.Value{types.NewBool(false), types.NewBool(true)},
			want: types.NewGenerator(&rangeGenerator{
				valueType: types.UintType,
				nextU:     0,
				startU:    0,
				stopU:     1,
				stepU:     1,
			}),
		},
		{
			name:    "non_numeric_arg",
			fn:      "range",
			args:    []execute.Value{types.NewStr("10")},
			wantErr: errors.TypeErrorFromMessage("range cannot be called with non-numeric values: \"str\""),
		},
		// TODO
	})
}

func TestRangeGenerator(t *testing.T) {
	tests := []struct {
		name string
		rg   *rangeGenerator
		want []execute.Value
	}{
		{
			name: "0_to_5",
			rg:   mustMakeRangeGenerator(t, types.NewInt(0), types.NewInt(5), types.NewInt(1)),
			want: []execute.Value{
				types.NewInt(0),
				types.NewInt(1),
				types.NewInt(2),
				types.NewInt(3),
				types.NewInt(4),
			},
		},
		{
			name: "0_to_5_by_2",
			rg:   mustMakeRangeGenerator(t, types.NewInt(0), types.NewInt(5), types.NewInt(2)),
			want: []execute.Value{
				types.NewInt(0),
				types.NewInt(2),
				types.NewInt(4),
			},
		},
		{
			name: "0_to_6_by_2",
			rg:   mustMakeRangeGenerator(t, types.NewInt(0), types.NewInt(6), types.NewInt(2)),
			want: []execute.Value{
				types.NewInt(0),
				types.NewInt(2),
				types.NewInt(4),
			},
		},
		{
			name: "1_to_5",
			rg:   mustMakeRangeGenerator(t, types.NewInt(1), types.NewInt(5), types.NewInt(1)),
			want: []execute.Value{
				types.NewInt(1),
				types.NewInt(2),
				types.NewInt(3),
				types.NewInt(4),
			},
		},
		{
			name: "1_to_5_by_2",
			rg:   mustMakeRangeGenerator(t, types.NewInt(1), types.NewInt(5), types.NewInt(2)),
			want: []execute.Value{
				types.NewInt(1),
				types.NewInt(3),
			},
		},
		{
			name: "0_to_5_float",
			rg:   mustMakeRangeGenerator(t, types.NewFloat(0), types.NewFloat(5), types.NewFloat(1)),
			want: []execute.Value{
				types.NewFloat(0),
				types.NewFloat(1),
				types.NewFloat(2),
				types.NewFloat(3),
				types.NewFloat(4),
			},
		},
		{
			name: "0_to_5_by_1.2_float",
			rg:   mustMakeRangeGenerator(t, types.NewFloat(0), types.NewFloat(5), types.NewFloat(1.2)),
			want: []execute.Value{
				types.NewFloat(0),
				types.NewFloat(1.2),
				types.NewFloat(2.4),
				types.NewFloat(3.6),
				types.NewFloat(4.8),
			},
		},
		{
			name: "0_to_5_uint",
			rg:   mustMakeRangeGenerator(t, types.NewUint(0), types.NewUint(5), types.NewUint(1)),
			want: []execute.Value{
				types.NewUint(0),
				types.NewUint(1),
				types.NewUint(2),
				types.NewUint(3),
				types.NewUint(4),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var got []execute.Value
			for tc.rg.HasNext() {
				n, err := tc.rg.Next()
				if err != nil {
					t.Fatalf("Next() returned unexpected error: %v", err)
				}
				got = append(got, n)
			}
			// TODO: the EquateApprox opt is necessary because 0_to_5_by_1.2_float yields 3.59999999... instead of 3.6; this should be fixed
			if diff := cmp.Diff(tc.want, got, allowUnexported, cmpopts.EquateApprox(1e-6, 0)); diff != "" {
				t.Errorf("range generator yielded incorrect items (-want +got):\n%s", diff)
			}
		})
	}
}

func mustMakeRangeGenerator(t *testing.T, start, stop, step execute.Value) *rangeGenerator {
	rg, err := newRangeGenerator(start, stop, step)
	if err != nil {
		t.Fatalf("newRangeGenerator() returned unexpected error: %v", err)
	}
	return rg
}
