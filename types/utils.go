package types

import (
	"encoding/binary"
	"math"

	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
)

var allTypes = []execute.Type{
	BoolType,
	BytesType,
	FloatType,
	FuncType,
	GeneratorType,
	IntType,
	IteratorType,
	ListType,
	MapType,
	ModuleType,
	NullType,
	StrType,
	UintType,
}

var typeHierarchy = map[execute.Type]int{
	FloatType: 0,
	IntType:   1,
	UintType:  2,
	BoolType:  3,
}

func CommonNumericType(t1 execute.Type, t2 execute.Type) (execute.Type, bool) {
	if t1 == t2 {
		return t1, true
	}
	t1p, t1ok := typeHierarchy[t1]
	t2p, t2ok := typeHierarchy[t2]
	if !t1ok || !t2ok {
		return nil, false
	}
	if t1p < t2p {
		return t1, true
	} else {
		return t2, true
	}
}

func compareNumbers[T float64 | int64 | uint64](v1, v2 T) int {
	if v1 == v2 {
		return 0
	} else if v1 < v2 {
		return -1
	}
	return 1
}

func numToBytes[T float64 | int64 | uint64](v T) []byte {
	var u uint64
	switch v := any(v).(type) {
	case float64:
		u = math.Float64bits(v)
	case int64:
		u = uint64(v)
	case uint64:
		u = v
	default:
		panic("unhandled type in numToBytes()")
	}
	var buf [8]byte
	binary.BigEndian.AppendUint64(buf[:], u)
	return buf[:]
}

func numericIndex(v execute.Value, t errors.Type) (int, error) {
	var ret int
	if vu, ok := v.(*Uint); ok {
		ret = int(vu.value)
	} else if vi, ok := v.(*Int); ok {
		ret = int(vi.value)
	} else if vb, ok := v.(*Bool); ok {
		ret = int(must(vb.ToInt()))
	} else {
		return 0, errors.NonNumericIndexError(v.Type(), t)
	}
	return ret, nil
}

// normalizeIndex converts a possible-negative index value to a positive index value. If the index
// is out-of-bounds based on the provided container length, the index is returned unchanged and the
// second return value is false.
func normalizeIndex(idx, cLen int) (int, bool) {
	if idx >= cLen || idx < -1*cLen {
		return idx, false
	}
	if idx >= 0 {
		return idx, true
	}
	return cLen + idx, true
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
