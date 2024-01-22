package matrix

import (
	"errors"
	"reflect"
)

type Variants int

const (
	NormalizeWeightsByMidPoint = 0
	NormalizeWithSum           = 1
	NormalizeWithMax           = 2
	Sengupta                   = 3
	AlphaSlices                = 4
	AggregateRatings           = 5
	AggregateDistances         = 6
	SqrtDistance               = 8
	CbrtDistance               = 9
	Positive                   = true
	Negative                   = false
)

func (m *Matrix) CalcWeightedMatrix() {
	for j := 0; j < m.countCriteria; j++ {
		for i := 0; i < m.countAlternatives; i++ {
			m.data[i].grade[j] = m.data[i].grade[j].Weighted(m.criteria[j].weight)
		}
	}
}

func (m *Matrix) FindIdeals(v Variants) {
	positive := Alternative{make([]Evaluated, m.countCriteria), m.countCriteria}
	negative := Alternative{make([]Evaluated, m.countCriteria), m.countCriteria}

	for i := 0; i < m.countAlternatives; i++ {
		for j := 0; j < m.countCriteria; j++ {
			if positive.grade[j] == nil {
				positive.grade[j] = m.data[i].grade[j]
			} else {
				if reflect.TypeOf(m.data[i].grade[j]) == reflect.TypeOf(Interval{}) && v == Sengupta {
					positive.grade[j] = positiveIdealRateInterval(positive.grade[j].ConvertToInterval(),
						m.data[i].grade[j].ConvertToInterval(), m.criteria[j].typeOfCriteria)
				} else if reflect.TypeOf(m.data[i].grade[j]) == reflect.TypeOf(&T1FS{}) {
					positive.grade[j] = positiveIdealRateT1FS(positive.grade[j].ConvertToT1FS(Default),
						m.data[i].grade[j].ConvertToT1FS(Default), m.criteria[j].typeOfCriteria)
				} else {
					positive.grade[j] = positiveIdealRateNumber(positive.grade[j],
						m.data[i].grade[j], m.criteria[j].typeOfCriteria)
				}
			}

			if negative.grade[j] == nil {
				negative.grade[j] = m.data[i].grade[j]
			} else {
				if reflect.TypeOf(m.data[i].grade[j]) == reflect.TypeOf(Interval{}) && v == Sengupta {
					negative.grade[j] = negativeIdealRateInterval(negative.grade[j].ConvertToInterval(),
						m.data[i].grade[j].ConvertToInterval(), m.criteria[j].typeOfCriteria)
				} else if reflect.TypeOf(m.data[i].grade[j]) == reflect.TypeOf(&T1FS{}) {
					negative.grade[j] = negativeIdealRateT1FS(negative.grade[j].ConvertToT1FS(Default),
						m.data[i].grade[j].ConvertToT1FS(Default), m.criteria[j].typeOfCriteria)
				} else {
					negative.grade[j] = negativeIdealRateNumber(negative.grade[j],
						m.data[i].grade[j], m.criteria[j].typeOfCriteria)
				}
			}
		}
	}

	m.positiveIdeal = &positive
	m.negativeIdeal = &negative
	m.idealsFind = true
}

func (m *Matrix) FindDistanceToIdeals(vt, vi, vn Variants) error {
	for i := 0; i < m.countAlternatives; i++ {
		var err error
		if m.highType != reflect.TypeOf(&T1FS{}) || vt == AlphaSlices {
			if vi == Default {
				if m.distancesToPositive[i], err = m.data[i].FindDistanceNumber(m.positiveIdeal, vn); err != nil {
					return errors.Join(err)
				}

				if m.distancesToNegative[i], err = m.data[i].FindDistanceNumber(m.negativeIdeal, vn); err != nil {
					return errors.Join(err)
				}
			} else if vi == Sengupta {
				if m.distancesToPositive[i], err = m.data[i].FindDistanceInterval(m.positiveIdeal,
					m.criteria, Positive, vn); err != nil {
					return errors.Join(err)
				}

				if m.distancesToNegative[i], err = m.data[i].FindDistanceInterval(m.negativeIdeal,
					m.criteria, Negative, vn); err != nil {
					return errors.Join(err)
				}
			} else {
				return errors.New("incomplete case of calc distance between intervals")
			}
		} else {
			if m.distancesToPositive[i], err = m.data[i].FindDistanceT1FS(m.positiveIdeal, vn); err != nil {
				return errors.Join(err)
			}

			if m.distancesToNegative[i], err = m.data[i].FindDistanceT1FS(m.negativeIdeal, vn); err != nil {
				return errors.Join(err)
			}
		}
	}
	m.distancesFind = true
	return nil
}

func (m *Matrix) CalcCloseness() {
	for i := 0; i < m.countAlternatives; i++ {
		if reflect.TypeOf(m.distancesToPositive[i]) == reflect.TypeOf(Interval{}) {
			neg := m.distancesToNegative[i].ConvertToInterval()
			pos := m.distancesToPositive[i].ConvertToInterval()
			m.relativeCloseness[i] = Interval{neg.Start / (neg.End + pos.End),
				neg.End / (neg.Start + pos.Start)}
		} else {
			m.relativeCloseness[i] = m.distancesToNegative[i].ConvertToNumbers() /
				(m.distancesToNegative[i].ConvertToNumbers() + m.distancesToPositive[i].ConvertToNumbers())
		}
	}
	m.closenessFind = true
}
