package matrix

import (
	"errors"
	"fmt"
	"math"
	"reflect"
)

type Interval struct {
	Start, End Number
}

func (i Interval) String() string {
	return fmt.Sprintf("[%.2f, %.2f]", i.Start, i.End)
}

func (i Interval) ConvertToNumbers() Number {
	return (i.Start + i.End) / 2
}

func (i Interval) ConvertToInterval() Interval {
	return i
}

func (i Interval) ConvertToT1FS(v Variants) *T1FS {
	if v == Default || v == Triangle {
		return NewT1FS(i.Start, i.ConvertToNumbers(), i.End)
	} else {
		return NewT1FS(i.Start, i.Start, i.End, i.End)
	}
}

func (i Interval) Weighted(weight Evaluated) Evaluated {
	if reflect.TypeOf(weight) == reflect.TypeOf(Number(0)) || reflect.TypeOf(weight) == reflect.TypeOf(i) {
		return Interval{i.Start * weight.ConvertToNumbers(), i.End * weight.ConvertToNumbers()}
	}
	return nil
}

func (i Interval) DistanceNumber(other Evaluated, v Variants) (Number, error) {
	if reflect.TypeOf(other) != reflect.TypeOf(NumbersMax) {
		return 0, errors.New("incomparable values")
	}

	if v == SqrtDistance {
		return (i.Start-other.ConvertToNumbers())*(i.Start-other.ConvertToNumbers()) +
			(i.End-other.ConvertToNumbers())*(i.End-other.ConvertToNumbers()), nil
	} else if v == CbrtDistance {
		return Number(math.Abs(math.Pow(float64(i.Start-other.ConvertToNumbers()), 3))) +
			Number(math.Abs(math.Pow(float64(i.End-other.ConvertToNumbers()), 3))), nil
	} else {
		return NumbersMin, errors.New("incomplete case of distance")
	}
}

func (i Interval) DistanceInterval(other Interval, typeOfCriterion bool, v Variants) (Interval, error) {
	ret := Interval{}
	if typeOfCriterion == Benefit {
		if v == SqrtDistance {
			ret.Start = (i.End - other.Start) * (i.End - other.Start)
			ret.End = (i.Start - other.End) * (i.Start - other.End)
		} else if v == CbrtDistance {
			ret.Start = Number(math.Abs(math.Pow(float64(i.End-other.Start), 3)))
			ret.End = Number(math.Abs(math.Pow(float64(i.Start-other.End), 3)))
		} else {
			return Interval{}, errors.New("incomplete case of distance")
		}
	} else if typeOfCriterion == Cost {
		if v == SqrtDistance {
			ret.Start = (i.Start - other.Start) * (i.Start - other.Start)
			ret.End = (i.End - other.End) * (i.End - other.End)
		} else if v == CbrtDistance {
			ret.Start = Number(math.Abs(math.Pow(float64(i.Start-other.Start), 3)))
			ret.End = Number(math.Abs(math.Pow(float64(i.End-other.End), 3)))
		} else {
			return Interval{}, errors.New("incomplete case of distance")
		}
	} else {
		return Interval{}, errors.New("invalid type of criterion")
	}
	return ret, nil
}

func (i Interval) Sum(other Evaluated) (Evaluated, error) {
	if reflect.TypeOf(other) != reflect.TypeOf(i) {
		return nil, errors.New("can't sum number and interval")
	} else {
		return Interval{i.Start + other.ConvertToInterval().Start,
			i.End + other.ConvertToInterval().End}, nil
	}
}

func (i Interval) SenguptaGeq(other Interval) bool {
	if math.Abs(float64(i.ConvertToNumbers()-other.ConvertToNumbers())) < 1e-3 {
		return i.End-i.Start < other.End-other.Start
	} else {
		return i.ConvertToNumbers() > other.ConvertToNumbers()
	}
}

func positiveIdealRateInterval(x, y Interval, typeOfCriterion bool) Interval {
	ideal := Interval{}
	if typeOfCriterion == Benefit {
		if x.End >= y.End {
			ideal.End = x.End
			ideal.Start = Number(math.Max(float64(x.Start), float64(y.End)))
		} else {
			ideal.End = y.End
			ideal.Start = Number(math.Max(float64(x.End), float64(y.Start)))
		}
	} else {
		if x.Start <= y.Start {
			ideal.End = x.Start
			ideal.Start = Number(math.Min(float64(x.End), float64(y.Start)))
		} else {
			ideal.End = y.Start
			ideal.Start = Number(math.Min(float64(x.Start), float64(y.End)))
		}
	}
	return ideal
}

func negativeIdealRateInterval(x, y Interval, typeOfCriterion bool) Interval {
	return positiveIdealRateInterval(x, y, !typeOfCriterion)
}
