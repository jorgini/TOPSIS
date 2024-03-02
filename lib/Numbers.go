package lib

import (
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

func (n Number) ConvertToNumber() Number {
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

func (n Number) ConvertToIT2FS(v Variants) *IT2FS {
	if v == Default || v == Triangle {
		return NewIT2FS([]Interval{{n, n}, {n, n}}, []Number{n})
	} else {
		return NewIT2FS([]Interval{{n, n}, {n, n}}, []Number{n, n})
	}
}

func (n Number) Weighted(weight Evaluated) Evaluated {
	if reflect.TypeOf(weight) == reflect.TypeOf(Number(0)) || reflect.TypeOf(weight) == reflect.TypeOf(Interval{}) {
		return n * weight.ConvertToNumber()
	}
	return nil
}

func (n Number) DiffNumber(other Evaluated, v Variants) (Number, error) {
	if reflect.TypeOf(other) != reflect.TypeOf(n) {
		return 0, IncompatibleTypes
	}

	if v == SqrtDistance {
		return (n - other.ConvertToNumber()) * (n - other.ConvertToNumber()), nil
	} else if v == CbrtDistance {
		return Number(math.Abs(math.Pow(float64(n-other.ConvertToNumber()), 3))), nil
	} else {
		return 0, InvalidCaseOfOperation
	}
}

func (n Number) DiffInterval(other Interval, typeOfCriterion bool, v Variants) (Interval, error) {
	return Interval{}, NoUsageMethod
}

func (n Number) Sum(other Evaluated) Evaluated {
	return n + other.ConvertToNumber()
}

func positiveIdealRateNumber(alts []Alternative, criteria []Criterion) (*Alternative, error) {
	positive := &Alternative{make([]Evaluated, len(alts[0].grade)), len(criteria)}

	for i, c := range criteria {
		for j := range alts {
			if reflect.TypeOf(alts[j].grade[i]) != reflect.TypeOf(NumbersMin) &&
				reflect.TypeOf(alts[j].grade[i]) != reflect.TypeOf(Interval{}) {
				return nil, IncompatibleTypes
			}

			if positive.grade[i] == nil && c.typeOfCriteria == Benefit {
				positive.grade[i] = NumbersMin
			} else if positive.grade[i] == nil && c.typeOfCriteria == Cost {
				positive.grade[i] = NumbersMax
			}

			if reflect.TypeOf(alts[j].grade[i]) == reflect.TypeOf(Interval{}) {
				if c.typeOfCriteria == Benefit {
					positive.grade[i] = Max(positive.grade[i].ConvertToNumber(), alts[j].grade[i].ConvertToInterval().End)
				} else {
					positive.grade[i] = Min(positive.grade[i].ConvertToNumber(), alts[j].grade[i].ConvertToInterval().Start)
				}
			} else if reflect.TypeOf(alts[j].grade[i]) == reflect.TypeOf(NumbersMin) {
				if c.typeOfCriteria == Benefit {
					positive.grade[i] = Max(positive.grade[i], alts[j].grade[i])
				} else {
					positive.grade[i] = Min(positive.grade[i], alts[j].grade[i])
				}
			}
		}
	}
	return positive, nil
}

func negativeIdealRateNumber(alts []Alternative, criteria []Criterion) (*Alternative, error) {
	return positiveIdealRateNumber(alts, ChangeTypes(criteria))
}
