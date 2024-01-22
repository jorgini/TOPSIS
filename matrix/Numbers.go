package matrix

import (
	"errors"
	"fmt"
	"math"
	"reflect"
)

type Number float64

const NumbersMax Number = 2147483647

const NumbersMin Number = -2147483647

func (n Number) String() string {
	return fmt.Sprintf("%.2f", float64(n))
}

func (n Number) ConvertToNumbers() Number {
	return n
}

func (n Number) ConvertToInterval() Interval {
	return Interval{n, n}
}

func (n Number) ConvertToT1FS(v Variants) *T1FS {
	if v == Default || v == Triangle {
		return NewT1FS(n, n, n)
	} else {
		return NewT1FS(n, n, n, n)
	}
}

func (n Number) Weighted(weight Evaluated) Evaluated {
	if reflect.TypeOf(weight) == reflect.TypeOf(Number(0)) || reflect.TypeOf(weight) == reflect.TypeOf(Interval{}) {
		return n * weight.ConvertToNumbers()
	}
	return nil
}

func (n Number) DistanceNumber(other Evaluated, v Variants) (Number, error) {
	if reflect.TypeOf(other) != reflect.TypeOf(n) {
		return 0, errors.New("incomparable values")
	}

	if v == SqrtDistance {
		return (n - other.ConvertToNumbers()) * (n - other.ConvertToNumbers()), nil
	} else if v == CbrtDistance {
		return Number(math.Abs(math.Pow(float64(n-other.ConvertToNumbers()), 3))), nil
	} else {
		return 0, errors.New("incomplete case of distance")
	}
}

func (n Number) DistanceInterval(other Interval, typeOfCriterion bool, v Variants) (Interval, error) {
	return Interval{}, errors.New("not usage method")
}

func (n Number) Sum(other Evaluated) (Evaluated, error) {
	if reflect.TypeOf(other) != reflect.TypeOf(n) {
		return nil, errors.New("can't sum number and interval")
	} else {
		n += other.ConvertToNumbers()
		return n, nil
	}
}

func positiveIdealRateNumber(x, y Evaluated, typeOfCriterion bool) Number {
	a := NumbersMin
	b := NumbersMin
	if reflect.TypeOf(x) == reflect.TypeOf(Interval{}) {
		tmp := x.ConvertToInterval()
		if typeOfCriterion == Benefit {
			a = tmp.End
		} else {
			a = tmp.Start
		}
	} else if reflect.TypeOf(x) == reflect.TypeOf(&T1FS{}) {
		tmp := x.ConvertToT1FS(Default)
		if typeOfCriterion == Benefit {
			a = tmp.vert[len(tmp.vert)-1]
		} else {
			a = tmp.vert[0]
		}
	}

	if reflect.TypeOf(y) == reflect.TypeOf(Interval{}) {
		tmp := y.ConvertToInterval()
		if typeOfCriterion == Benefit {
			b = tmp.End
		} else {
			b = tmp.Start
		}
	} else if reflect.TypeOf(x) == reflect.TypeOf(&T1FS{}) {
		tmp := y.ConvertToT1FS(Default)
		if typeOfCriterion == Benefit {
			b = tmp.vert[len(tmp.vert)-1]
		} else {
			b = tmp.vert[0]
		}
	}

	if typeOfCriterion == Benefit {
		return Number(math.Max(float64(a), float64(b)))
	}
	return Number(math.Min(float64(a), float64(b)))
}

func negativeIdealRateNumber(x, y Evaluated, typeOfCriterion bool) Number {
	return positiveIdealRateNumber(x, y, !typeOfCriterion)
}
