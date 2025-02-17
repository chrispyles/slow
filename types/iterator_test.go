package types

import (
	"testing"

	typestesting "github.com/chrispyles/slow/types/internal/testing"
)

func TesIteratorrType(t *testing.T) {
	// TODO: NewTestCases
	tc := typestesting.TypeTestCase{
		Type:          IteratorType,
		WantString:    "iterator",
		WantIsNumeric: false,
	}
	tc.Run(t)
}

// TODO: test type implementation
