package lib

import (
	"errors"
	"math"
	"reflect"
)

type Alternative struct {
	grade           []Evaluated
	countOfCriteria int
}

func (a *Alternative) String() string {
	s := "[ "
	for i := 0; i < a.countOfCriteria; i++ {
		s += a.grade[i].String() + " "
	}
	return s + " ]"
}

func (a *Alternative) NumberMetric(to *Alternative, v Variants) (Number, error) {
	if a.countOfCriteria != to.countOfCriteria {
		return 0, InvalidSize
	}

	result := Number(0)

	for i := 0; i < a.countOfCriteria; i++ {
		if tmp, err := a.grade[i].DiffNumber(to.grade[i].ConvertToNumber(), v); err != nil {
			return 0, errors.Join(err)
		} else {
			result += tmp
		}
	}

	if v == SqrtDistance {
		result = Number(math.Sqrt(float64(result)))
	} else {
		result = Number(math.Cbrt(float64(result)))
	}

	return result, nil
}

func (a *Alternative) IntervalMetric(to *Alternative, criteria []Criterion,
	idealType bool, v Variants) (Interval, error) {
	if a.countOfCriteria != to.countOfCriteria {
		return Interval{}, InvalidSize
	}

	var result Evaluated
	result = Interval{Number(0), Number(0)}

	for i, c := range criteria {
		var typeOfCriterion bool
		if idealType == Positive {
			typeOfCriterion = c.typeOfCriteria
		} else {
			typeOfCriterion = !c.typeOfCriteria
		}
		if tmp, err := a.grade[i].DiffInterval(to.grade[i].ConvertToInterval(), typeOfCriterion, v); err != nil {
			return Interval{}, errors.Join(err)
		} else {
			result = result.Sum(tmp)
		}
	}

	if v == SqrtDistance {
		result = Interval{Number(math.Sqrt(float64(result.ConvertToInterval().Start))),
			Number(math.Sqrt(float64(result.ConvertToInterval().End)))}
	} else {
		result = Interval{Number(math.Cbrt(float64(result.ConvertToInterval().Start))),
			Number(math.Cbrt(float64(result.ConvertToInterval().End)))}
	}

	return result.ConvertToInterval(), nil
}

func (a *Alternative) FSMetric(to *Alternative, v Variants) (Number, error) {
	if a.countOfCriteria != to.countOfCriteria {
		return 0, InvalidSize
	}

	result := Number(0)

	for i := 0; i < a.countOfCriteria; i++ {
		if tmp, err := a.grade[i].DiffNumber(to.grade[i], v); err != nil {
			return 0, errors.Join(err)
		} else {
			result += tmp
		}
	}

	if v == SqrtDistance {
		result = Number(math.Sqrt(float64(result)))
	} else {
		result = Number(math.Cbrt(float64(result)))
	}

	return result, nil
}

func (a *Alternative) Sum() Evaluated {
	var sum Evaluated

	if reflect.TypeOf(a.grade[0]) == reflect.TypeOf(NumbersMin) {
		sum = Number(0.0)
	} else if reflect.TypeOf(a.grade[0]) == reflect.TypeOf(Interval{}) {
		sum = Interval{0.0, 0.0}
	} else if reflect.TypeOf(a.grade[0]) == reflect.TypeOf(&T1FS{}) {
		sum = NewT1FS(0.0, 0.0, 0.0)
	}

	for _, el := range a.grade {
		sum = sum.Sum(el)
	}
	return sum
}
