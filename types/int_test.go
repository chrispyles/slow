package types

import (
	"testing"

	typestesting "github.com/chrispyles/slow/types/internal/testing"
)

func TestIntType(t *testing.T) {
	// TODO: NewTestCases
	tc := typestesting.TypeTestCase{
		Type:          IntType,
		WantString:    "int",
		WantIsNumeric: true,
	}
	tc.Run(t)
}

// TODO: test type implementation
