package operators

import (
	"github.com/chrispyles/slow/src/errors"
	"github.com/chrispyles/slow/src/execute"
	"github.com/chrispyles/slow/src/types"
)

func (o *UnaryOperator) Value(v execute.Value) (execute.Value, error) {
	switch o {
	case UnOp_POS:
		if v.Type() != types.FloatType &&
			v.Type() != types.IntType &&
			v.Type() != types.UintType &&
			v.Type() != types.BoolType {
			return nil, errors.IncompatibleType(v.Type(), o.String())
		}
		return v.CloneIfPrimitive(), nil
	case UnOp_NEG:
		switch v.Type() {
		case types.FloatType:
			return types.NewFloat(-1 * must(v.ToFloat())), nil
		case types.IntType:
			return types.NewInt(-1 * must(v.ToInt())), nil
		case types.UintType:
			fallthrough
		case types.BoolType:
			return types.NewInt(-1 * must(v.ToInt())), nil
		default:
			return nil, errors.IncompatibleType(v.Type(), o.String())
		}
	case UnOp_NOT:
		// Each type's ToBool method determines the value's truthiness.
		return types.NewBool(!v.ToBool()), nil
	case UnOp_INCR:
		return BinOp_PLUS.Value(v, types.NewUint(1))
	case UnOp_DECR:
		return BinOp_MINUS.Value(v, types.NewUint(1))
	}
	panic("unandled unary operator in UnaryOperator.Value()")
}

func (o *UnaryOperator) IsReassignmentOperator() bool {
	return o == UnOp_INCR || o == UnOp_DECR
}

func (o *UnaryOperator) String() string {
	return o.chars
}
