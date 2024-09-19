package builtins

import (
	"testing"
)

func TestBuiltins_range(t *testing.T) {
	doBuiltinTest(t, []builtinTest{
		// {
		// 	name: "len_success",
		// 	fn:   "len",
		// 	args: []execute.Value{&slowtesting.MockValue{LengthRet: 10}},
		// 	want: types.NewUint(10),
		// },
		// TODO
	})
}
