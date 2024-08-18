package types

import "github.com/chrispyles/slow/execute"

var typeHierarchy = map[execute.Type]int{
	FloatType: 0,
	IntType:   1,
	UintType:  2,
	BoolType:  3, // TODO: validate that bool is treated like uint
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
