package builtins

import (
	"fmt"

	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/types"
)

func rangeImpl(args ...execute.Value) (execute.Value, error) {
	if len(args) == 0 {
		return nil, errors.CallError("range", len(args), 1)
	}
	var lower execute.Value
	var step execute.Value
	upper := args[0]
	if len(args) > 1 {
		lower = upper
		upper = args[1]
	} else {
		lower = types.NewUint(0)
	}
	if len(args) > 2 {
		step = args[2]
	} else {
		step = types.NewUint(1)
	}
	if len(args) > 3 {
		return nil, errors.CallError("range", len(args), 3)
	}
	rg, err := newRangeGenerator(lower, upper, step)
	if err != nil {
		return nil, err
	}
	return types.NewGenerator(rg), nil
}

type rangeGenerator struct {
	valueType execute.Type

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

func newRangeGenerator(start, stop, step execute.Value) (*rangeGenerator, error) {
	for _, v := range []execute.Value{start, stop, step} {
		if !v.Type().IsNumeric() {
			return nil, errors.TypeErrorFromMessage(fmt.Sprintf("range cannot be called with non-numeric values: %q", v.Type()))
		}
	}
	commonType, ok := types.CommonNumericType(start.Type(), stop.Type())
	if !ok {
		panic("types.CommonNumericType() returned no common type in newRangeGenerator()")
	}
	commonType, ok = types.CommonNumericType(commonType, step.Type())
	if !ok {
		panic("types.CommonNumericType() returned no common type in newRangeGenerator()")
	}
	switch commonType {
	case types.FloatType:
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
		return &rangeGenerator{
			valueType: types.FloatType,
			nextF:     startC,
			startF:    startC,
			stopF:     stopC,
			stepF:     stepC,
		}, nil
	case types.IntType:
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
		return &rangeGenerator{
			valueType: types.IntType,
			nextI:     startC,
			startI:    startC,
			stopI:     stopC,
			stepI:     stepC,
		}, nil
	case types.BoolType:
		fallthrough
	case types.UintType:
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
		return &rangeGenerator{
			valueType: types.UintType,
			nextU:     startC,
			startU:    startC,
			stopU:     stopC,
			stepU:     stepC,
		}, nil
	}
	panic("unexpected commonType in newRangeGenerator()")
}

func (g *rangeGenerator) HasNext() bool {
	switch g.valueType {
	case types.FloatType:
		return g.nextF < g.stopF
	case types.IntType:
		return g.nextI < g.stopI
	case types.UintType:
		return g.nextU < g.stopU
	}
	panic("unexpected valueType in rangeGenerator.HasNext()")
}

func (g *rangeGenerator) Next() (execute.Value, error) {
	switch g.valueType {
	case types.FloatType:
		var curr float64
		curr, g.nextF = g.nextF, g.nextF+g.stepF
		return types.NewFloat(curr), nil
	case types.IntType:
		var curr int64
		curr, g.nextI = g.nextI, g.nextI+g.stepI
		return types.NewInt(curr), nil
	case types.UintType:
		var curr uint64
		curr, g.nextU = g.nextU, g.nextU+g.stepU
		return types.NewUint(curr), nil
	}
	panic("unexpected type in rangeGenerator.Next()")
}
