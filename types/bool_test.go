package types

import (
	"fmt"
	"testing"

	"github.com/chrispyles/slow/execute"
	typestesting "github.com/chrispyles/slow/types/internal/testing"
)

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
			t.Run(fmt.Sprintf("%+v", tc.in), func(t *testing.T) {
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
		// TODO
	})

	t.Run("GetAttribute", func(t *testing.T) {
		// TODO
	})

	t.Run("GetIndex", func(t *testing.T) {
		// TODO
	})

	t.Run("HasAttribute", func(t *testing.T) {
		// TODO
	})

	t.Run("HashBytes", func(t *testing.T) {
		// TODO
	})

	t.Run("Length", func(t *testing.T) {
		// TODO
	})

	t.Run("SetAttribute", func(t *testing.T) {
		// TODO
	})

	t.Run("SetIndex", func(t *testing.T) {
		// TODO
	})

	t.Run("String", func(t *testing.T) {
		// TODO
	})

	t.Run("ToBool", func(t *testing.T) {
		// TODO
	})

	t.Run("ToBytes", func(t *testing.T) {
		// TODO
	})

	t.Run("ToCallable", func(t *testing.T) {
		// TODO
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
