package lib

import (
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

func (i Interval) ConvertToNumber() Number {
	return (i.Start + i.End) / 2
}

func (i Interval) ConvertToInterval() Interval {
	return i
}

func (i Interval) ConvertToT1FS(v Variants) *T1FS {
	if v == Default || v == Triangle {
		return NewT1FS(i.Start, i.ConvertToNumber(), i.End)
	} else {
		return NewT1FS(i.Start, i.Start, i.End, i.End)
	}
}

func (i Interval) ConvertToIT2FS(v Variants) *IT2FS {
	if v == Default || v == Triangle {
		return NewIT2FS([]Interval{{i.Start, i.Start}, {i.End, i.End}}, []Number{i.ConvertToNumber()})
	} else {
		return NewIT2FS([]Interval{{i.Start, i.Start}, {i.End, i.End}}, []Number{i.Start, i.End})
	}
}

func (i Interval) Weighted(weight Evaluated) Evaluated {
	if reflect.TypeOf(weight) == reflect.TypeOf(Number(0)) || reflect.TypeOf(weight) == reflect.TypeOf(i) {
		return Interval{i.Start * weight.ConvertToNumber(), i.End * weight.ConvertToNumber()}
	}
	return nil
}

func (i Interval) DiffNumber(other Evaluated, v Variants) (Number, error) {
	if reflect.TypeOf(other) != reflect.TypeOf(NumbersMax) {
		return 0, IncompatibleTypes
	}

	if v == SqrtDistance {
		return (i.Start-other.ConvertToNumber())*(i.Start-other.ConvertToNumber()) +
			(i.End-other.ConvertToNumber())*(i.End-other.ConvertToNumber()), nil
	} else if v == CbrtDistance {
		return Number(math.Abs(math.Pow(float64(i.Start-other.ConvertToNumber()), 3))) +
			Number(math.Abs(math.Pow(float64(i.End-other.ConvertToNumber()), 3))), nil
	} else {
		return NumbersMin, InvalidCaseOfOperation
	}
}

func (i Interval) DiffInterval(other Interval, typeOfCriterion bool, v Variants) (Interval, error) {
	ret := Interval{}
	if typeOfCriterion == Benefit {
		if v == SqrtDistance {
			ret.Start = (i.End - other.Start) * (i.End - other.Start)
			ret.End = (i.Start - other.End) * (i.Start - other.End)
		} else if v == CbrtDistance {
			ret.Start = Number(math.Abs(math.Pow(float64(i.End-other.Start), 3)))
			ret.End = Number(math.Abs(math.Pow(float64(i.Start-other.End), 3)))
		} else {
			return Interval{}, InvalidCaseOfOperation
		}
	} else if typeOfCriterion == Cost {
		if v == SqrtDistance {
			ret.Start = (i.Start - other.Start) * (i.Start - other.Start)
			ret.End = (i.End - other.End) * (i.End - other.End)
		} else if v == CbrtDistance {
			ret.Start = Number(math.Abs(math.Pow(float64(i.Start-other.Start), 3)))
			ret.End = Number(math.Abs(math.Pow(float64(i.End-other.End), 3)))
		} else {
			return Interval{}, InvalidCaseOfOperation
		}
	} else {
		return Interval{}, InvalidCaseOfOperation
	}
	return ret, nil
}

func (i Interval) Sum(other Evaluated) Evaluated {
	return Interval{i.Start + other.ConvertToInterval().Start,
		i.End + other.ConvertToInterval().End}
}

func (i Interval) SenguptaGeq(other Interval) bool {
	if math.Abs(float64(i.ConvertToNumber()-other.ConvertToNumber())) < 1e-3 {
		return i.End-i.Start < other.End-other.Start
	} else {
		return i.ConvertToNumber() > other.ConvertToNumber()
	}
}

func positiveIdealRateInterval(alts []Alternative, criteria []Criterion) (*Alternative, error) {
	positive := &Alternative{make([]Evaluated, len(alts[0].grade)), len(criteria)}

	for i, c := range criteria {
		for j := range alts {
			if reflect.TypeOf(alts[j].grade[i]) != reflect.TypeOf(Interval{}) {
				return nil, IncompatibleTypes
			}

			if positive.grade[i] == nil {
				positive.grade[i] = alts[j].grade[i]
				continue
			}

			if c.typeOfCriteria == Benefit {
				positive.grade[i] = Max(positive.grade[i], alts[j].grade[i])
			} else {
				positive.grade[i] = Min(positive.grade[i], alts[j].grade[i])
			}
		}
	}

	return positive, nil
}

func negativeIdealRateInterval(alts []Alternative, criteria []Criterion) (*Alternative, error) {
	return positiveIdealRateInterval(alts, ChangeTypes(criteria))
}
