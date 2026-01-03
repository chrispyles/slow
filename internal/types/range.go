package types

import (
	"fmt"

	"github.com/chrispyles/slow/internal/errors"
	"github.com/chrispyles/slow/internal/execute"
)

type RangeIterator struct {
	valueType        execute.Type
	incr             bool
	truncToContainer bool

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
	var startType, stopType, stepType execute.Type = UintType, UintType, UintType
	if start != nil {
		startType = start.Type()
	}
	if stop != nil {
		stopType = stop.Type()
	}
	if step == nil {
		step = NewUint(1)
	}
	stepType = step.Type()
	for _, v := range []execute.Type{startType, stopType, stepType} {
		if !v.IsNumeric() {
			return nil, errors.TypeErrorFromMessage(fmt.Sprintf("ranges cannot be created with non-numeric values: %q", v))
		}
	}
	commonType, ok := CommonNumericType(startType, stopType)
	if !ok {
		panic("CommonNumericType() returned no common type in newRangeGenerator()")
	}
	commonType, ok = CommonNumericType(commonType, stepType)
	if !ok {
		panic("CommonNumericType() returned no common type in newRangeGenerator()")
	}
	incr := must(step.ToFloat()) > 0 // all numeric types convert to float w/o losing sign
	var truncToContainer bool
	if incr {
		truncToContainer = stop == nil
	} else {
		truncToContainer = start == nil
	}
	if start == nil {
		start = NewUint(0)
	}
	if stop == nil {
		stop = NewUint(0)
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
			valueType:        FloatType,
			incr:             incr,
			truncToContainer: truncToContainer,
			nextF:            startC,
			startF:           startC,
			stopF:            stopC,
			stepF:            stepC,
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
			valueType:        IntType,
			incr:             incr,
			truncToContainer: truncToContainer,
			nextI:            startC,
			startI:           startC,
			stopI:            stopC,
			stepI:            stepC,
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
			valueType:        UintType,
			incr:             incr,
			truncToContainer: truncToContainer,
			nextU:            startC,
			startU:           startC,
			stopU:            stopC,
			stepU:            stepC,
		}, nil
	}
	panic("unexpected commonType in newRangeGenerator()")
}

func (g *RangeIterator) HasNext() bool {
	// TODO: panic if truncToContainer is true????
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
	if g.truncToContainer {
		return nil, errors.NewValueError("a range without endpoints may only be used for indexing")
	}
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

func (g *RangeIterator) WithContainerLen(l uint64) *Generator {
	copy := *g
	if !copy.incr && l > 0 {
		// If decrementing, set the starting index to len - 1 so we don't start w/ an out-of-bounds index.
		l -= 1
	}
	if copy.truncToContainer {
		switch copy.valueType {
		case FloatType:
			if copy.startF != copy.nextF {
				panic("WithContainerLen called after iteration started")
			}
			if copy.incr {
				copy.stopF = float64(l)
			} else {
				copy.startF = float64(l)
				copy.nextF = float64(l)
			}
		case IntType:
			if copy.startI != copy.nextI {
				panic("WithContainerLen called after iteration started")
			}
			if copy.incr {
				copy.stopI = int64(l)
			} else {
				copy.startI = int64(l)
				copy.nextI = int64(l)
			}
		case UintType:
			if copy.startU != copy.nextU {
				panic("WithContainerLen called after iteration started")
			}
			if copy.incr {
				copy.stopU = l
			} else {
				copy.startU = l
				copy.nextU = l
			}
		}
		copy.truncToContainer = false
	}
	return NewGenerator(&copy)
}
