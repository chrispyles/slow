package operators

import (
	"math"

	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/types"
)

type UnaryOperator string

const (
	UnOp_EMPTY UnaryOperator = ""
	UnOp_POS   UnaryOperator = "+"
	UnOp_NEG   UnaryOperator = "-"
	UnOp_NOT   UnaryOperator = "!"
	UnOp_INCR  UnaryOperator = "++"
	UnOp_DECR  UnaryOperator = "--"
)

var allUnOps = map[UnaryOperator]bool{
	UnOp_POS:  true,
	UnOp_NEG:  true,
	UnOp_NOT:  true,
	UnOp_INCR: true,
	UnOp_DECR: true,
}

// TODO: this is duplication w/ above is annoying
func ToUnaryOp(maybeOp string) (UnaryOperator, bool) {
	op := UnaryOperator(maybeOp)
	if allUnOps[op] {
		return op, true
	}
	return UnOp_EMPTY, false
}

func (o UnaryOperator) Value(v execute.Value) (execute.Value, error) {
	switch o {
	case UnOp_POS:
		// TODO
		return nil, nil
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
	}
	// TODO: other operators
	return nil, nil
}

func (o UnaryOperator) String() string {
	return string(o)
}

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

// TODO: this is duplication w/ above is annoying
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

type binOpCaster struct {
	dest      execute.Type
	castLeft  bool
	castRight bool
}

func (c *binOpCaster) singleCast(val execute.Value) (execute.Value, error) {
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

func (c *binOpCaster) Cast(l, r execute.Value) (execute.Value, execute.Value, error) {
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

var typeHierarchy = map[execute.Type]int{
	types.FloatType: 0,
	types.IntType:   1,
	types.UintType:  2,
	types.BoolType:  3,
	// TODO: account for all types
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

	var caster *binOpCaster
	if o == BinOp_DIV {
		caster = &binOpCaster{types.FloatType, true, true}
	} else {
		ltp, rtp := typeHierarchy[lt], typeHierarchy[rt]
		if ltp < rtp {
			caster = &binOpCaster{lt, false, true}
		} else {
			caster = &binOpCaster{rt, true, false}
		}
	}
	doCast := func() (execute.Value, execute.Value, error) {
		return caster.Cast(l, r)
	}

	if o.IsComparison() {
		// Non-numeric types are always != to another value if it is of a different type, and are
		// incomparable.
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
			switch o {
			case BinOp_EQ:
				return types.NewBool(lc.Equals(rc)), nil
			case BinOp_NEQ:
				return types.NewBool(!lc.Equals(rc)), nil
				// TODO: other operators
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

	// TODO: handle more types in all of these
	switch o {
	case BinOp_PLUS:
		switch caster.dest {
		case types.FloatType:
			return types.NewFloat(must(lc.ToFloat()) + must(rc.ToFloat())), nil
		case types.IntType:
			return types.NewInt(must(lc.ToInt()) + must(rc.ToInt())), nil
		case types.UintType:
			return types.NewUint(must(lc.ToUint()) + must(rc.ToUint())), nil
		}
	case BinOp_MINUS:
		switch caster.dest {
		case types.FloatType:
			return types.NewFloat(must(lc.ToFloat()) - must(rc.ToFloat())), nil
		case types.IntType:
			return types.NewInt(must(lc.ToInt()) - must(rc.ToInt())), nil
		case types.UintType:
			return types.NewUint(must(lc.ToUint()) - must(rc.ToUint())), nil
		}
	case BinOp_TIMES:
		switch caster.dest {
		case types.FloatType:
			return types.NewFloat(must(lc.ToFloat()) * must(rc.ToFloat())), nil
		case types.IntType:
			return types.NewInt(must(lc.ToInt()) * must(rc.ToInt())), nil
		case types.UintType:
			return types.NewUint(must(lc.ToUint()) * must(rc.ToUint())), nil
		}
	case BinOp_DIV:
		switch caster.dest {
		case types.FloatType:
			return types.NewFloat(must(lc.ToFloat()) / must(rc.ToFloat())), nil
		case types.IntType:
			return types.NewInt(must(lc.ToInt()) / must(rc.ToInt())), nil
		case types.UintType:
			return types.NewUint(must(lc.ToUint()) / must(rc.ToUint())), nil
		}
	case BinOp_FDIV:
		switch caster.dest {
		case types.FloatType:
			return types.NewInt(int64(math.Floor(must(lc.ToFloat()) / must(rc.ToFloat())))), nil
		case types.IntType:
			return types.NewInt(must(lc.ToInt()) / must(rc.ToInt())), nil
		case types.UintType:
			return types.NewUint(must(lc.ToUint()) / must(rc.ToUint())), nil
		}
	case BinOp_EXP:
		// TODO: what happens if this receives non-numeric type?
		lf, rf := must(l.ToFloat()), must(r.ToFloat())
		val := math.Pow(lf, rf)
		switch caster.dest {
		case types.FloatType:
			return types.NewFloat(val), nil
		case types.IntType:
			return types.NewInt(int64(val)), nil
		case types.UintType:
			return types.NewUint(uint64(val)), nil
		}
	}

	// TODO: return error
	return nil, nil
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

type TernaryOperator string

const (
	TernOp_EMPTY TernaryOperator = ""
)

func ToTernaryOp(maybeOp string) (TernaryOperator, bool) {
	return TernOp_EMPTY, false
}

func must[T any](v T, _ error) T {
	return v
}
