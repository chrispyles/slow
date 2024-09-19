package builtins

import (
	"testing"

	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/types"
)

func TestBuiltins_type(t *testing.T) {
	doBuiltinTest(t, []builtinTest{
		{
			name: "bool",
			fn:   "type",
			args: []execute.Value{types.NewBool(true)},
			want: types.NewStr("bool"),
		},
		{
			name: "float",
			fn:   "type",
			args: []execute.Value{types.NewFloat(1)},
			want: types.NewStr("float"),
		},
		{
			name: "func",
			fn:   "type",
			args: []execute.Value{types.NewFunc("", nil, nil)},
			want: types.NewStr("func"),
		},
		{
			name: "generator",
			fn:   "type",
			args: []execute.Value{types.NewGenerator(nil)},
			want: types.NewStr("generator"),
		},
		{
			name: "int",
			fn:   "type",
			args: []execute.Value{types.NewInt(1)},
			want: types.NewStr("int"),
		},
		{
			name: "list",
			fn:   "type",
			args: []execute.Value{types.NewList(nil)},
			want: types.NewStr("list"),
		},
		{
			name: "map",
			fn:   "type",
			args: []execute.Value{types.NewMap()},
			want: types.NewStr("map"),
		},
		{
			name: "null",
			fn:   "type",
			args: []execute.Value{types.Null},
			want: types.NewStr("null"),
		},
		{
			name: "str",
			fn:   "type",
			args: []execute.Value{types.NewStr("")},
			want: types.NewStr("str"),
		},
		{
			name: "uint",
			fn:   "type",
			args: []execute.Value{types.NewUint(1)},
			want: types.NewStr("uint"),
		},
		{
			name:    "too_many_args",
			fn:      "type",
			args:    []execute.Value{types.NewInt(0), types.NewInt(1)},
			wantErr: errors.CallError("type", 2, 1),
		},
	})
}
