package matrix

import (
	"errors"
	"math"
)

type Alternative struct {
	grade           []Evaluated
	countOfCriteria int
}

func (a *Alternative) String() string {
	s := ""
	for i := 0; i < a.countOfCriteria; i++ {
		s += a.grade[i].String() + " "
	}
	return s
}

func (a *Alternative) FindDistanceNumber(to *Alternative, v Variants) (Number, error) {
	if a.countOfCriteria != to.countOfCriteria {
		return 0, errors.New("incomparable alternatives")
	}

	result := Number(0)

	for i := 0; i < a.countOfCriteria; i++ {
		if tmp, err := a.grade[i].DistanceNumber(to.grade[i].ConvertToNumbers(), v); err != nil {
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

func (a *Alternative) FindDistanceInterval(to *Alternative, criteria []Criterion,
	idealType bool, v Variants) (Interval, error) {
	if a.countOfCriteria != to.countOfCriteria {
		return Interval{}, errors.New("incomparable alternatives")
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
		if tmp, err := a.grade[i].DistanceInterval(to.grade[i].ConvertToInterval(), typeOfCriterion, v); err != nil {
			return Interval{}, errors.Join(err)
		} else {
			if result, err = result.Sum(tmp); err != nil {
				return Interval{}, errors.Join(err)
			}
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

func (a *Alternative) FindDistanceT1FS(to *Alternative, v Variants) (Number, error) {
	if a.countOfCriteria != to.countOfCriteria {
		return 0, errors.New("incomparable alternatives")
	}

	result := Number(0)

	for i := 0; i < a.countOfCriteria; i++ {
		if tmp, err := a.grade[i].DistanceNumber(to.grade[i].ConvertToT1FS(Default), v); err != nil {
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
