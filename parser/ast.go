package parser

import (
	"github.com/chrispyles/slow/ast"
	"github.com/chrispyles/slow/execute"
)

func Parse(s string) (execute.AST, error) {
	b, err := parse(s)
	if err != nil {
		return nil, err
	}
	return ast.New(b), nil
}
