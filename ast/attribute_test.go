package ast

import (
	"testing"

	asttesting "github.com/chrispyles/slow/ast/internal/testing"
	"github.com/chrispyles/slow/execute"
	slowtesting "github.com/chrispyles/slow/testing"
	"github.com/chrispyles/slow/types"
)

func TestAttributeNode(t *testing.T) {
	mv := &slowtesting.MockValue{Attributes: map[string]execute.Value{"foo": types.NewInt(1)}}
	asttesting.RunTestCase(t, asttesting.TestCase{
		Name:        "success",
		Node:        &AttributeNode{Left: &VariableNode{Name: "mockValue"}, Right: "foo"},
		Env:         slowtesting.MustMakeEnv(t, map[string]execute.Value{"mockValue": mv}),
		Want:        types.NewInt(1),
		WantSameEnv: true,
	})
	_, wantErr := mv.GetAttribute("bar")
	if wantErr == nil {
		t.Fatalf("GetAttribute() for nonexistent attribute returned nil error")
	}
	asttesting.RunTestCase(t, asttesting.TestCase{
		Name:        "error",
		Node:        &AttributeNode{Left: &VariableNode{Name: "mockValue"}, Right: "bar"},
		Env:         slowtesting.MustMakeEnv(t, map[string]execute.Value{"mockValue": mv}),
		WantErr:     wantErr,
		WantSameEnv: true,
	})
}
