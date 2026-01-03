package types

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/chrispyles/slow/src/errors"
	"github.com/chrispyles/slow/src/execute"
	slowtesting "github.com/chrispyles/slow/src/testing"
	testhelpers "github.com/chrispyles/slow/src/testing/helpers"
	typestesting "github.com/chrispyles/slow/src/types/internal/testing"
	"github.com/google/go-cmp/cmp"
)

func TestBytesType(t *testing.T) {
	tc := typestesting.TypeTestCase{
		Type: BytesType,
		NewTestCases: []typestesting.NewTestCase{
			{
				In: &slowtesting.MockValue{
					ToBytesRet: []byte{0x01, 0x02},
				},
				Want: NewBytes([]byte{0x01, 0x02}),
			},
			{
				In: &slowtesting.MockValue{
					ToBytesErr: errors.NewValueError("no"),
				},
				WantErr: errors.NewValueError("no"),
			},
		},
		WantString:    "bytes",
		WantIsNumeric: true,
	}
	tc.Run(t)
}

func TestBytes(t *testing.T) {
	t.Run("CloneIfPrimitive", func(t *testing.T) {
		in := NewBytes(nil)
		got := in.CloneIfPrimitive()
		testhelpers.CheckDiff(t, "CloneIfPrimitive()", in, got, cmp.AllowUnexported(*in))
		if reflect.ValueOf(in).Pointer() == reflect.ValueOf(got).Pointer() {
			t.Errorf("CloneIfPrimitive() did not create a clone")
		}
	})

	t.Run("CompareTo", func(t *testing.T) {
		for _, tc := range []struct {
			in     *Bytes
			other  execute.Value
			want   int
			wantOk bool
		}{
			{
				in:     NewBytes([]byte{0xA0}),
				other:  NewBytes([]byte{0xA0}),
				want:   0,
				wantOk: true,
			},
			{
				in:     NewBytes([]byte{0xA0}),
				other:  NewBytes([]byte{0x40}),
				want:   1,
				wantOk: true,
			},
			{
				in:     NewBytes([]byte{0xA0}),
				other:  NewBytes([]byte{0xFF}),
				want:   -1,
				wantOk: true,
			},
			{
				in:     NewBytes([]byte{0xA0}),
				other:  NewBool(false),
				want:   0,
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
			in    *Bytes
			other execute.Value
			want  bool
		}{
			{
				in:    NewBytes([]byte{0xA0}),
				other: NewBytes([]byte{0xA0}),
				want:  true,
			},
			{
				in:    NewBytes([]byte{0xA0}),
				other: NewBytes([]byte{0xA1}),
				want:  false,
			},
			{
				in:    NewBytes([]byte{0xA0}),
				other: NewBool(true),
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
		got, err := NewBytes([]byte{0xA0}).GetAttribute("foo")
		want := errors.NewAttributeError(BytesType, "foo")
		testhelpers.CheckDiff(t, "GetAttribute() error", want, err, allowUnexported)
		if got, want := got, (execute.Value)(nil); got != want {
			t.Errorf("GetAttribute() = %v, want %v", got, want)
		}
	})

	t.Run("GetIndex", func(t *testing.T) {
		v := NewBytes([]byte{0x00, 0x01, 0x02, 0x03})
		for _, tc := range []struct {
			idx     execute.Value
			want    execute.Value
			wantErr error
		}{
			{
				idx:  NewInt(1),
				want: NewBytes([]byte{0x01}),
			},
			{
				idx:  NewInt(-2),
				want: NewBytes([]byte{0x02}),
			},
			{
				idx:  NewUint(1),
				want: NewBytes([]byte{0x01}),
			},
			{
				idx:  NewBool(false),
				want: NewBytes([]byte{0x00}),
			},
			{
				idx:  NewInt(-int64(len(v.value))),
				want: NewBytes([]byte{0x00}),
			},
			{
				idx:     NewInt(int64(len(v.value))),
				wantErr: errors.NewIndexError(fmt.Sprintf("%d", len(v.value))),
			},
			{
				idx:     NewInt(-int64(len(v.value)) - 1),
				wantErr: errors.NewIndexError(fmt.Sprintf("%d", -1*len(v.value)-1)),
			},
			{
				idx:     NewFloat(1),
				wantErr: errors.NonNumericIndexError(FloatType, BytesType),
			},
			{
				idx:     NewStr("1"),
				wantErr: errors.NonNumericIndexError(StrType, BytesType),
			},
		} {
			t.Run(fmt.Sprintf("%+v", tc.idx), func(t *testing.T) {
				got, err := v.GetIndex(tc.idx)
				testhelpers.CheckDiff(t, "GetIndex() error", tc.wantErr, err, allowUnexported)
				testhelpers.CheckDiff(t, "GetIndex()", tc.want, got, allowUnexported)
			})
		}
	})

	t.Run("HasAttribute", func(t *testing.T) {
		got := NewBytes([]byte{0xA0}).HasAttribute("foo")
		if got, want := got, false; got != want {
			t.Errorf("HasAttribute() = %v, want %v", got, want)
		}
	})

	t.Run("HashBytes", func(t *testing.T) {
		got, err := NewBytes([]byte{0xA0, 0xB0}).HashBytes()
		testhelpers.CheckDiff(t, "HashBytes() error", nil, err, allowUnexported)
		testhelpers.CheckDiff(t, "HashBytes()", []byte{0xA0, 0xB0}, got, allowUnexported)
	})

	t.Run("Length", func(t *testing.T) {
		got, err := NewBytes([]byte{0xA0, 0xB0}).Length()
		testhelpers.CheckDiff(t, "Length() error", nil, err, allowUnexported)
		testhelpers.CheckDiff(t, "Length()", uint64(2), got, allowUnexported)
	})

	t.Run("SetAttribute", func(t *testing.T) {
		err := NewBytes([]byte{0xA0}).SetAttribute("foo", Null)
		want := errors.NewAttributeError(BytesType, "foo")
		testhelpers.CheckDiff(t, "SetAttribute() error", want, err, allowUnexported)
	})

	t.Run("SetIndex", func(t *testing.T) {
		err := NewBytes([]byte{0xA0}).SetIndex(NewInt(0), Null)
		want := errors.SetIndexNotSupported(BytesType)
		testhelpers.CheckDiff(t, "SetIndex() error", want, err, allowUnexported)
	})

	t.Run("String", func(t *testing.T) {
		got := NewBytes([]byte{0xA0, 0xB0}).String()
		testhelpers.CheckDiff(t, "String()", "0xA0B0", got, allowUnexported)
	})

	t.Run("ToBool", func(t *testing.T) {
		got := NewBytes([]byte{0x00, 0x01, 0x00}).ToBool()
		if got, want := got, true; got != want {
			t.Errorf("ToBool() = %v, want %v", got, want)
		}
		got = NewBytes([]byte{0x00, 0x00, 0x00}).ToBool()
		if got, want := got, false; got != want {
			t.Errorf("ToBool() = %v, want %v", got, want)
		}
	})

	t.Run("ToBytes", func(t *testing.T) {
		got, err := NewBytes([]byte{0xA0, 0xB0}).ToBytes()
		testhelpers.CheckDiff(t, "ToBytes() error", nil, err, allowUnexported)
		testhelpers.CheckDiff(t, "ToBytes()", []byte{0xA0, 0xB0}, got, allowUnexported)
	})

	t.Run("ToCallable", func(t *testing.T) {
		got, err := NewBytes([]byte{0xA0}).ToCallable()
		want := errors.NewTypeError(BytesType, FuncType)
		testhelpers.CheckDiff(t, "ToCallable() error", want, err, allowUnexported)
		if got, want := got, (execute.Callable)(nil); got != want {
			t.Errorf("ToCallable() = %v, want %v", got, want)
		}
	})

	t.Run("ToFloat", func(t *testing.T) {
		// TODO
	})

	t.Run("ToInt", func(t *testing.T) {
		// TODO
	})

	t.Run("ToIterator", func(t *testing.T) {
		// TODO
	})

	t.Run("ToStr", func(t *testing.T) {
		// TODO
	})

	t.Run("ToUint", func(t *testing.T) {
		// TODO
	})

	t.Run("Type", func(t *testing.T) {
		// TODO
	})

	t.Run("type_methods", func(t *testing.T) {
		// TODO
	})
}
