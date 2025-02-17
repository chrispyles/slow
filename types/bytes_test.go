package types

import (
	"testing"

	typestesting "github.com/chrispyles/slow/types/internal/testing"
)

func TestBytesType(t *testing.T) {
	// TODO: NewTestCases
	tc := typestesting.TypeTestCase{
		Type:          BytesType,
		WantString:    "bytes",
		WantIsNumeric: true,
	}
	tc.Run(t)
}

// TODO: test type implementation
