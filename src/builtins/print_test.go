package builtins

import (
	"testing"

	"github.com/chrispyles/slow/src/execute"
	slowtesting "github.com/chrispyles/slow/src/testing"
	"github.com/chrispyles/slow/src/types"
)

func TestBuiltins_print(t *testing.T) {
	doBuiltinTest(t, []builtinTest{
		{
			name: "string",
			fn:   "print",
			args: []execute.Value{
				types.NewStr("foo"),
			},
			want:       types.Null,
			wantPrints: []string{"foo\n"},
		},
		{
			name: "value",
			fn:   "print",
			args: []execute.Value{
				&slowtesting.MockValue{StringRet: "MOCK_VALUE"},
			},
			want:       types.Null,
			wantPrints: []string{"MOCK_VALUE\n"},
		},
		{
			name: "many",
			fn:   "print",
			args: []execute.Value{
				&slowtesting.MockValue{StringRet: "MV1"},
				&slowtesting.MockValue{StringRet: "MV2"},
				&slowtesting.MockValue{StringRet: "MV3"},
			},
			want:       types.Null,
			wantPrints: []string{"MV1MV2MV3\n"},
		},
	})
}
