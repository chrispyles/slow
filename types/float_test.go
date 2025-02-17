package types

import (
	"testing"

	typestesting "github.com/chrispyles/slow/types/internal/testing"
)

func TestFloatType(t *testing.T) {
	// TODO: NewTestCases
	tc := typestesting.TypeTestCase{
		Type:          FloatType,
		WantString:    "float",
		WantIsNumeric: true,
	}
	tc.Run(t)
}

// TODO: test type implementation
