package types

import (
	"testing"

	typestesting "github.com/chrispyles/slow/types/internal/testing"
)

func TestStrType(t *testing.T) {
	// TODO: NewTestCases
	tc := typestesting.TypeTestCase{
		Type:          StrType,
		WantString:    "str",
		WantIsNumeric: false,
	}
	tc.Run(t)
}

// TODO: test type implementation
