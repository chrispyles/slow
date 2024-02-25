package types

import (
	"testing"

	"github.com/chrispyles/slow/execute"
)

var allTypes = []execute.Type{
	BoolType,
	FloatType,
	FuncType,
	IntType,
	IteratorType,
	StrType,
	UintType,
}

func TestTypes(t *testing.T) {
	tests := []struct {
		ty          execute.Type
		wantMatches map[execute.Type]bool
		wantString  string
	}{
		{
			ty: BoolType,
			wantMatches: map[execute.Type]bool{
				BoolType:     true,
				FloatType:    false,
				FuncType:     false,
				IntType:      false,
				IteratorType: false,
				StrType:      false,
				UintType:     false,
			},
			wantString: "bool",
		},
		{
			ty: FloatType,
			wantMatches: map[execute.Type]bool{
				BoolType:     false,
				FloatType:    true,
				FuncType:     false,
				IntType:      false,
				IteratorType: false,
				StrType:      false,
				UintType:     false,
			},
			wantString: "float",
		},
		{
			ty: FuncType,
			wantMatches: map[execute.Type]bool{
				BoolType:     false,
				FloatType:    false,
				FuncType:     true,
				IntType:      false,
				IteratorType: false,
				StrType:      false,
				UintType:     false,
			},
			wantString: "func",
		},
		{
			ty: IntType,
			wantMatches: map[execute.Type]bool{
				BoolType:     false,
				FloatType:    false,
				FuncType:     false,
				IntType:      true,
				IteratorType: false,
				StrType:      false,
				UintType:     false,
			},
			wantString: "int",
		},
		{
			ty: IteratorType,
			wantMatches: map[execute.Type]bool{
				BoolType:     false,
				FloatType:    false,
				FuncType:     false,
				IntType:      false,
				IteratorType: true,
				StrType:      false,
				UintType:     false,
			},
			wantString: "iterator",
		},
		{
			ty: StrType,
			wantMatches: map[execute.Type]bool{
				BoolType:     false,
				FloatType:    false,
				FuncType:     false,
				IntType:      false,
				IteratorType: false,
				StrType:      true,
				UintType:     false,
			},
			wantString: "str",
		},
		{
			ty: UintType,
			wantMatches: map[execute.Type]bool{
				BoolType:     false,
				FloatType:    false,
				FuncType:     false,
				IntType:      false,
				IteratorType: false,
				StrType:      false,
				UintType:     true,
			},
			wantString: "uint",
		},
	}

	for _, tc := range tests {
		t.Run(tc.ty.String(), func(t *testing.T) {
			for _, ty2 := range allTypes {
				want, ok := tc.wantMatches[ty2]
				if !ok {
					t.Fatalf("no wantMatches value for type %s", ty2)
				}
				if got := tc.ty.Matches(ty2); got != want {
					t.Errorf("Matches returned incorrect value: got %v, want %v", got, want)
				}
			}

			if got, want := tc.ty.String(), tc.wantString; got != want {
				t.Errorf("String returned incorrect value: got %q, want %q", got, want)
			}
		})
	}
}
