package operators

import (
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/types"
)

const (
	lessThan    = -1
	equalTo     = 0
	greaterThan = 1
)

func compareNumeric(left, right execute.Value) (int, error) {
	if left.Type() != right.Type() {
		panic("compareNumeric called with values of different types")
	}
	switch left.Type() {
	case types.BoolType:
		l, r := left.ToBool(), right.ToBool()
		if l == r {
			return equalTo, nil
		} else if !l && r {
			return lessThan, nil
		}
		return greaterThan, nil
	case types.FloatType:
		l, r := must(left.ToFloat()), must(right.ToFloat())
		return doComparison(l, r), nil
	case types.IntType:
		l, r := must(left.ToInt()), must(right.ToInt())
		return doComparison(l, r), nil
	case types.UintType:
		l, r := must(left.ToUint()), must(right.ToUint())
		return doComparison(l, r), nil
	}
	panic("compareNumeric called with non-numeric values")
}

func doComparison[T uint64 | int64 | float64](l, r T) int {
	if l == r {
		return equalTo
	} else if l < r {
		return lessThan
	} else {
		return greaterThan
	}
}
