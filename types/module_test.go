package types

import (
	"testing"

	typestesting "github.com/chrispyles/slow/types/internal/testing"
)

func TestModuleType(t *testing.T) {
	// TODO: NewTestCases
	tc := typestesting.TypeTestCase{
		Type:          ModuleType,
		WantString:    "module",
		WantIsNumeric: false,
	}
	tc.Run(t)
}

// TODO: test type implementation
