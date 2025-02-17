package types

import (
	"testing"

	typestesting "github.com/chrispyles/slow/types/internal/testing"
)

func TestUintType(t *testing.T) {
	// TODO: NewTestCases
	tc := typestesting.TypeTestCase{
		Type:          UintType,
		WantString:    "uint",
		WantIsNumeric: true,
	}
	tc.Run(t)
}

// TODO: test type implementation
