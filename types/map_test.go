package types

import (
	"testing"

	typestesting "github.com/chrispyles/slow/types/internal/testing"
)

func TestMapType(t *testing.T) {
	// TODO: NewTestCases
	tc := typestesting.TypeTestCase{
		Type:          MapType,
		WantString:    "map",
		WantIsNumeric: false,
	}
	tc.Run(t)
}

// TODO: test type implementation
