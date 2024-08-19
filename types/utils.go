package types

import (
	"encoding/binary"
	"math"

	"github.com/chrispyles/slow/execute"
)

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
	switch any(v).(type) {
	case float64:
		u = math.Float64bits(any(v).(float64))
	case int64:
		u = uint64(v)
	case uint64:
		u = any(v).(uint64)
	default:
		panic("unhandled type in numToBytes()")
	}
	var buf [8]byte
	binary.BigEndian.AppendUint64(buf[:], u)
	return buf[:]
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
