package testing

import (
	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/types"
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
				types.Str{},
				types.Uint{},
			},
			addl...)...,
	)
}
