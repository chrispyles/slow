package operators

import (
	"iter"
	"maps"
)

var (
	unaryOperators  = make(map[string]*UnaryOperator)
	binaryOperators = make(map[string]*BinaryOperator)
)

func AllUnaryOperators() iter.Seq[*UnaryOperator] {
	return maps.Values(unaryOperators)
}

func AllBinaryOperators() iter.Seq[*BinaryOperator] {
	return maps.Values(binaryOperators)
}

type unaryOperatorType string
type binaryOperatorType string

const (
	unOp  unaryOperatorType  = "u"
	binOp binaryOperatorType = "b"
)

type Operator[T unaryOperatorType | binaryOperatorType] struct {
	opType T
	chars  string
}

type UnaryOperator Operator[unaryOperatorType]
type BinaryOperator Operator[binaryOperatorType]

func newUnaryOperator(chars string) *UnaryOperator {
	op := &Operator[unaryOperatorType]{unOp, chars}
	uop := (*UnaryOperator)(op)
	unaryOperators[chars] = uop
	return uop
}

func newBinaryOperator(chars string) *BinaryOperator {
	op := &Operator[binaryOperatorType]{binOp, chars}
	bop := (*BinaryOperator)(op)
	binaryOperators[chars] = bop
	return bop
}

func ToUnaryOp(maybeOp string) (*UnaryOperator, bool) {
	if op, ok := unaryOperators[maybeOp]; ok {
		return op, true
	}
	return nil, false
}

func ToBinaryOp(maybeOp string) (*BinaryOperator, bool) {
	if op, ok := binaryOperators[maybeOp]; ok {
		return op, true
	}
	return nil, false
}

// unary operators
var (
	UnOp_POS  = newUnaryOperator("+")
	UnOp_NEG  = newUnaryOperator("-")
	UnOp_NOT  = newUnaryOperator("!")
	UnOp_INCR = newUnaryOperator("++")
	UnOp_DECR = newUnaryOperator("--")
)

// binary operators
var (
	// arithmetic operators
	BinOp_PLUS  = newBinaryOperator("+")
	BinOp_MINUS = newBinaryOperator("-")
	BinOp_TIMES = newBinaryOperator("*")
	BinOp_DIV   = newBinaryOperator("/")
	BinOp_MOD   = newBinaryOperator("%")
	BinOp_FDIV  = newBinaryOperator("//")
	BinOp_EXP   = newBinaryOperator("**")
	BinOp_AND   = newBinaryOperator("&&")
	BinOp_OR    = newBinaryOperator("||")
	BinOp_XOR   = newBinaryOperator("^^")

	// reassignment operators
	BinOp_RPLUS  = newBinaryOperator("+=")
	BinOp_RMINUS = newBinaryOperator("-=")
	BinOp_RTIMES = newBinaryOperator("*=")
	BinOp_RDIV   = newBinaryOperator("/=")
	BinOp_RMOD   = newBinaryOperator("%=")
	BinOp_RFDIV  = newBinaryOperator("//=")
	BinOp_REXP   = newBinaryOperator("**=")
	BinOp_RAND   = newBinaryOperator("&&=")
	BinOp_ROR    = newBinaryOperator("||=")
	BinOp_RXOR   = newBinaryOperator("^^=")

	// comparison operators
	BinOp_EQ  = newBinaryOperator("==")
	BinOp_NEQ = newBinaryOperator("!=")
	BinOp_LT  = newBinaryOperator("<")
	BinOp_LEQ = newBinaryOperator("<=")
	BinOp_GT  = newBinaryOperator(">")
	BinOp_GEQ = newBinaryOperator(">=")
)
