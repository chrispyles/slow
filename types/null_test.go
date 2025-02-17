package types

import (
	"testing"

	typestesting "github.com/chrispyles/slow/types/internal/testing"
)

func TestNullType(t *testing.T) {
	// TODO: NewTestCases
	tc := typestesting.TypeTestCase{
		Type:          NullType,
		WantString:    "null",
		WantIsNumeric: false,
	}
	tc.Run(t)
}

// TODO: test type implementation
