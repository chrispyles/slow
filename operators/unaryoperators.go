package operators

import (
	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/types"
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
			return types.NewInt(-1 * must(v.ToInt())), nil
		default:
			return nil, errors.IncompatibleType(v.Type(), o.String())
		}
	case UnOp_NOT:
		// Each type's ToBool method determines the value's truthiness.
		return types.NewBool(!v.ToBool()), nil
	}
	// TODO: other operators
	return nil, nil
}

func (o *UnaryOperator) String() string {
	return o.chars
}
