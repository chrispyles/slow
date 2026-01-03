package builtins

import (
	"testing"

	"github.com/chrispyles/slow/src/errors"
	"github.com/chrispyles/slow/src/execute"
	slowtesting "github.com/chrispyles/slow/src/testing"
	"github.com/chrispyles/slow/src/types"
)

func TestBuiltins_len(t *testing.T) {
	doBuiltinTest(t, []builtinTest{
		{
			name: "success",
			fn:   "len",
			args: []execute.Value{&slowtesting.MockValue{LengthRet: 10}},
			want: types.NewUint(10),
		},
		{
			name: "err",
			fn:   "len",
			args: []execute.Value{&slowtesting.MockValue{
				LengthErr: errors.NoLengthError(slowtesting.NewMockType()),
			}},
			wantErr: errors.NoLengthError(slowtesting.NewMockType()),
		},
		{
			name:    "few_args",
			fn:      "len",
			args:    []execute.Value{},
			wantErr: errors.CallError("len", 0, 1),
		},
		{
			name: "many_args",
			fn:   "len",
			args: []execute.Value{
				&slowtesting.MockValue{},
				&slowtesting.MockValue{},
			},
			wantErr: errors.CallError("len", 2, 1),
		},
	})
}
