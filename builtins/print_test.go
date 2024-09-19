package builtins

import (
	"testing"

	"github.com/chrispyles/slow/execute"
	slowtesting "github.com/chrispyles/slow/testing"
	"github.com/chrispyles/slow/types"
)

func TestBuiltins_print(t *testing.T) {
	doBuiltinTest(t, []builtinTest{
		{
			name: "string",
			fn:   "print",
			args: []execute.Value{
				types.NewStr("foo"),
			},
			want:         types.Null,
			wantPrintlns: []string{"foo"},
		},
		{
			name: "value",
			fn:   "print",
			args: []execute.Value{
				&slowtesting.MockValue{StringRet: "MOCK_VALUE"},
			},
			want:         types.Null,
			wantPrintlns: []string{"MOCK_VALUE"},
		},
		{
			name: "many",
			fn:   "print",
			args: []execute.Value{
				&slowtesting.MockValue{StringRet: "MV1"},
				&slowtesting.MockValue{StringRet: "MV2"},
				&slowtesting.MockValue{StringRet: "MV3"},
			},
			want:         types.Null,
			wantPrintlns: []string{"MV1MV2MV3"},
		},
	})
}
