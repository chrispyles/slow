package execute

import "github.com/chrispyles/slow/errors"

type Environment struct {
	values map[string]Value
	parent *Environment
}

func NewEnvironment() *Environment {
	return &Environment{values: make(map[string]Value)}
}

func (e *Environment) Declare(n string) error {
	if _, ok := e.values[n]; ok {
		return errors.NewDeclarationError(n)
	}
	e.values[n] = nil
	return nil
}

func (e *Environment) Get(n string) (Value, error) {
	v, ok := e.values[n]
	if !ok {
		if e.parent != nil {
			return e.parent.Get(n)
		}
		return nil, errors.NewNameError(n)
	}
	return v, nil
}

func (e *Environment) NewFrame() *Environment {
	c := NewEnvironment()
	c.parent = e
	return c
}

func (e *Environment) Set(n string, v Value) (Value, error) {
	if _, ok := e.values[n]; ok {
		e.values[n] = v
		return v, nil
	}
	if e.parent != nil {
		return e.parent.Set(n, v)
	}
	return nil, errors.NewNameError(n) // TODO: add message like "var has not been declared"
}
