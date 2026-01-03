package interpreter

import (
	"bufio"
	"strings"
	"testing"

	"github.com/chrispyles/slow/internal/execute"
	"github.com/google/go-cmp/cmp"
)

func TestRun(t *testing.T) {
	t.Run("noninteractive", func(t *testing.T) {
		evalCalls := setup(t)

		Run("foo", nil)

		if diff := cmp.Diff([]string{"foo"}, *evalCalls); diff != "" {
			t.Errorf("Run() called eval incorrectly (-want +got):\n%s", diff)
		}
	})

	t.Run("interactive", func(t *testing.T) {
		evalCalls := setup(t)

		input := strings.NewReader("bar\nbaz\n\n")

		defer func() {
			if err := recover(); err == nil {
				t.Errorf("Run() did not run forever")
			}
			if diff := cmp.Diff([]string{"foo", "bar\n", "baz\n"}, *evalCalls); diff != "" {
				t.Errorf("Run() called eval incorrectly (-want +got):\n%s", diff)
			}
		}()
		Run("foo", input)
	})
}

func setup(t *testing.T) *[]string {
	origEval := eval
	evalCalls := &[]string{}
	eval = func(c string, env *execute.Environment, print bool) {
		*evalCalls = append(*evalCalls, c)
	}
	origRead := read
	read = func(r *bufio.Reader) (string, error) {
		s, err := r.ReadString('\n')
		if err != nil {
			// Call panic to kill the Run function when out of input since it runs indefinitely. The
			// calling function should recover from this panic.
			panic(err)
		}
		return s, nil
	}
	t.Cleanup(func() {
		eval = origEval
		read = origRead
	})
	return evalCalls
}
