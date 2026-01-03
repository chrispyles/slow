package execute_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/chrispyles/slow/internal/execute"
	slowtesting "github.com/chrispyles/slow/internal/testing"
	"github.com/google/go-cmp/cmp"
)

func TestBlock(t *testing.T) {
	var val execute.Value = &slowtesting.MockValue{}
	expr := &slowtesting.MockExpression{ExecuteRet: val}
	b := execute.Block{expr}
	env := execute.NewEnvironment()
	got, err := b.Execute(env)
	if err != nil {
		t.Errorf("b.Execute() returned an unexpected error: %v", err)
	}
	if got != val {
		t.Errorf("b.Execute() returned wrong value: got %v, want %v", got, val)
	}
	if diff := cmp.Diff([]uintptr{reflect.ValueOf(env).Pointer()}, expr.Calls); diff != "" {
		t.Errorf("b.Execute() called expr.Execute() incorrectly (-want +got):\n%s", diff)
	}

	exprErr := errors.New("doh")
	b = execute.Block{
		&slowtesting.MockExpression{ExecuteRet: val},
		&slowtesting.MockExpression{ExecuteErr: exprErr},
		&slowtesting.MockExpression{ExecuteRet: val},
	}
	got, err = b.Execute(env)
	if err != exprErr {
		t.Errorf("b.Execute() returned wrong error: got %v, want %v", err, exprErr)
	}
	if got != nil {
		t.Errorf("b.Execute() returned wrong value: got %v, want %v", got, nil)
	}
	for i, e := range b {
		want := 1
		if i > 1 {
			want = 0
		}
		if got := len(e.(*slowtesting.MockExpression).Calls); got != want {
			t.Errorf("b.Execute() called expr Execute() wrong number of times: got %v, want %v", got, want)
		}
	}
}
