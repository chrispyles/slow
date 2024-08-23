package builtins

import (
	"testing"

	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
	slowtesting "github.com/chrispyles/slow/testing"
	"github.com/chrispyles/slow/types"
	"github.com/google/go-cmp/cmp"
)

var allowUnexported = cmp.AllowUnexported(
	errors.AttributeError{},
	errors.DeclarationError{},
	errors.EOFError{},
	errors.IndexError{},
	errors.NameError{},
	errors.SyntaxError{},
	errors.TypeError{},

	types.Bool{},
	types.Float{},
	types.Func{},
	types.Int{},
	types.Iterator{},
	types.List{},
	types.Str{},
	types.Uint{},
)

func TestBuiltins(t *testing.T) {
	tests := []struct {
		name         string
		fn           string
		args         []execute.Value
		want         execute.Value
		wantPrintlns []string
		wantErr      error
	}{
		// TODO: exit
		{
			name: "len_success",
			fn:   "len",
			args: []execute.Value{&slowtesting.MockValue{LengthRet: 10}},
			want: types.NewUint(10),
		},
		{
			name: "len_err",
			fn:   "len",
			args: []execute.Value{&slowtesting.MockValue{
				LengthErr: errors.NoLengthError(slowtesting.NewMockType()),
			}},
			wantErr: errors.NoLengthError(slowtesting.NewMockType()),
		},
		{
			name:    "len_few_args",
			fn:      "len",
			args:    []execute.Value{},
			wantErr: errors.CallError("len", 0, 1),
		},
		{
			name: "len_many_args",
			fn:   "len",
			args: []execute.Value{
				&slowtesting.MockValue{},
				&slowtesting.MockValue{},
			},
			wantErr: errors.CallError("len", 2, 1),
		},
		{
			name: "print_string",
			fn:   "print",
			args: []execute.Value{
				types.NewStr("foo"),
			},
			want:         types.Null,
			wantPrintlns: []string{"foo"},
		},
		{
			name: "print_value",
			fn:   "print",
			args: []execute.Value{
				&slowtesting.MockValue{StringRet: "MOCK_VALUE"},
			},
			want:         types.Null,
			wantPrintlns: []string{"MOCK_VALUE"},
		},
		{
			name: "print_many",
			fn:   "print",
			args: []execute.Value{
				&slowtesting.MockValue{StringRet: "MV1"},
				&slowtesting.MockValue{StringRet: "MV2"},
				&slowtesting.MockValue{StringRet: "MV3"},
			},
			want:         types.Null,
			wantPrintlns: []string{"MV1MV2MV3"},
		},
		// TODO: range
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var printed []string
			println = func(s string) {
				printed = append(printed, s)
			}
			env := RootEnvironment.NewFrame()
			fn, err := env.Get(tc.fn)
			if err != nil {
				t.Fatalf("Get() returned unexpected error: %v", err)
			}
			c, err := fn.ToCallable()
			if err != nil {
				t.Fatalf("fn.ToCallable() returned unexpected error: %v", err)
			}
			got, err := c.Call(env, tc.args...)
			if diff := cmp.Diff(tc.wantErr, err, allowUnexported); diff != "" {
				t.Errorf("c.Call() returned incorrect error (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.want, got, allowUnexported); diff != "" {
				t.Errorf("c.Call() returned incorrect value (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.wantPrintlns, printed); diff != "" {
				t.Errorf("println called incorrectly (-want +got):\n%s", diff)
			}
		})
	}
}

func TestBuiltins_type(t *testing.T) {
	tests := []struct {
		value execute.Value
		want  string
	}{
		{
			types.NewBool(true),
			"bool",
		},
		{
			types.NewFloat(1),
			"float",
		},
		{
			types.NewFunc("", nil, nil),
			"func",
		},
		{
			types.NewGenerator(nil),
			"generator",
		},
		{
			types.NewInt(1),
			"int",
		},
		{
			types.NewList(nil),
			"list",
		},
		{
			types.NewMap(),
			"map",
		},
		{
			types.Null,
			"null",
		},
		{
			types.NewStr(""),
			"str",
		},
		{
			types.NewUint(1),
			"uint",
		},
	}
	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			env := RootEnvironment.NewFrame()
			f, _ := env.Get("type")
			c, _ := f.ToCallable()
			gotv, err := c.Call(env, tc.value)
			if err != nil {
				t.Fatalf("c.Call() returned unexpected error: %v", err)
			}
			if got, _ := gotv.(*types.Str).ToStr(); got != tc.want {
				t.Errorf("builtin type() returned %v, want %v", got, tc.want)
			}
		})
	}
}

func TestNewRootEnvironmentIsFrozen(t *testing.T) {
	// Attempt to reassign a variable that is bound to a built-in, so we know it's already declared.
	_, err := RootEnvironment.Set("import", &slowtesting.MockValue{})
	if err == nil {
		t.Errorf("env.Set did not error")
	}
}
