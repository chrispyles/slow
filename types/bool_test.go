package types

import (
	"testing"

	typestesting "github.com/chrispyles/slow/types/internal/testing"
)

func TestBoolType(t *testing.T) {
	// TODO: NewTestCases
	tc := typestesting.TypeTestCase{
		Type:          BoolType,
		WantString:    "bool",
		WantIsNumeric: true,
	}
	tc.Run(t)
}

// TODO: test type implementation
