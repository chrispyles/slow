package operators

import (
	"math"

	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/types"
)

type BinaryOperator string

// TODO: and, or, xor
const (
	// TODO: +=, -=, *=, /=, %=, //=
	BinOp_EMPTY BinaryOperator = ""
	BinOp_PLUS  BinaryOperator = "+"
	BinOp_MINUS BinaryOperator = "-"
	BinOp_TIMES BinaryOperator = "*"
	BinOp_DIV   BinaryOperator = "/"
	BinOp_MOD   BinaryOperator = "%"
	BinOp_FDIV  BinaryOperator = "//"
	BinOp_EXP   BinaryOperator = "**"
	BinOp_EQ    BinaryOperator = "=="
	BinOp_NEQ   BinaryOperator = "!="
	BinOp_LT    BinaryOperator = "<"
	BinOp_LEQ   BinaryOperator = "<="
	BinOp_GT    BinaryOperator = ">"
	BinOp_GEQ   BinaryOperator = ">="
)

var allBinOps = map[BinaryOperator]bool{
	BinOp_PLUS:  true,
	BinOp_MINUS: true,
	BinOp_TIMES: true,
	BinOp_DIV:   true,
	BinOp_MOD:   true,
	BinOp_FDIV:  true,
	BinOp_EXP:   true,
	BinOp_EQ:    true,
	BinOp_NEQ:   true,
	BinOp_LT:    true,
	BinOp_LEQ:   true,
	BinOp_GT:    true,
	BinOp_GEQ:   true,
}

func ToBinaryOp(maybeOp string) (BinaryOperator, bool) {
	op := BinaryOperator(maybeOp)
	if allBinOps[op] {
		return op, true
	}
	return BinOp_EMPTY, false
}

func (o BinaryOperator) Value(l, r execute.Value) (execute.Value, error) {
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

	caster, ok := NewTypeCaster(lt, rt)
	if o == BinOp_DIV {
		// TODO: this should error if either rt or lt is not numeric
		caster, ok = newFloatCaster()
	}
	if !ok {
		return nil, errors.IncompatibleTypes(lt, rt, o.String())
	}

	doCast := func() (execute.Value, execute.Value, error) {
		return caster.Cast(l, r)
	}

	if o.IsComparison() {
		// Non-numeric types are always != to another value if it is of a different type, and are
		// incomparable otherwise.
		if (!lt.IsNumeric() || !rt.IsNumeric()) && lt != rt {
			switch o {
			case BinOp_EQ:
				return types.NewBool(false), nil
			case BinOp_NEQ:
				return types.NewBool(true), nil
			default:
				return nil, errors.IncompatibleTypes(lt, rt, o.String())
			}
		}
		if lt.IsNumeric() && rt.IsNumeric() {
			lc, rc, err := doCast()
			if err != nil {
				return nil, err
			}
			r, err := compareNumeric(lc, rc)
			if err != nil {
				return nil, err
			}
			switch o {
			case BinOp_EQ:
				return types.NewBool(r == equalTo), nil
			case BinOp_NEQ:
				return types.NewBool(r != equalTo), nil
			case BinOp_LT:
				return types.NewBool(r == lessThan), nil
			case BinOp_LEQ:
				return types.NewBool(r == lessThan || r == equalTo), nil
			case BinOp_GT:
				return types.NewBool(r == greaterThan), nil
			case BinOp_GEQ:
				return types.NewBool(r == greaterThan || r == equalTo), nil
			}
		}
		switch o {
		// TODO: use Value.Equals()
		case BinOp_EQ:
			switch lt {
			case types.FuncType:
				return types.NewBool(l == r), nil
			case types.NullType:
				// null is a singleton, so if two values are of NullType, this must be true.
				return types.NewBool(true), nil
			case types.StrType:
				return types.NewBool(must(l.ToStr()) == must(r.ToStr())), nil
			}
		case BinOp_NEQ:
			switch lt {
			case types.FuncType:
				return types.NewBool(l != r), nil
			case types.NullType:
				// null is a singleton, so if two values are of NullType, this must be false.
				return types.NewBool(false), nil
			case types.StrType:
				return types.NewBool(must(l.ToStr()) != must(r.ToStr())), nil
			}
		case BinOp_LT:
			switch lt {
			case types.FuncType:
				fallthrough
			case types.NullType:
				return nil, errors.IncomparableType(lt, o.String())
			case types.StrType:
				return types.NewBool(must(l.ToStr()) < must(r.ToStr())), nil
			}
		case BinOp_LEQ:
			switch lt {
			case types.FuncType:
				fallthrough
			case types.NullType:
				return nil, errors.IncomparableType(lt, o.String())
			case types.StrType:
				return types.NewBool(must(l.ToStr()) <= must(r.ToStr())), nil
			}
		case BinOp_GT:
			switch lt {
			case types.FuncType:
				fallthrough
			case types.NullType:
				return nil, errors.IncomparableType(lt, o.String())
			case types.StrType:
				return types.NewBool(must(l.ToStr()) > must(r.ToStr())), nil
			}
		case BinOp_GEQ:
			switch lt {
			case types.FuncType:
				fallthrough
			case types.NullType:
				return nil, errors.IncomparableType(lt, o.String())
			case types.StrType:
				return types.NewBool(must(l.ToStr()) >= must(r.ToStr())), nil
			}
		}

		return nil, errors.IncompatibleTypes(lt, rt, o.String())
	}

	lc, rc, err := doCast()
	if err != nil {
		return nil, err
	}

	// Return an error if there is an attempt to divide by zero.
	if caster.dest.IsNumeric() && must(rc.ToFloat()) == 0 && (o == BinOp_DIV || o == BinOp_MOD || o == BinOp_FDIV) {
		return nil, errors.NewZeroDivisionError()
	}

	// TODO: handle more types in all of these
	switch o {
	case BinOp_PLUS:
		switch caster.dest {
		case types.BoolType:
			return types.NewUint(must(lc.ToUint()) + must(rc.ToUint())), nil
		case types.FloatType:
			return types.NewFloat(must(lc.ToFloat()) + must(rc.ToFloat())), nil
		case types.IntType:
			return types.NewInt(must(lc.ToInt()) + must(rc.ToInt())), nil
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
		// TODO: what happens if this receives non-numeric type?
		if !caster.dest.IsNumeric() {
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

	// TODO: return error
	return nil, errors.IncompatibleTypes(lt, rt, o.String())
}

func (o BinaryOperator) IsComparison() bool {
	return o == BinOp_EQ || o == BinOp_NEQ || o == BinOp_LT || o == BinOp_LEQ || o == BinOp_GT || o == BinOp_GEQ
}

func (o BinaryOperator) String() string {
	return string(o)
}

var operatorPrecedence = map[BinaryOperator]int{
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
func (o BinaryOperator) Compare(other BinaryOperator) bool {
	l, r := operatorPrecedence[o], operatorPrecedence[other]
	return l < r
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
