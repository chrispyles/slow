package modules

import "github.com/chrispyles/slow/execute"

type Module interface {
	Name() string
	Import() (*execute.Environment, error)
}

var modules = map[string]Module{
	"fs": &fsModule{},
}

var AllModules []string

func Get(name string) (Module, bool) {
	m, ok := modules[name]
	return m, ok
}

func init() {
	AllModules = make([]string, len(modules))
	i := 0
	for m := range modules {
		AllModules[i] = m
		i++
	}
}
