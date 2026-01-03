package operators

import (
	"fmt"
	"testing"

	"github.com/chrispyles/slow/internal/errors"
	"github.com/chrispyles/slow/internal/execute"
	slowtesting "github.com/chrispyles/slow/internal/testing"
	slowcmpopts "github.com/chrispyles/slow/internal/testing/cmpopts"
	"github.com/chrispyles/slow/internal/types"
	"github.com/google/go-cmp/cmp"
)

func TestUnaryOperator_Value(t *testing.T) {
	type testCase struct {
		op      *UnaryOperator
		val     execute.Value
		want    execute.Value
		wantErr error
	}
	tests := []testCase{}

	op := UnOp_POS
	for _, v := range makeValues(1) {
		if v.Type() == types.StrType {
			tests = append(tests, testCase{op, v, nil, errors.IncompatibleType(v.Type(), op.String())})
			continue
		}
		tests = append(tests, testCase{op, v, v.CloneIfPrimitive(), nil})
	}

	// Add a test to check that UnOp_POS does nothing to negative operands.
	tests = append(tests, testCase{op, types.NewInt(-1), types.NewInt(-1), nil})

	op = UnOp_NEG
	for _, v := range makeValues(1) {
		if v.Type() == types.StrType {
			tests = append(tests, testCase{op, v, nil, errors.IncompatibleType(v.Type(), op.String())})
			continue
		}
		var want execute.Value = types.NewInt(-1)
		if v.Type() == types.FloatType {
			want = types.NewFloat(-1)
		}
		tests = append(tests, testCase{op, v, want, nil})
	}

	op = UnOp_NOT
	for _, v := range makeValues(2) {
		tests = append(tests, testCase{op, v, types.NewBool(false), nil})
	}
	for _, v := range makeValues(0) {
		if v.Type() == types.StrType {
			tests = append(tests, testCase{op, v, types.NewBool(false), nil})
			continue
		}
		tests = append(tests, testCase{op, v, types.NewBool(true), nil})
	}

	// Add UnOp_NOT tests for a non-numeric types.
	val := &slowtesting.MockValue{ToBoolRet: true}
	tests = append(tests, testCase{op, val, types.NewBool(false), nil})
	val = &slowtesting.MockValue{ToBoolRet: false}
	tests = append(tests, testCase{op, val, types.NewBool(true), nil})

	op = UnOp_INCR
	for _, v := range makeValues(2) {
		if v.Type() == types.StrType {
			tests = append(tests, testCase{op, v, nil, errors.IncompatibleType(v.Type(), op.String())})
			continue
		}
		var want execute.Value
		if v.Type() == types.BoolType {
			want = types.NewUint(2)
		} else {
			var err error
			want, err = v.Type().New(types.NewUint(3))
			if err != nil {
				t.Fatalf("type %s New() returned an error: %v", v.Type(), err)
			}
		}
		tests = append(tests, testCase{op, v, want, nil})
	}

	op = UnOp_DECR
	for _, v := range makeValues(2) {
		if v.Type() == types.StrType {
			tests = append(tests, testCase{op, v, nil, errors.IncompatibleType(v.Type(), op.String())})
			continue
		}
		var want execute.Value
		if v.Type() == types.BoolType {
			want = types.NewUint(0)
		} else {
			var err error
			want, err = v.Type().New(types.NewUint(1))
			if err != nil {
				t.Fatalf("type %s New() returned an error: %v", v.Type(), err)
			}
		}
		tests = append(tests, testCase{op, v, want, nil})
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s_%s_%s", tc.op, tc.val.Type(), tc.val), func(t *testing.T) {
			got, err := tc.op.Value(tc.val)
			if (tc.wantErr != nil) != (err != nil) {
				if diff := cmp.Diff(tc.wantErr, err); diff != "" {
					t.Errorf("Value() returned incorrect error (-want +got):\n%s", diff)
				}
			}
			if diff := cmp.Diff(tc.want, got, slowcmpopts.AllowUnexported()); diff != "" {
				t.Errorf("Value() returned diff (-want +got):\n%s", diff)
			}
			// if wantC, ok := tc.want.(*slowtesting.MockValue); ok {
			// 	gotC, ok := got.(*slowtesting.MockValue)
			// 	if !ok {
			// 		t.Fatalf("Value() did not return a MockValue: %v", got)
			// 	}
			// 	if gotC != wantC {
			// 		t.Errorf("Value() returned wrong MockValue: %v", gotC)
			// 	}
			// 	if got, want := gotC.CloneIfPrimitiveCalls, 1; got != want {
			// 		t.Errorf("Value() called MockValue.CloneIfPrimitive() %d times, want %d", got, want)
			// 	}
			// }
		})
	}
}

func TestUnaryOperator_IsReassignmentOperator(t *testing.T) {
	reassignmentOperators := map[string]bool{
		"++": true,
		"--": true,
	}
	for _, op := range unaryOperators {
		t.Run(op.String(), func(t *testing.T) {
			if got, want := op.IsReassignmentOperator(), reassignmentOperators[op.String()]; got != want {
				t.Errorf("IsReassignmentOperator() = %v, want %v", got, want)
			}
		})
	}
}

func TestUnaryOperator_String(t *testing.T) {
	for opS, op := range unaryOperators {
		t.Run(op.String(), func(t *testing.T) {
			if got, want := op.String(), opS; got != want {
				t.Errorf("String() = %q, want %q", got, want)
			}
		})
	}
}

// func makeValues(v uint64) []execute.Value {
// 	return []execute.Value{
// 		types.NewFloat(float64(v)),
// 		types.NewInt(int64(v)),
// 		types.NewUint(v),
// 		types.NewBool(v > 0),
// 		types.NewStr(fmt.Sprintf("%d", v)),
// 	}
// }
