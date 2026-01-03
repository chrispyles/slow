package operators

import (
	"fmt"
	"testing"

	"github.com/chrispyles/slow/internal/execute"
	slowcmpopts "github.com/chrispyles/slow/internal/testing/cmpopts"
	"github.com/chrispyles/slow/internal/types"
	"github.com/google/go-cmp/cmp"
)

func TestTypeCaster(t *testing.T) {
	tests := []struct {
		left                execute.Value
		right               execute.Value
		wantContructorNotOk bool
		wantDest            execute.Type
		wantLeftCast        execute.Value
		wantRightCast       execute.Value
	}{
		{
			left:          types.NewFloat(1),
			right:         types.NewFloat(2),
			wantDest:      types.FloatType,
			wantLeftCast:  types.NewFloat(1),
			wantRightCast: types.NewFloat(2),
		},
		{
			left:          types.NewInt(1),
			right:         types.NewFloat(2),
			wantDest:      types.FloatType,
			wantLeftCast:  types.NewFloat(1),
			wantRightCast: types.NewFloat(2),
		},
		{
			left:          types.NewInt(1),
			right:         types.NewUint(2),
			wantDest:      types.IntType,
			wantLeftCast:  types.NewInt(1),
			wantRightCast: types.NewInt(2),
		},
		{
			left:          types.NewBool(true),
			right:         types.NewUint(2),
			wantDest:      types.UintType,
			wantLeftCast:  types.NewUint(1),
			wantRightCast: types.NewUint(2),
		},
		{
			left:          types.NewBool(true),
			right:         types.NewFloat(2),
			wantDest:      types.FloatType,
			wantLeftCast:  types.NewFloat(1),
			wantRightCast: types.NewFloat(2),
		},
		{
			left:          types.NewBool(true),
			right:         types.NewBool(false),
			wantDest:      types.BoolType,
			wantLeftCast:  types.NewBool(true),
			wantRightCast: types.NewBool(false),
		},
		// N.B. Even though str is not a numeric type, newTypeCaster should always succeed in returning
		// a typeCaster if the two input types are the same.
		{
			left:          types.NewStr("foo"),
			right:         types.NewStr("bar"),
			wantDest:      types.StrType,
			wantLeftCast:  types.NewStr("foo"),
			wantRightCast: types.NewStr("bar"),
		},
		// N.B. Because str is not a numeric type and both types are not the same, this should not be
		// able to construct a typeCaster.
		{
			left:                types.NewStr("foo"),
			right:               types.NewInt(1),
			wantContructorNotOk: true,
		},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s_%s_%s_%s", tc.left.Type(), tc.left, tc.right.Type(), tc.right), func(t *testing.T) {
			caster, ok := newTypeCaster(tc.left.Type(), tc.right.Type())
			if got, want := ok, !tc.wantContructorNotOk; got != want {
				t.Fatalf("newTypeCaster()[1] = %v, want %v", got, want)
			}
			if !ok {
				return
			}
			gotL, gotR := caster.Cast(tc.left, tc.right)
			if diff := cmp.Diff(tc.wantLeftCast, gotL, slowcmpopts.AllowUnexported()); diff != "" {
				t.Errorf("Cast() returned incorrect left value (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.wantRightCast, gotR, slowcmpopts.AllowUnexported()); diff != "" {
				t.Errorf("Cast() returned incorrect right value (-want +got):\n%s", diff)
			}
		})
	}
}
