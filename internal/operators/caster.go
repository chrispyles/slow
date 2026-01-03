package operators

import (
	"github.com/chrispyles/slow/internal/execute"
	"github.com/chrispyles/slow/internal/types"
)

type typeCaster struct {
	dest      execute.Type
	castLeft  bool
	castRight bool
}

func newTypeCaster(leftType execute.Type, rightType execute.Type) (*typeCaster, bool) {
	dest, ok := types.CommonNumericType(leftType, rightType)
	if !ok {
		return nil, false
	}
	return &typeCaster{dest, dest != leftType, dest != rightType}, true
}

func newFloatCaster() (*typeCaster, bool) {
	return &typeCaster{types.FloatType, true, true}, true
}

func (c *typeCaster) singleCast(val execute.Value) (execute.Value, error) {
	var res execute.Value
	var err error
	switch c.dest {
	case types.BoolType:
		res = types.NewBool(val.ToBool())
	case types.FloatType:
		var v float64
		v, err = val.ToFloat()
		res = types.NewFloat(v)
	case types.IntType:
		var v int64
		v, err = val.ToInt()
		res = types.NewInt(v)
	case types.UintType:
		var v uint64
		v, err = val.ToUint()
		res = types.NewUint(v)
	}
	return res, err
}

func (c *typeCaster) Cast(l, r execute.Value) (execute.Value, execute.Value) {
	var lc execute.Value
	var rc execute.Value
	var err error
	if c.castLeft {
		lc, err = c.singleCast(l)
	} else {
		lc = l
	}
	if err != nil {
		// The typeCaster should fail during construction if the type cast is not possible, so
		// singleCast should never return an error.
		panic(err)
	}
	if c.castRight {
		rc, err = c.singleCast(r)
	} else {
		rc = r
	}
	if err != nil {
		// See comment in previous error handling block.
		panic(err)
	}
	return lc, rc
}

func (c *typeCaster) Dest() execute.Type {
	return c.dest
}
