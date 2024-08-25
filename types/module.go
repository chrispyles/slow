package types

import (
	"fmt"

	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
)

type Module struct {
	name string
	env  *execute.Environment
}

func NewModule(name string, env *execute.Environment) *Module {
	return &Module{name, env}
}

func (v *Module) CloneIfPrimitive() execute.Value {
	return v
}

func (v *Module) CompareTo(o execute.Value) (int, bool) {
	return 0, false
}

func (v *Module) Equals(o execute.Value) bool {
	om, ok := o.(*Module)
	if !ok {
		return false
	}
	return v == om
}

func (v *Module) GetAttribute(a string) (execute.Value, error) {
	return v.env.Get(a)
}

func (v *Module) HasAttribute(a string) bool {
	return v.env.Has(a)
}

func (v *Module) HashBytes() ([]byte, error) {
	return nil, errors.UnhashableTypeError(v.Type())
}

func (v *Module) Length() (uint64, error) {
	return 0, errors.NoLengthError(v.Type())
}

func (v *Module) SetAttribute(a string, _ execute.Value) error {
	if v.HasAttribute(a) {
		return errors.AssignmentError(v.Type(), a)
	}
	return errors.NewAttributeError(v.Type(), a)
}

func (v *Module) String() string {
	return fmt.Sprintf("<module %s>", v.name)
}

func (v *Module) ToBool() bool {
	return true
}

func (v *Module) ToCallable() (execute.Callable, error) {
	return nil, errors.NewTypeError(v.Type(), FuncType)
}

func (v *Module) ToFloat() (float64, error) {
	return 0, errors.NewTypeError(v.Type(), FloatType)
}

func (v *Module) ToInt() (int64, error) {
	return 0, errors.NewTypeError(v.Type(), IntType)
}

func (v *Module) ToIterator() (execute.Iterator, error) {
	return nil, errors.NewTypeError(v.Type(), IteratorType)
}

func (v *Module) ToStr() (string, error) {
	return "", errors.NewTypeError(v.Type(), StrType)
}

func (v *Module) ToUint() (uint64, error) {
	return 0, errors.NewTypeError(v.Type(), UintType)
}

func (v *Module) Type() execute.Type {
	return ModuleType
}
