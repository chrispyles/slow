package types

import (
	"testing"

	typestesting "github.com/chrispyles/slow/types/internal/testing"
)

func TestListType(t *testing.T) {
	// TODO: NewTestCases
	tc := typestesting.TypeTestCase{
		Type:          ListType,
		WantString:    "list",
		WantIsNumeric: false,
	}
	tc.Run(t)
}

// TODO: test type implementation
