package types

import (
	"testing"

	typestesting "github.com/chrispyles/slow/types/internal/testing"
)

func TestFuncType(t *testing.T) {
	// TODO: NewTestCases
	tc := typestesting.TypeTestCase{
		Type:          FuncType,
		WantString:    "func",
		WantIsNumeric: false,
	}
	tc.Run(t)
}

// TODO: test type implementation
