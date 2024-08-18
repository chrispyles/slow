package operators

import "github.com/chrispyles/slow/execute"
import "github.com/chrispyles/slow/types"

type TypeCaster struct {
	dest      execute.Type
	castLeft  bool
	castRight bool
}

func NewTypeCaster(leftType execute.Type, rightType execute.Type) (*TypeCaster, bool) {
	dest, ok := types.CommonNumericType(leftType, rightType)
	if !ok {
		return nil, false
	}
	// TODO: validate that this WAI
	return &TypeCaster{dest, dest != leftType, dest != rightType}, true
}

func newFloatCaster() (*TypeCaster, bool) {
	return &TypeCaster{types.FloatType, true, true}, true
}

func (c *TypeCaster) singleCast(val execute.Value) (execute.Value, error) {
	var res execute.Value
	var err error
	switch c.dest {
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

func (c *TypeCaster) Cast(l, r execute.Value) (execute.Value, execute.Value, error) {
	var lc execute.Value
	var rc execute.Value
	var err error
	if c.castLeft {
		lc, err = c.singleCast(l)
	} else {
		lc = l
	}
	if err != nil {
		return nil, nil, err
	}
	if c.castRight {
		rc, err = c.singleCast(r)
	} else {
		rc = r
	}
	return lc, rc, err
}

func (c *TypeCaster) Dest() execute.Type {
	return c.dest
}
