package ast

import (
	"fmt"
	"testing"

	asttesting "github.com/chrispyles/slow/internal/ast/internal/testing"
	"github.com/chrispyles/slow/internal/errors"
	slowtesting "github.com/chrispyles/slow/internal/testing"
	"github.com/chrispyles/slow/internal/types"
)

func TestCastNode(t *testing.T) {
	tests := []asttesting.TestCase{
		{
			Name: "success",
			Node: &CastNode{
				Expr: &ConstantNode{
					Value: &slowtesting.MockValue{},
				},
				Type: &slowtesting.MockType{NewRet: types.NewUint(42)},
			},
			Want: types.NewUint(42),
		},
		{
			Name: "new_err",
			Node: &CastNode{
				Expr: &ConstantNode{
					Value: &slowtesting.MockValue{},
				},
				Type: &slowtesting.MockType{NewErr: errors.NewValueError("doh")},
			},
			WantErr: errors.NewValueError("doh"),
		},
	}
	for ty := range castingUnsupportedTypes {
		tests = append(tests, asttesting.TestCase{
			Name: fmt.Sprintf("unsupported_%s", ty),
			Node: &CastNode{
				Expr: &ConstantNode{
					Value: &slowtesting.MockValue{},
				},
				Type: ty,
			},
			WantErr: errors.InvalidTypeCastTarget(ty),
		})
	}
	for _, tc := range tests {
		asttesting.RunTestCase(t, tc)
	}
}
