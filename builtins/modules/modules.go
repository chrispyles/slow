package modules

import "github.com/chrispyles/slow/execute"

type Module interface {
	Name() string
	Import() (*execute.Environment, error)
}

var modules = map[string]Module{
	"fs": &fsModule{},
}

func Get(name string) (Module, bool) {
	m, ok := modules[name]
	return m, ok
}
