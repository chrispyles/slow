package types

import (
	"testing"

	typestesting "github.com/chrispyles/slow/types/internal/testing"
)

func TesIteratorType(t *testing.T) {
	// TODO: NewTestCases
	tc := typestesting.TypeTestCase{
		Type:          IteratorType,
		WantString:    "iterator",
		WantIsNumeric: false,
	}
	tc.Run(t)
}

func TestIterator(t *testing.T) {
	t.Run("CompareTo", func(t *testing.T) {
		// TODO
	})

	t.Run("Equals", func(t *testing.T) {
		// TODO
	})

	t.Run("GetAttribute", func(t *testing.T) {
		// TODO
	})

	t.Run("GetIndex", func(t *testing.T) {
		// TODO
	})

	t.Run("HasAttribute", func(t *testing.T) {
		// TODO
	})

	t.Run("HashBytes", func(t *testing.T) {
		// TODO
	})

	t.Run("Length", func(t *testing.T) {
		// TODO
	})

	t.Run("SetAttribute", func(t *testing.T) {
		// TODO
	})

	t.Run("SetIndex", func(t *testing.T) {
		// TODO
	})

	t.Run("String", func(t *testing.T) {
		// TODO
	})

	t.Run("ToBool", func(t *testing.T) {
		// TODO
	})

	t.Run("ToBytes", func(t *testing.T) {
		// TODO
	})

	t.Run("ToCallable", func(t *testing.T) {
		// TODO
	})

	t.Run("ToFloat", func(t *testing.T) {
		// TODO
	})

	t.Run("ToInt", func(t *testing.T) {
		// TODO
	})

	t.Run("ToIterator", func(t *testing.T) {
		// TODO
	})

	t.Run("ToStr", func(t *testing.T) {
		// TODO
	})

	t.Run("ToUint", func(t *testing.T) {
		// TODO
	})

	t.Run("Type", func(t *testing.T) {
		// TODO
	})

	t.Run("type_methods", func(t *testing.T) {
		// TODO
	})
}
