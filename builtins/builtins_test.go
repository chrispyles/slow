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

// TODO: tests for exit()

func TestNewRootEnvironment(t *testing.T) {
	tests := []struct {
		name    string
		fn      string
		args    []execute.Value
		want    execute.Value
		wantErr error
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
		},
		{
			name: "print_int",
		},
		{
			name: "print_float",
		},
		{
			name: "print_many",
		},
		// TODO: print, range, type
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			env := NewRootEnvironment()
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
		})
	}
}

func TestNewRootEnvironmentIsFrozen(t *testing.T) {
	env := NewRootEnvironment()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("env.Set did not panic")
		}
	}()
	env.Set("foo", &slowtesting.MockValue{})
}
