package operators

import (
	"math"

	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/types"
)

var logicalOperators = map[*BinaryOperator]bool{
	BinOp_AND: true,
	BinOp_OR:  true,
	BinOp_XOR: true,
}

var reassignmentToArithmeticOperator = map[*BinaryOperator]*BinaryOperator{
	BinOp_RPLUS:  BinOp_PLUS,
	BinOp_RMINUS: BinOp_MINUS,
	BinOp_RTIMES: BinOp_TIMES,
	BinOp_RDIV:   BinOp_DIV,
	BinOp_RMOD:   BinOp_MOD,
	BinOp_RFDIV:  BinOp_FDIV,
	BinOp_REXP:   BinOp_EXP,
	BinOp_RAND:   BinOp_AND,
	BinOp_ROR:    BinOp_OR,
	BinOp_RXOR:   BinOp_XOR,
}

var incomparableTypes = map[execute.Type]bool{
	types.FuncType:      true,
	types.GeneratorType: true,
	types.IteratorType:  true,
	types.ListType:      true,
	types.NullType:      true,
}

func (o *BinaryOperator) Value(l, r execute.Value) (execute.Value, error) {
	// If this is a reassignment operator, convert it to its arithmetic version to calculate the new
	// value.
	if ao, ok := reassignmentToArithmeticOperator[o]; ok {
		o = ao
	}

	lt, rt := l.Type(), r.Type()
	if o == BinOp_MOD {
		// Floats with no remainder can be treated as ints.
		if lt == types.FloatType && !l.(*types.Float).HasRemainder() {
			lt = types.IntType
		}
		if rt == types.FloatType && !r.(*types.Float).HasRemainder() {
			rt = types.IntType
		}

		if lt != types.IntType && lt != types.UintType {
			return nil, errors.IncompatibleType(lt, o.String())
		}
		if rt != types.IntType && rt != types.UintType {
			return nil, errors.IncompatibleType(rt, o.String())
		}

		if lt == types.UintType && rt == types.UintType {
			return types.NewUint(must(l.ToUint()) % must(r.ToUint())), nil
		}
		return types.NewInt(must(l.ToInt()) % must(r.ToInt())), nil
	}

	if logicalOperators[o] {
		lb, rb := l.ToBool(), r.ToBool()
		switch o {
		case BinOp_AND:
			if !lb {
				return l.CloneIfPrimitive(), nil
			}
			return r.CloneIfPrimitive(), nil
		case BinOp_OR:
			if lb {
				return l.CloneIfPrimitive(), nil
			}
			return r.CloneIfPrimitive(), nil
		case BinOp_XOR:
			return types.NewBool((lb || rb) && !(lb && rb)), nil
		}
	}

	caster, ok := newTypeCaster(lt, rt)
	if o == BinOp_DIV {
		if !lt.IsNumeric() || !rt.IsNumeric() {
			return nil, errors.IncompatibleTypes(lt, rt, o.String())
		}
		caster, ok = newFloatCaster()
	}
	if !ok {
		return nil, errors.IncompatibleTypes(lt, rt, o.String())
	}

	doCast := func() (execute.Value, execute.Value, error) {
		return caster.Cast(l, r)
	}

	if o.IsComparison() {
		if incomparableTypes[lt] || incomparableTypes[rt] {
			// These types are all pass-by-reference, and will be equal iff they are the same instance.
			switch o {
			case BinOp_EQ:
				return types.NewBool(lt == rt), nil
			case BinOp_NEQ:
				return types.NewBool(lt != rt), nil
			default:
				return nil, errors.IncompatibleTypes(lt, rt, o.String())
			}
		}

		comparator, comparable := l.CompareTo(r)
		if !comparable {
			return nil, errors.IncompatibleTypes(lt, rt, o.String())
		}

		switch o {
		case BinOp_EQ:
			return types.NewBool(comparator == 0), nil
		case BinOp_NEQ:
			return types.NewBool(comparator != 0), nil
		case BinOp_LT:
			return types.NewBool(comparator < 0), nil
		case BinOp_LEQ:
			return types.NewBool(comparator <= 0), nil
		case BinOp_GT:
			return types.NewBool(comparator > 0), nil
		case BinOp_GEQ:
			return types.NewBool(comparator >= 0), nil
		default:
			panic("unhandled comparison operator in BinaryOperator.Value()")
		}
	}

	lc, rc, err := doCast()
	if err != nil {
		return nil, err
	}

	// Return an error if there is an attempt to divide by zero.
	if caster.dest.IsNumeric() && must(rc.ToFloat()) == 0 && (o == BinOp_DIV || o == BinOp_MOD || o == BinOp_FDIV) {
		return nil, errors.NewZeroDivisionError()
	}

	switch o {
	case BinOp_PLUS:
		switch caster.dest {
		case types.BoolType:
			return types.NewUint(must(lc.ToUint()) + must(rc.ToUint())), nil
		case types.FloatType:
			return types.NewFloat(must(lc.ToFloat()) + must(rc.ToFloat())), nil
		case types.IntType:
			return types.NewInt(must(lc.ToInt()) + must(rc.ToInt())), nil
		case types.StrType:
			return types.NewStr(must(lc.ToStr()) + must(rc.ToStr())), nil
		case types.UintType:
			return types.NewUint(must(lc.ToUint()) + must(rc.ToUint())), nil
		}
	case BinOp_MINUS:
		switch caster.dest {
		case types.BoolType:
			return types.NewUint(must(lc.ToUint()) - must(rc.ToUint())), nil
		case types.FloatType:
			return types.NewFloat(must(lc.ToFloat()) - must(rc.ToFloat())), nil
		case types.IntType:
			return types.NewInt(must(lc.ToInt()) - must(rc.ToInt())), nil
		case types.UintType:
			return types.NewUint(must(lc.ToUint()) - must(rc.ToUint())), nil
		}
	case BinOp_TIMES:
		switch caster.dest {
		case types.BoolType:
			return types.NewUint(must(lc.ToUint()) * must(rc.ToUint())), nil
		case types.FloatType:
			return types.NewFloat(must(lc.ToFloat()) * must(rc.ToFloat())), nil
		case types.IntType:
			return types.NewInt(must(lc.ToInt()) * must(rc.ToInt())), nil
		case types.UintType:
			return types.NewUint(must(lc.ToUint()) * must(rc.ToUint())), nil
		}
	case BinOp_DIV:
		// BinOp_DIV always sets the destination type to float
		return types.NewFloat(must(lc.ToFloat()) / must(rc.ToFloat())), nil
	case BinOp_FDIV:
		switch caster.dest {
		case types.BoolType:
			return types.NewUint(must(lc.ToUint()) / must(rc.ToUint())), nil
		case types.FloatType:
			return types.NewInt(int64(math.Floor(must(lc.ToFloat()) / must(rc.ToFloat())))), nil
		case types.IntType:
			return types.NewInt(must(lc.ToInt()) / must(rc.ToInt())), nil
		case types.UintType:
			return types.NewUint(must(lc.ToUint()) / must(rc.ToUint())), nil
		}
	case BinOp_EXP:
		if !lt.IsNumeric() || !rt.IsNumeric() {
			return nil, errors.IncompatibleType(caster.dest, o.String())
		}
		lf, rf := must(l.ToFloat()), must(r.ToFloat())
		val := math.Pow(lf, rf)
		switch caster.dest {
		case types.BoolType:
			return types.NewUint(uint64(val)), nil
		case types.FloatType:
			return types.NewFloat(val), nil
		case types.IntType:
			return types.NewInt(int64(val)), nil
		case types.UintType:
			return types.NewUint(uint64(val)), nil
		}
	}

	return nil, errors.IncompatibleTypes(lt, rt, o.String())
}

func (o *BinaryOperator) IsComparison() bool {
	return o == BinOp_EQ ||
		o == BinOp_NEQ ||
		o == BinOp_LT ||
		o == BinOp_LEQ ||
		o == BinOp_GT ||
		o == BinOp_GEQ
}

func (o *BinaryOperator) IsReassignmentOperator() bool {
	_, ok := reassignmentToArithmeticOperator[o]
	return ok
}

func (o *BinaryOperator) String() string {
	return o.chars
}

var operatorPrecedence = map[*BinaryOperator]int{
	BinOp_EXP:   -2,
	BinOp_TIMES: -1,
	BinOp_DIV:   -1,
	BinOp_MOD:   -1,
	BinOp_FDIV:  -1,
	BinOp_PLUS:  0,
	BinOp_MINUS: 0,
	BinOp_EQ:    1,
	BinOp_NEQ:   1,
	BinOp_LT:    1,
	BinOp_LEQ:   1,
	BinOp_GT:    1,
	BinOp_GEQ:   1,
}

// Compare returns true if this BinaryOperator takes precendence over other (i.e. this operation
// should be evaluated first).
func (o *BinaryOperator) Compare(other *BinaryOperator) bool {
	if ao, ok := reassignmentToArithmeticOperator[o]; ok {
		o = ao
	}
	if ao, ok := reassignmentToArithmeticOperator[other]; ok {
		other = ao
	}
	l, r := operatorPrecedence[o], operatorPrecedence[other]
	return l < r
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
