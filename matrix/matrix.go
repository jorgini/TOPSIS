package matrix

import (
	"errors"
	"fmt"
	"reflect"
)

type Matrix struct {
	data                []Alternative
	countAlternatives   int
	countCriteria       int
	criteria            []Criterion
	criteriaSet         bool
	positiveIdeal       *Alternative
	negativeIdeal       *Alternative
	idealsFind          bool
	distancesToPositive []Numbers
	distancesToNegative []Numbers
	distancesFind       bool
	relativeCloseness   []Numbers
	closenessFind       bool
}

func (m *Matrix) String() string {
	s := "Decision matrix:\n"
	for i := range m.data {
		s += m.data[i].String() + "\n"
	}

	if m.criteriaSet {
		s += "\nWeights of criteria:\n"
		for i := 0; i < m.countCriteria; i++ {
			s += m.criteria[i].weight.String() + " "
		}
	}

	if m.idealsFind {
		s += "\nPositive-ideal alternative:\n" + m.positiveIdeal.String() +
			"\nNegative-ideal alternative:\n" + m.negativeIdeal.String()
	}

	if m.distancesFind {
		s += "\nDistances to positive:\n"
		for i := 0; i < m.countAlternatives; i++ {
			s += m.distancesToPositive[i].String() + " "
		}
		s += "\nDistance to negative:\n"
		for i := 0; i < m.countAlternatives; i++ {
			s += m.distancesToNegative[i].String() + " "
		}
	}

	if m.closenessFind {
		s += "\nRelative closeness:\n"
		for i := 0; i < m.countAlternatives; i++ {
			s += m.relativeCloseness[i].String() + " "
		}
		s += "\n"
	}
	return s
}

func (m *Matrix) Result() string {
	s := "Result:\n"
	set := make([]int, m.countAlternatives)
	for i := 0; i < m.countAlternatives; i++ {
		set[i] = i
	}

	for i := 0; i < m.countAlternatives; i++ {
		s += fmt.Sprint(i+1) + ": "
		max := 0.0
		indMax := 0
		for j := range set {
			if m.relativeCloseness[set[j]] > Numbers(max) {
				max = float64(m.relativeCloseness[set[j]])
				indMax = set[j]
			}
		}
		s += m.data[indMax].String() + "\t" + m.relativeCloseness[indMax].String() + "\n"
		set = append(set[:indMax], set[indMax+1:]...)
	}
	return s
}

func NewMatrix(x, y int) *Matrix {
	m := Matrix{
		data:                make([]Alternative, x),
		countAlternatives:   x,
		countCriteria:       y,
		criteria:            make([]Criterion, y),
		distancesToNegative: make([]Numbers, x),
		distancesToPositive: make([]Numbers, x),
		relativeCloseness:   make([]Numbers, x),
	}

	for i := range m.data {
		m.data[i] = Alternative{make([]Evaluated, y), y}
	}
	return &m
}

func (m *Matrix) SetValue(value Evaluated, i, j int) error {
	if i < m.countAlternatives && j < m.countCriteria {
		m.data[i].grade[j] = value
	} else {
		return errors.New("invalid arguments")
	}
	return nil
}

func (m *Matrix) SetCriterion(weight Evaluated, typeOF bool, i int) error {
	if i < m.countCriteria {
		m.criteria[i].set(weight, typeOF)
	} else {
		return errors.New("invalid arguments")
	}
	m.criteriaSet = true
	return nil
}

func TypingMatrices(matrices ...*Matrix) {
	hasInterval := false

	for k := range matrices {
		for i := range matrices[k].data {
			for j := range matrices[k].data[i].grade {
				if reflect.TypeOf(matrices[k].data[i].grade[j]) == reflect.TypeOf(Interval{}) {
					hasInterval = true
				}
			}
		}
	}

	if hasInterval {
		if matrices != nil {
			for k := range matrices {
				for i := range matrices[k].data {
					for j := range matrices[k].data[i].grade {
						matrices[k].data[i].grade[j] = matrices[k].data[i].grade[j].ConvertToInterval()
					}
				}
			}
		}
	}
}

func Aggregate(matrices []*Matrix, weights []Evaluated) {
	TypingMatrices(matrices...)

	for k := range matrices {
		for i := range matrices[k].data {
			for j := range matrices[k].data[i].grade {
				matrices[k].data[i].grade[j] = matrices[k].data[i].grade[j].Weighted(weights[k])
			}
		}
	}
}
