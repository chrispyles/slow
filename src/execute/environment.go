package execute

import (
	"fmt"
	"maps"

	"github.com/chrispyles/slow/src/errors"
)

type Environment struct {
	values map[string]Value
	consts map[string]bool
	parent *Environment
	frozen bool
}

func NewEnvironment() *Environment {
	return &Environment{values: make(map[string]Value), consts: make(map[string]bool)}
}

// FromMap returns a frozen Environment from the provided map.
func FromMap(values map[string]Value) *Environment {
	e := &Environment{values: values}
	e.Freeze()
	return e
}

// Copy returns a shallow copy of this environment.
func (e *Environment) Copy() *Environment {
	if e == nil {
		return nil
	}
	return &Environment{
		values: maps.Clone(e.values),
		consts: maps.Clone(e.consts),
		parent: e.parent,
		frozen: e.frozen,
	}
}

func (e *Environment) Declare(n string) error {
	if e.frozen {
		panic("can't declare or set variables in a frozen environment")
	}
	if _, ok := e.values[n]; ok {
		return errors.NewDeclarationError(n)
	}
	e.values[n] = nil
	return nil
}

func (e *Environment) DeclareConst(n string, v Value) (Value, error) {
	if err := e.Declare(n); err != nil {
		return nil, err
	}
	v, err := e.Set(n, v)
	if err != nil {
		return nil, err
	}
	e.consts[n] = true
	return v, nil
}

func (e *Environment) Get(n string) (Value, error) {
	v, ok := e.values[n]
	if !ok {
		if e.parent != nil {
			return e.parent.Get(n)
		}
		return nil, errors.NewNameError(n)
	}
	if v == nil {
		return nil, errors.NewValueError(fmt.Sprintf("variable %q is uninitialized", n))
	}
	return v, nil
}

func (e *Environment) Has(n string) bool {
	_, ok := e.values[n]
	return ok
}

func (e *Environment) NewFrame() *Environment {
	c := NewEnvironment()
	c.parent = e
	return c
}

func (e *Environment) Set(n string, v Value) (Value, error) {
	if e.frozen {
		return nil, errors.NewRuntimeError("cannot assign variables in a frozen environment")
	}
	if e.consts[n] {
		return nil, errors.TypeErrorFromMessage(fmt.Sprintf("cannot reassign constant %q", n))
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
