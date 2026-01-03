package cmpopts

import (
	"testing"
	"unsafe"

	"github.com/chrispyles/slow/src/errors"
	"github.com/chrispyles/slow/src/execute"
	"github.com/chrispyles/slow/src/testing/helpers"
	"github.com/chrispyles/slow/src/types"
	"github.com/google/go-cmp/cmp"
)

func AllowUnexported(addl ...interface{}) cmp.Option {
	return cmp.AllowUnexported(
		append(
			[]interface{}{
				errors.SlowError{},
				execute.Environment{},
				types.Bool{},
				types.Bytes{},
				types.Float{},
				types.Func{},
				types.Generator{},
				types.Int{},
				types.Iterator{},
				types.List{},
				types.Module{},
				types.RangeIterator{},
				types.Str{},
				types.Uint{},
			},
			addl...)...,
	)
}

// Adapted from https://github.com/google/go-cmp/issues/162
func EquateFuncs() cmp.Option {
	return cmp.Comparer(func(x, y types.FuncImpl) bool {
		px := *(*unsafe.Pointer)(unsafe.Pointer(&x))
		py := *(*unsafe.Pointer)(unsafe.Pointer(&y))
		return px == py
	})
}

func CheckDiff(t *testing.T, name string, want, got interface{}, opts ...cmp.Option) {
	helpers.CheckDiff(t, name, want, got, AllowUnexported())
}
