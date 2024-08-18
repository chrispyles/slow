package execute

import "github.com/chrispyles/slow/errors"

type Environment struct {
	values map[string]Value
	parent *Environment
	frozen bool
}

func NewEnvironment() *Environment {
	return &Environment{values: make(map[string]Value)}
}

func (e *Environment) Declare(n string) error {
	if e.frozen {
		// TODO: should this return an error instead of panicking?
		panic("can't bind create or set variables in a frozen environment")
	}
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
	if e.frozen {
		// TODO: should this return an error instead of panicking?
		panic("can't bind create or set variables in a frozen environment")
	}
	if _, ok := e.values[n]; ok {
		e.values[n] = v
		return v, nil
	}
	// The condition below assumes that if the parent frame is frozen, all ancestor frames are also
	// frozen.
	if e.parent != nil && !e.parent.frozen {
		return e.parent.Set(n, v)
	}
	return nil, errors.NewNameError(n)
}

func (e *Environment) Freeze() {
	e.frozen = true
}
