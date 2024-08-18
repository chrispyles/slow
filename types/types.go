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

func (t *boolType) String() string {
	return "bool"
}

var BoolType = &boolType{}

// -------------------------------------------------------------------------------------------------

type floatType struct{}

func (t *floatType) IsNumeric() bool {
	return true
}

func (t *floatType) Matches(o execute.Type) bool {
	_, ok := o.(*floatType)
	return ok
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

func (t *listType) String() string {
	return "list"
}

var ListType = &listType{}

// -------------------------------------------------------------------------------------------------

type nullType struct{}

func (t *nullType) IsNumeric() bool {
	return false
}

func (t *nullType) Matches(o execute.Type) bool {
	_, ok := o.(*nullType)
	return ok
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

func (t *uintType) String() string {
	return "uint"
}

var UintType = &uintType{}
