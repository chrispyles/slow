package types

import (
	"testing"

	typestesting "github.com/chrispyles/slow/types/internal/testing"
)

func TestGeneratorType(t *testing.T) {
	// TODO: NewTestCases
	tc := typestesting.TypeTestCase{
		Type:          GeneratorType,
		WantString:    "generator",
		WantIsNumeric: false,
	}
	tc.Run(t)
}

// TODO: test type implementation
