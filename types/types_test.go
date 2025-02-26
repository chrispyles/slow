package types

import (
	"github.com/chrispyles/slow/errors"
	"github.com/google/go-cmp/cmp"
)

var allowUnexported = cmp.AllowUnexported(
	errors.SlowError{},
	Bool{},
	Bytes{},
	Float{},
	Func{},
	Generator{},
	Int{},
	Iterator{},
	List{},
	Module{},
	Str{},
	Uint{},
)
