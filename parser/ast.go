package parser

import (
	"github.com/chrispyles/slow/config"
	"github.com/chrispyles/slow/execute"
	"github.com/sanity-io/litter"
)

type AST struct {
	Nodes execute.Block
}

func NewAST(s string) (execute.AST, error) {
	b, err := parse(s)
	if err != nil {
		return nil, err
	}
	return &AST{b}, nil
}

func (a *AST) Execute(e *execute.Environment) (execute.Value, error) {
	var val execute.Value
	var err error
	for _, n := range a.Nodes {
		val, err = n.Execute(e)
		if err != nil {
			return nil, err
		}

		if *config.Debug {
			print("<AST EXECUTE LITTER> ")
			litter.Dump(val)
		}
	}
	return val, err
}

func (a *AST) String() string {
	return litter.Sdump(a)
}

func init() {
	litter.Config.HidePrivateFields = false
}
