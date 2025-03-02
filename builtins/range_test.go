package builtins

import (
	"testing"

	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/types"
)

func TestBuiltins_range(t *testing.T) {
	doBuiltinTest(t, []builtinTest{
		{
			name: "one_arg",
			fn:   "range",
			args: []execute.Value{types.NewInt(10)},
			want: mustMakeRangeGenerator(
				t,
				types.NewInt(0),
				types.NewInt(10),
				types.NewInt(1),
			),
		},
		{
			name: "two_args",
			fn:   "range",
			args: []execute.Value{types.NewInt(10), types.NewInt(20)},
			want: mustMakeRangeGenerator(
				t,
				types.NewInt(10),
				types.NewInt(20),
				types.NewInt(1),
			),
		},
		{
			name: "three_args",
			fn:   "range",
			args: []execute.Value{types.NewInt(10), types.NewInt(20), types.NewInt(2)},
			want: mustMakeRangeGenerator(
				t,
				types.NewInt(10),
				types.NewInt(20),
				types.NewInt(2),
			),
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
	})
}

func mustMakeRangeGenerator(t *testing.T, start, stop, step execute.Value) execute.Value {
	rg, err := types.NewRangeGenerator(start, stop, step)
	if err != nil {
		t.Fatalf("newRangeGenerator() returned unexpected error: %v", err)
	}
	return rg
}
