package modules

import (
	"os"
	"testing"

	"github.com/chrispyles/slow/internal/errors"
	"github.com/chrispyles/slow/internal/execute"
	slowcmpopts "github.com/chrispyles/slow/internal/testing/cmpopts"
	"github.com/chrispyles/slow/internal/types"
	"github.com/google/go-cmp/cmp"
)

type testCase struct {
	name    string
	args    []execute.Value
	want    execute.Value
	wantErr error
}

func makeTestCallback(fn string, tc testCase) func(*testing.T) {
	return func(t *testing.T) {
		m := &fsModule{}
		env, err := m.Import()
		if err != nil {
			t.Fatalf("m.Import() returned an unexepcted error: %v", err)
		}
		f, err := env.Get(fn)
		if err != nil {
			t.Fatalf("failed to get function: %v", err)
		}
		fc, err := f.ToCallable()
		if err != nil {
			t.Fatalf("failed to convert Value to callable: %v", err)
		}
		got, err := fc.Call(env, tc.args...)
		if diff := cmp.Diff(tc.wantErr, err, slowcmpopts.AllowUnexported()); diff != "" {
			t.Errorf("function returned an unexpected error (-want +got):\n%s", diff)
		}
		if diff := cmp.Diff(tc.want, got, slowcmpopts.AllowUnexported()); diff != "" {
			t.Errorf("function returned an unexpected diff (-want +got):\n%s", diff)
		}
	}
}

func Test_fs_Name(t *testing.T) {
	m := &fsModule{}
	if got, want := m.Name(), "fs"; got != want {
		t.Errorf("m.Name() = %q, want %q", got, want)
	}
}

func Test_fs_read(t *testing.T) {
	tests := []testCase{
		{
			name: "success",
			args: []execute.Value{types.NewStr("testdata/a_file.txt")},
			want: types.NewStr("this is a file\n"),
		},
		{
			name:    "no_args",
			args:    []execute.Value{},
			wantErr: errors.CallError("fs.read", 0, 1),
		},
		{
			name:    "too_many_args",
			args:    []execute.Value{types.NewStr("testdata/a_file.txt"), types.NewStr("testdata/a_file.txt")},
			wantErr: errors.CallError("fs.read", 2, 1),
		},
		{
			name:    "file_does_not_exist",
			args:    []execute.Value{types.NewStr("testdata/another_file.txt")},
			wantErr: errors.WrapFileError(os.ErrNotExist, "testdata/another_file.txt"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, makeTestCallback("read", tc))
	}
}

func Test_fs_readBytes(t *testing.T) {
	tests := []testCase{
		{
			name: "success",
			args: []execute.Value{types.NewStr("testdata/a_file.txt")},
			want: types.NewBytes([]byte("this is a file\n")),
		},
		{
			name:    "no_args",
			args:    []execute.Value{},
			wantErr: errors.CallError("fs.readBytes", 0, 1),
		},
		{
			name:    "too_many_args",
			args:    []execute.Value{types.NewStr("testdata/a_file.txt"), types.NewStr("testdata/a_file.txt")},
			wantErr: errors.CallError("fs.readBytes", 2, 1),
		},
		{
			name:    "file_does_not_exist",
			args:    []execute.Value{types.NewStr("testdata/another_file.txt")},
			wantErr: errors.WrapFileError(os.ErrNotExist, "testdata/another_file.txt"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, makeTestCallback("readBytes", tc))
	}
}
