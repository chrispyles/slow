package types

import "github.com/chrispyles/slow/execute"

// -------------------------------------------------------------------------------------------------

type boolType struct{}

func (t *boolType) IsNumeric() bool {
	return true
}

func (t *boolType) Matches(o execute.Type) bool {
	_, ok := o.(*boolType)
	return ok
}

func (t *boolType) New(v execute.Value) (execute.Value, error) {
	return NewBool(v.ToBool()), nil
}

func (t *boolType) String() string {
	return "bool"
}

var BoolType = &boolType{}

// -------------------------------------------------------------------------------------------------

type bytesType struct{}

func (t *bytesType) IsNumeric() bool {
	return true
}

func (t *bytesType) Matches(o execute.Type) bool {
	_, ok := o.(*bytesType)
	return ok
}

func (t *bytesType) New(v execute.Value) (execute.Value, error) {
	vc, err := v.ToBytes()
	if err != nil {
		return nil, err
	}
	return NewBytes(vc), nil
}

func (t *bytesType) String() string {
	return "bytes"
}

var BytesType = &bytesType{}

// -------------------------------------------------------------------------------------------------

type floatType struct{}

func (t *floatType) IsNumeric() bool {
	return true
}

func (t *floatType) Matches(o execute.Type) bool {
	_, ok := o.(*floatType)
	return ok
}

func (t *floatType) New(v execute.Value) (execute.Value, error) {
	vc, err := v.ToFloat()
	if err != nil {
		return nil, err
	}
	return NewFloat(vc), nil
}

func (t *floatType) String() string {
	return "float"
}

var FloatType = &floatType{}

// -------------------------------------------------------------------------------------------------

type funcType struct{}

func (t *funcType) IsNumeric() bool {
	return false
}

func (t *funcType) Matches(o execute.Type) bool {
	_, ok := o.(*funcType)
	return ok
}

func (t *funcType) New(v execute.Value) (execute.Value, error) {
	panic("funcType.New() is not supported")
}

func (t *funcType) String() string {
	return "func"
}

var FuncType = &funcType{}

// -------------------------------------------------------------------------------------------------

type generatorType struct{}

func (t *generatorType) IsNumeric() bool {
	return false
}

func (t *generatorType) Matches(o execute.Type) bool {
	_, ok := o.(*generatorType)
	return ok
}

func (t *generatorType) New(v execute.Value) (execute.Value, error) {
	panic("generatorType.New() is not supported")
}

func (t *generatorType) String() string {
	return "generator"
}

var GeneratorType = &generatorType{}

// -------------------------------------------------------------------------------------------------

type intType struct{}

func (t *intType) IsNumeric() bool {
	return true
}

func (t *intType) Matches(o execute.Type) bool {
	_, ok := o.(*intType)
	return ok
}

func (t *intType) New(v execute.Value) (execute.Value, error) {
	vc, err := v.ToInt()
	if err != nil {
		return nil, err
	}
	return NewInt(vc), nil
}

func (t *intType) String() string {
	return "int"
}

var IntType = &intType{}

// -------------------------------------------------------------------------------------------------

type iteratorType struct{}

func (t *iteratorType) IsNumeric() bool {
	return false
}

func (t *iteratorType) Matches(o execute.Type) bool {
	_, ok := o.(*iteratorType)
	return ok
}

func (t *iteratorType) New(v execute.Value) (execute.Value, error) {
	panic("iteratorType.New() is not supported")
}

func (t *iteratorType) String() string {
	return "iterator"
}

var IteratorType = &iteratorType{}

// -------------------------------------------------------------------------------------------------

type listType struct{}

func (t *listType) IsNumeric() bool {
	return false
}

func (t *listType) Matches(o execute.Type) bool {
	_, ok := o.(*listType)
	return ok
}

func (t *listType) New(v execute.Value) (execute.Value, error) {
	if v == nil {
		return NewList(nil), nil
	}
	panic("listType.New() is not supported with a non-nil argument")
}

func (t *listType) String() string {
	return "list"
}

var ListType = &listType{}

// -------------------------------------------------------------------------------------------------

type mapType struct{}

func (t *mapType) IsNumeric() bool {
	return false
}

func (t *mapType) Matches(o execute.Type) bool {
	_, ok := o.(*mapType)
	return ok
}

func (t *mapType) New(v execute.Value) (execute.Value, error) {
	if v == nil {
		return NewMap(), nil
	}
	panic("mapType.New() is not supported with a non-nil argument")
}

func (t *mapType) String() string {
	return "map"
}

var MapType = &mapType{}

// -------------------------------------------------------------------------------------------------

type moduleType struct{}

func (t *moduleType) IsNumeric() bool {
	return false
}

func (t *moduleType) Matches(o execute.Type) bool {
	_, ok := o.(*moduleType)
	return ok
}

func (t *moduleType) New(v execute.Value) (execute.Value, error) {
	panic("moduleType.New() is not supported")
}

func (t *moduleType) String() string {
	return "module"
}

var ModuleType = &moduleType{}

// -------------------------------------------------------------------------------------------------

type nullType struct{}

func (t *nullType) IsNumeric() bool {
	return false
}

func (t *nullType) Matches(o execute.Type) bool {
	_, ok := o.(*nullType)
	return ok
}

func (t *nullType) New(v execute.Value) (execute.Value, error) {
	panic("nullType.New() is not supported")
}

func (t *nullType) String() string {
	return "null"
}

var NullType = &nullType{}

// -------------------------------------------------------------------------------------------------

type strType struct{}

func (t *strType) IsNumeric() bool {
	return false
}

func (t *strType) Matches(o execute.Type) bool {
	_, ok := o.(*strType)
	return ok
}

func (t *strType) New(v execute.Value) (execute.Value, error) {
	vc, err := v.ToStr()
	if err != nil {
		return nil, err
	}
	return NewStr(vc), nil
}

func (t *strType) String() string {
	return "str"
}

var StrType = &strType{}

// -------------------------------------------------------------------------------------------------

type uintType struct{}

func (t *uintType) IsNumeric() bool {
	return true
}

func (t *uintType) Matches(o execute.Type) bool {
	_, ok := o.(*uintType)
	return ok
}

func (t *uintType) New(v execute.Value) (execute.Value, error) {
	vc, err := v.ToUint()
	if err != nil {
		return nil, err
	}
	return NewUint(vc), nil
}

func (t *uintType) String() string {
	return "uint"
}

var UintType = &uintType{}
