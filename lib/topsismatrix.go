package lib

import (
	"fmt"
	"reflect"
)

type TopsisMatrix struct {
	*Matrix
	positiveIdeal       *Alternative
	negativeIdeal       *Alternative
	idealsFind          bool
	distancesToPositive []Evaluated
	distancesToNegative []Evaluated
	distancesFind       bool
	relativeCloseness   []Evaluated
	closenessFind       bool
}

func (tm *TopsisMatrix) GetCoefs() []Evaluated {
	return tm.relativeCloseness
}

func (tm *TopsisMatrix) String() string {
	s := "Decision matrix:\n"
	for i := range tm.data {
		s += tm.data[i].String() + "\n"
	}

	if tm.criteriaSet {
		s += "\nWeights of criteria:\n"
		for i := 0; i < tm.countCriteria; i++ {
			s += tm.criteria[i].String() + " "
		}
	}

	if tm.idealsFind {
		s += "\nPositive-ideal alternative:\n" + tm.positiveIdeal.String() +
			"\nNegative-ideal alternative:\n" + tm.negativeIdeal.String()
	}

	if tm.distancesFind {
		s += "\nDistances to positive:\n"
		for i := 0; i < tm.countAlternatives; i++ {
			s += tm.distancesToPositive[i].String() + " "
		}
		s += "\nDistance to negative:\n"
		for i := 0; i < tm.countAlternatives; i++ {
			s += tm.distancesToNegative[i].String() + " "
		}
	}

	if tm.closenessFind {
		s += "\nRelative closeness:\n"
		for i := 0; i < tm.countAlternatives; i++ {
			s += tm.relativeCloseness[i].String() + " "
		}
		s += "\n"
	}
	return s
}

func (tm *TopsisMatrix) Result() string {
	s := "Result:\n"
	set := make([]Evaluated, tm.countAlternatives)
	for i := 0; i < tm.countAlternatives; i++ {
		set[i] = tm.relativeCloseness[i]
	}

	for i := 0; i < tm.countAlternatives; i++ {
		s += fmt.Sprint(i+1) + ": "
		indMax := 0
		if reflect.TypeOf(set[i]) == reflect.TypeOf(NumbersMin) {
			max := NumbersMin

			for j := range set {
				if set[j].ConvertToNumber() > max {
					max = set[j].ConvertToNumber()
					indMax = j
				}
			}
		} else {
			max := Interval{NumbersMin, NumbersMin}

			for j := range set {
				if set[j].ConvertToInterval().SenguptaGeq(max) {
					max = set[j].ConvertToInterval()
					indMax = j
				}
			}
		}
		s += tm.data[indMax].String() + "\t" + tm.relativeCloseness[indMax].String() + "\n"
		set[indMax] = NumbersMin
	}
	return s
}

func NewTopsisMatrix(x, y int) *TopsisMatrix {
	return &TopsisMatrix{
		Matrix:              NewMatrix(x, y),
		distancesToNegative: make([]Evaluated, x),
		distancesToPositive: make([]Evaluated, x),
		relativeCloseness:   make([]Evaluated, x),
	}
}

func ConvertToTopsisMatrix(m *Matrix) *TopsisMatrix {
	return &TopsisMatrix{
		Matrix:              m,
		distancesToNegative: make([]Evaluated, m.countAlternatives),
		distancesToPositive: make([]Evaluated, m.countAlternatives),
		relativeCloseness:   make([]Evaluated, m.countAlternatives),
	}
}

func AggregateDistances(matrices []*TopsisMatrix, weights []Evaluated) (*TopsisMatrix, error) {
	x := matrices[0].countAlternatives
	result := NewTopsisMatrix(matrices[0].countAlternatives, matrices[0].countCriteria)
	for k := range matrices {
		if x != matrices[k].countAlternatives {
			return nil, InvalidSize
		}
		for i, dist := range matrices[k].distancesToPositive {
			result.distancesToPositive[i] = result.distancesToPositive[i].Sum(dist.Weighted(weights[k]))
		}
	}
	return result, nil
}
