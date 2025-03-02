package types

import (
	"testing"

	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
	testhelpers "github.com/chrispyles/slow/testing/helpers"
)

func TestRange_Construction(t *testing.T) {
	tests := []struct {
		name       string
		args       []execute.Value
		wantNewErr error
		want       []execute.Value
	}{
		{
			name: "common_float_type",
			args: []execute.Value{NewUint(10), NewInt(20), NewFloat(1.5)},
			want: []execute.Value{
				NewFloat(10),
				NewFloat(11.5),
				NewFloat(13),
				NewFloat(14.5),
				NewFloat(16),
				NewFloat(17.5),
				NewFloat(19),
			},
		},
		{
			name: "common_int_type",
			args: []execute.Value{NewUint(10), NewInt(20), NewInt(2)},
			want: []execute.Value{
				NewInt(10),
				NewInt(12),
				NewInt(14),
				NewInt(16),
				NewInt(18),
			},
		},
		{
			name: "common_uint_type",
			args: []execute.Value{NewUint(10), NewUint(20), NewUint(2)},
			want: []execute.Value{
				NewUint(10),
				NewUint(12),
				NewUint(14),
				NewUint(16),
				NewUint(18),
			},
		},
		{
			name: "common_bool_type",
			args: []execute.Value{NewBool(false), NewBool(true), NewBool(true)},
			want: []execute.Value{
				NewUint(0),
			},
		},
		{
			name:       "non_numeric_arg",
			args:       []execute.Value{NewStr("10"), NewInt(20), NewInt(2)},
			wantNewErr: errors.TypeErrorFromMessage("range cannot be called with non-numeric values: \"str\""),
		},
		{
			name: "0_to_5",
			args: []execute.Value{NewInt(0), NewInt(5), NewInt(1)},
			want: []execute.Value{
				NewInt(0),
				NewInt(1),
				NewInt(2),
				NewInt(3),
				NewInt(4),
			},
		},
		{
			name: "0_to_5_by_2",
			args: []execute.Value{NewInt(0), NewInt(5), NewInt(2)},
			want: []execute.Value{
				NewInt(0),
				NewInt(2),
				NewInt(4),
			},
		},
		{
			name: "0_to_6_by_2",
			args: []execute.Value{NewInt(0), NewInt(6), NewInt(2)},
			want: []execute.Value{
				NewInt(0),
				NewInt(2),
				NewInt(4),
			},
		},
		{
			name: "1_to_5",
			args: []execute.Value{NewInt(1), NewInt(5), NewInt(1)},
			want: []execute.Value{
				NewInt(1),
				NewInt(2),
				NewInt(3),
				NewInt(4),
			},
		},
		{
			name: "1_to_5_by_2",
			args: []execute.Value{NewInt(1), NewInt(5), NewInt(2)},
			want: []execute.Value{
				NewInt(1),
				NewInt(3),
			},
		},
		{
			name: "0_to_5_float",
			args: []execute.Value{NewFloat(0), NewFloat(5), NewFloat(1)},
			want: []execute.Value{
				NewFloat(0),
				NewFloat(1),
				NewFloat(2),
				NewFloat(3),
				NewFloat(4),
			},
		},
		{
			name: "0_to_5_by_1.5_float",
			args: []execute.Value{NewFloat(0), NewFloat(5), NewFloat(1.5)},
			want: []execute.Value{
				NewFloat(0),
				NewFloat(1.5),
				NewFloat(3),
				NewFloat(4.5),
			},
		},
		{
			name: "0_to_5_uint",
			args: []execute.Value{NewUint(0), NewUint(5), NewUint(1)},
			want: []execute.Value{
				NewUint(0),
				NewUint(1),
				NewUint(2),
				NewUint(3),
				NewUint(4),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewRangeGenerator(tc.args[0], tc.args[1], tc.args[2])
			testhelpers.CheckDiff(t, "NewRangeGenerator() error", tc.wantNewErr, err, allowUnexported)
			if tc.wantNewErr != nil {
				return
			}
			g := got.(*Generator)
			var gotVals []execute.Value
			for g.HasNext() {
				v, err := g.Next()
				if err != nil {
					t.Fatalf("g.Next() returned unexpected error: %v", err)
				}
				gotVals = append(gotVals, v)
			}
			testhelpers.CheckDiff(t, "NewRangeGenerator() values", tc.want, gotVals, allowUnexported)
		})
	}
}
