package types

import (
	"fmt"

	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
)

type RangeIterator struct {
	valueType execute.Type
	incr      bool

	nextF  float64
	startF float64
	stopF  float64
	stepF  float64

	nextI  int64
	startI int64
	stopI  int64
	stepI  int64

	nextU  uint64
	startU uint64
	stopU  uint64
	stepU  uint64
}

func NewRangeGenerator(start, stop, step execute.Value) (execute.Value, error) {
	rg, err := makeRange(start, stop, step)
	if err != nil {
		return nil, err
	}
	return NewGenerator(rg), nil
}

func makeRange(start, stop, step execute.Value) (*RangeIterator, error) {
	for _, v := range []execute.Value{start, stop, step} {
		if !v.Type().IsNumeric() {
			return nil, errors.TypeErrorFromMessage(fmt.Sprintf("range cannot be called with non-numeric values: %q", v.Type()))
		}
	}
	commonType, ok := CommonNumericType(start.Type(), stop.Type())
	if !ok {
		panic("CommonNumericType() returned no common type in newRangeGenerator()")
	}
	commonType, ok = CommonNumericType(commonType, step.Type())
	if !ok {
		panic("CommonNumericType() returned no common type in newRangeGenerator()")
	}
	switch commonType {
	case FloatType:
		startC, err := start.ToFloat()
		if err != nil {
			return nil, err
		}
		stopC, err := stop.ToFloat()
		if err != nil {
			return nil, err
		}
		stepC, err := step.ToFloat()
		if err != nil {
			return nil, err
		}
		return &RangeIterator{
			valueType: FloatType,
			incr:      stepC >= 0,
			nextF:     startC,
			startF:    startC,
			stopF:     stopC,
			stepF:     stepC,
		}, nil
	case IntType:
		startC, err := start.ToInt()
		if err != nil {
			return nil, err
		}
		stopC, err := stop.ToInt()
		if err != nil {
			return nil, err
		}
		stepC, err := step.ToInt()
		if err != nil {
			return nil, err
		}
		return &RangeIterator{
			valueType: IntType,
			incr:      stepC >= 0,
			nextI:     startC,
			startI:    startC,
			stopI:     stopC,
			stepI:     stepC,
		}, nil
	case BoolType:
		fallthrough
	case UintType:
		startC, err := start.ToUint()
		if err != nil {
			return nil, err
		}
		stopC, err := stop.ToUint()
		if err != nil {
			return nil, err
		}
		stepC, err := step.ToUint()
		if err != nil {
			return nil, err
		}
		return &RangeIterator{
			valueType: UintType,
			incr:      stepC >= 0,
			nextU:     startC,
			startU:    startC,
			stopU:     stopC,
			stepU:     stepC,
		}, nil
	}
	panic("unexpected commonType in newRangeGenerator()")
}

func (g *RangeIterator) HasNext() bool {
	switch g.valueType {
	case FloatType:
		if g.incr {
			return g.nextF < g.stopF
		} else {
			return g.nextF > g.stopF
		}
	case IntType:
		if g.incr {
			return g.nextI < g.stopI
		} else {
			return g.nextI > g.stopI
		}
	case UintType:
		if g.incr {
			return g.nextU < g.stopU
		} else {
			return g.nextU > g.stopU
		}
	}
	panic("unexpected valueType in rangeGenerator.HasNext()")
}

func (g *RangeIterator) Next() (execute.Value, error) {
	switch g.valueType {
	case FloatType:
		var curr float64
		curr, g.nextF = g.nextF, g.nextF+g.stepF
		return NewFloat(curr), nil
	case IntType:
		var curr int64
		curr, g.nextI = g.nextI, g.nextI+g.stepI
		return NewInt(curr), nil
	case UintType:
		var curr uint64
		curr, g.nextU = g.nextU, g.nextU+g.stepU
		return NewUint(curr), nil
	}
	panic("unexpected type in rangeGenerator.Next()")
}
