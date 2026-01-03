package errors_test

import (
	"testing"

	"github.com/chrispyles/slow/src/errors"
	slowtesting "github.com/chrispyles/slow/src/testing"
)

var (
	mt1 = &slowtesting.MockType{StringRet: "MockType1"}
	mt2 = &slowtesting.MockType{StringRet: "MockType2"}
)

func TestTypeError(t *testing.T) {
	e := errors.NewTypeError(mt1, mt2)

	got, want := e.Error(), "TypeError: type \"MockType1\" cannot be used as type \"MockType2\""
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}

func TestIncompatibleType(t *testing.T) {
	e := errors.IncompatibleType(mt1, "+")

	got, want := e.Error(), "TypeError: type \"MockType1\" cannot be used with operator \"+\""
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}

func TestIncompatibleTypes(t *testing.T) {
	e := errors.IncompatibleTypes(mt1, mt2, "+")

	got, want := e.Error(), "TypeError: types \"MockType1\" and \"MockType2\" cannot be used together with operator \"+\""
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}

func TestIncomparableType(t *testing.T) {
	e := errors.IncomparableType(mt1, ">")

	got, want := e.Error(), "TypeError: type \"MockType1\" cannot be used with comparison operator \">\""
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}

func TestCallError(t *testing.T) {
	e := errors.CallError("foo", 1, 2)

	got, want := e.Error(), "TypeError: function foo accepts 2 arguments but 1 were given"
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}

func TestNoLengthError(t *testing.T) {
	e := errors.NoLengthError(mt1)

	got, want := e.Error(), "TypeError: type \"MockType1\" does not have a length"
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}

func TestUnhashableTypeError(t *testing.T) {
	e := errors.UnhashableTypeError(slowtesting.NewMockType())

	got, want := e.Error(), "TypeError: type \"MockType\" is not hashable"
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}

func TestInvalidTypeCastTarget(t *testing.T) {
	e := errors.InvalidTypeCastTarget(slowtesting.NewMockType())

	got, want := e.Error(), "TypeError: type \"MockType\" does not support type casting"
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}

func TestTypeErrorFromMessage(t *testing.T) {
	e := errors.TypeErrorFromMessage("foo")

	got, want := e.Error(), "TypeError: foo"
	if got != want {
		t.Errorf("Error() returned incorrect value: got %q, want %q", got, want)
	}
}
