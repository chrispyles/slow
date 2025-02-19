package types

import (
	"fmt"
	"testing"

	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
	testhelpers "github.com/chrispyles/slow/testing/helpers"
	typestesting "github.com/chrispyles/slow/types/internal/testing"
	"github.com/google/go-cmp/cmp"
)

var allowUnexported = cmp.AllowUnexported(errors.SlowError{})

func TestBoolType(t *testing.T) {
	// TODO: NewTestCases
	tc := typestesting.TypeTestCase{
		Type:          BoolType,
		WantString:    "bool",
		WantIsNumeric: true,
	}
	tc.Run(t)
}

func TestBool(t *testing.T) {
	t.Run("CompareTo", func(t *testing.T) {
		for _, tc := range []struct {
			in     *Bool
			other  execute.Value
			want   int
			wantOk bool
		}{
			{
				in:     NewBool(true),
				other:  NewBool(true),
				want:   0,
				wantOk: true,
			},
			{
				in:     NewBool(false),
				other:  NewBool(true),
				want:   -1,
				wantOk: true,
			},
			{
				in:     NewBool(true),
				other:  NewBool(false),
				want:   1,
				wantOk: true,
			},
			{
				in:     NewBool(false),
				other:  NewBool(false),
				want:   0,
				wantOk: true,
			},
			{
				in:     NewBool(true),
				other:  NewFloat(-1),
				want:   1,
				wantOk: true,
			},
			{
				in:     NewBool(true),
				other:  NewFloat(1),
				want:   0,
				wantOk: true,
			},
			{
				in:     NewBool(true),
				other:  NewFloat(2),
				want:   -1,
				wantOk: true,
			},
			{
				in:     NewBool(true),
				other:  NewInt(-1),
				want:   1,
				wantOk: true,
			},
			{
				in:     NewBool(true),
				other:  NewInt(1),
				want:   0,
				wantOk: true,
			},
			{
				in:     NewBool(true),
				other:  NewInt(2),
				want:   -1,
				wantOk: true,
			},
			{
				in:     NewBool(true),
				other:  NewUint(0),
				want:   1,
				wantOk: true,
			},
			{
				in:     NewBool(true),
				other:  NewUint(1),
				want:   0,
				wantOk: true,
			},
			{
				in:     NewBool(true),
				other:  NewUint(2),
				want:   -1,
				wantOk: true,
			},
			{
				in:     NewBool(false),
				other:  NewFloat(-1),
				want:   1,
				wantOk: true,
			},
			{
				in:     NewBool(false),
				other:  NewFloat(0),
				want:   0,
				wantOk: true,
			},
			{
				in:     NewBool(false),
				other:  NewFloat(1),
				want:   -1,
				wantOk: true,
			},
			{
				in:     NewBool(false),
				other:  NewInt(-1),
				want:   1,
				wantOk: true,
			},
			{
				in:     NewBool(false),
				other:  NewInt(0),
				want:   0,
				wantOk: true,
			},
			{
				in:     NewBool(false),
				other:  NewInt(1),
				want:   -1,
				wantOk: true,
			},
			{
				in:     NewBool(false),
				other:  NewUint(0),
				want:   0,
				wantOk: true,
			},
			{
				in:     NewBool(false),
				other:  NewUint(1),
				want:   -1,
				wantOk: true,
			},
			{
				in:     NewBool(false),
				other:  NewStr("hi"),
				wantOk: false,
			},
		} {
			t.Run(fmt.Sprintf("%+v__%+v", tc.in, tc.other), func(t *testing.T) {
				got, ok := tc.in.CompareTo(tc.other)
				if got, want := ok, tc.wantOk; got != want {
					t.Errorf("CompareTo() returned incorrect ok value: got %v, want %v", got, want)
				}
				if got, want := got, tc.want; got != want {
					t.Errorf("CompareTo() returned incorrect value: got %v, want %v", got, want)
				}
			})
		}
	})

	t.Run("Equals", func(t *testing.T) {
		for _, tc := range []struct {
			in    *Bool
			other execute.Value
			want  bool
		}{
			{
				in:    NewBool(true),
				other: NewBool(true),
				want:  true,
			},
			{
				in:    NewBool(false),
				other: NewBool(true),
				want:  false,
			},
			{
				in:    NewBool(true),
				other: NewFloat(1),
				want:  false,
			},
		} {
			t.Run(fmt.Sprintf("%+v__%+v", tc.in, tc.other), func(t *testing.T) {
				if got, want := tc.in.Equals(tc.other), tc.want; got != want {
					t.Errorf("Equals() returned incorrect value: got %v, want %v", got, want)
				}
			})
		}
	})

	t.Run("GetAttribute", func(t *testing.T) {
		got, err := NewBool(true).GetAttribute("foo")
		want := errors.NewAttributeError(BoolType, "foo")
		testhelpers.CheckDiff(t, "GetAttribute() error", want, err, allowUnexported)
		if got, want := got, (execute.Value)(nil); got != want {
			t.Errorf("GetAttribute() = %v, want %v", got, want)
		}
	})

	t.Run("GetIndex", func(t *testing.T) {
		got, err := NewBool(true).GetIndex(NewInt(0))
		want := errors.IndexingNotSupported(BoolType)
		testhelpers.CheckDiff(t, "GetIndex() error", want, err, allowUnexported)
		if got, want := got, (execute.Value)(nil); got != want {
			t.Errorf("GetIndex() = %v, want %v", got, want)
		}
	})

	t.Run("HasAttribute", func(t *testing.T) {
		got := NewBool(true).HasAttribute("foo")
		if got, want := got, false; got != want {
			t.Errorf("HasAttribute() = %v, want %v", got, want)
		}
	})

	t.Run("HashBytes", func(t *testing.T) {
		got, err := NewBool(true).HashBytes()
		if err != nil {
			t.Errorf("HashBytes() returned non-nil error: %v", err)
		}
		testhelpers.CheckDiff(t, "HashBytes()", []byte{0x01}, got)
		got, err = NewBool(false).HashBytes()
		if err != nil {
			t.Errorf("HashBytes() returned non-nil error: %v", err)
		}
		testhelpers.CheckDiff(t, "HashBytes()", []byte{0x00}, got)
	})

	t.Run("Length", func(t *testing.T) {
		got, err := NewBool(true).Length()
		want := errors.NoLengthError(BoolType)
		testhelpers.CheckDiff(t, "Length() error", want, err, allowUnexported)
		if got, want := got, uint64(0); got != want {
			t.Errorf("Length() = %v, want %v", got, want)
		}
	})

	t.Run("SetAttribute", func(t *testing.T) {
		err := NewBool(true).SetAttribute("foo", NewInt(1))
		want := errors.NewAttributeError(BoolType, "foo")
		testhelpers.CheckDiff(t, "SetAttribute() error", want, err, allowUnexported)
	})

	t.Run("SetIndex", func(t *testing.T) {
		err := NewBool(true).SetIndex(NewInt(0), NewInt(1))
		want := errors.IndexingNotSupported(BoolType)
		testhelpers.CheckDiff(t, "SetIndex() error", want, err, allowUnexported)
	})

	t.Run("String", func(t *testing.T) {
		got := NewBool(true).String()
		if got, want := got, "true"; got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
		got = NewBool(false).String()
		if got, want := got, "false"; got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})

	t.Run("ToBool", func(t *testing.T) {
		got := NewBool(true).ToBool()
		if got, want := got, true; got != want {
			t.Errorf("ToBool() = %v, want %v", got, want)
		}
		got = NewBool(false).ToBool()
		if got, want := got, false; got != want {
			t.Errorf("ToBool() = %v, want %v", got, want)
		}
	})

	t.Run("ToBytes", func(t *testing.T) {
		got, err := NewBool(true).ToBytes()
		if err != nil {
			t.Errorf("ToBytes() returned non-nil error: %v", err)
		}
		testhelpers.CheckDiff(t, "ToBytes()", []byte{0x01}, got)
		got, err = NewBool(false).ToBytes()
		if err != nil {
			t.Errorf("ToBytes() returned non-nil error: %v", err)
		}
		testhelpers.CheckDiff(t, "ToBytes()", []byte{0x00}, got)
	})

	t.Run("ToCallable", func(t *testing.T) {
		got, err := NewBool(true).ToCallable()
		want := errors.NewTypeError(BoolType, FuncType)
		testhelpers.CheckDiff(t, "ToCallable() error", want, err, allowUnexported)
		if got, want := got, (execute.Callable)(nil); got != want {
			t.Errorf("ToCallable() = %v, want %v", got, want)
		}
	})

	t.Run("ToFloat", func(t *testing.T) {
		got, err := NewBool(true).ToFloat()
		if err != nil {
			t.Errorf("ToFloat() return an unexpected error: %v", err)
		}
		if got, want := got, 1.0; got != want {
			t.Errorf("ToFloat() = %v, want %v", got, want)
		}
		got, err = NewBool(false).ToFloat()
		if err != nil {
			t.Errorf("ToFloat() return an unexpected error: %v", err)
		}
		if got, want := got, 0.0; got != want {
			t.Errorf("ToFloat() = %v, want %v", got, want)
		}
	})

	t.Run("ToInt", func(t *testing.T) {
		got, err := NewBool(true).ToInt()
		if err != nil {
			t.Errorf("ToInt() return an unexpected error: %v", err)
		}
		if got, want := got, int64(1); got != want {
			t.Errorf("ToInt() = %v, want %v", got, want)
		}
		got, err = NewBool(false).ToInt()
		if err != nil {
			t.Errorf("ToInt() return an unexpected error: %v", err)
		}
		if got, want := got, int64(0); got != want {
			t.Errorf("ToInt() = %v, want %v", got, want)
		}
	})

	t.Run("ToIterator", func(t *testing.T) {
		got, err := NewBool(true).ToIterator()
		want := errors.NewTypeError(BoolType, IteratorType)
		testhelpers.CheckDiff(t, "ToIterator() error", want, err, allowUnexported)
		if got, want := got, (execute.Iterator)(nil); got != want {
			t.Errorf("ToIterator() = %v, want %v", got, want)
		}
	})

	t.Run("ToStr", func(t *testing.T) {
		got, err := NewBool(true).ToStr()
		if err != nil {
			t.Errorf("ToStr() return an unexpected error: %v", err)
		}
		if got, want := got, "true"; got != want {
			t.Errorf("ToStr() = %v, want %v", got, want)
		}
		got, err = NewBool(false).ToStr()
		if err != nil {
			t.Errorf("ToStr() return an unexpected error: %v", err)
		}
		if got, want := got, "false"; got != want {
			t.Errorf("ToStr() = %v, want %v", got, want)
		}
	})

	t.Run("ToUint", func(t *testing.T) {
		got, err := NewBool(true).ToUint()
		if err != nil {
			t.Errorf("ToUint() return an unexpected error: %v", err)
		}
		if got, want := got, uint64(1); got != want {
			t.Errorf("ToUint() = %v, want %v", got, want)
		}
		got, err = NewBool(false).ToUint()
		if err != nil {
			t.Errorf("ToUint() return an unexpected error: %v", err)
		}
		if got, want := got, uint64(0); got != want {
			t.Errorf("ToUint() = %v, want %v", got, want)
		}
	})

	t.Run("Type", func(t *testing.T) {
		got := NewBool(true).Type()
		if got, want := got, BoolType; got != want {
			t.Errorf("Type() = %v, want %v", got, want)
		}
	})

	t.Run("type_methods", func(t *testing.T) {
		// Type bool has no methods.
	})
}
