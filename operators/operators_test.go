package operators

import "testing"

func TestToUnaryOp(t *testing.T) {
	for opS, op := range unaryOperators {
		t.Run(opS, func(t *testing.T) {
			got, ok := ToUnaryOp(opS)
			if got, want := ok, true; got != want {
				t.Errorf("ToUnaryOp()[1] = %v, want %v", got, want)
			}
			if want := op; got != want {
				t.Errorf("ToUnaryOp()[0] = %v, want %v", got, want)
			}
		})
	}
	t.Run("invalid_op", func(t *testing.T) {
		got, ok := ToUnaryOp("foo")
		if got, want := ok, false; got != want {
			t.Errorf("ToUnaryOp()[1] = %v, want %v", got, want)
		}
		if want := (*UnaryOperator)(nil); got != want {
			t.Errorf("ToUnaryOp()[0] = %v, want %v", got, want)
		}
	})
}

func TestToBinaryOp(t *testing.T) {
	for opS, op := range binaryOperators {
		t.Run(opS, func(t *testing.T) {
			got, ok := ToBinaryOp(opS)
			if got, want := ok, true; got != want {
				t.Errorf("ToBinaryOp()[1] = %v, want %v", got, want)
			}
			if want := op; got != want {
				t.Errorf("ToBinaryOp()[0] = %v, want %v", got, want)
			}
		})
	}
	t.Run("invalid_op", func(t *testing.T) {
		got, ok := ToBinaryOp("foo")
		if got, want := ok, false; got != want {
			t.Errorf("ToBinaryOp()[1] = %v, want %v", got, want)
		}
		if want := (*BinaryOperator)(nil); got != want {
			t.Errorf("ToBinaryOp()[0] = %v, want %v", got, want)
		}
	})
}
