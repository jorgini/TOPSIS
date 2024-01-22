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
	distancesToPositive []Evaluated
	distancesToNegative []Evaluated
	distancesFind       bool
	relativeCloseness   []Evaluated
	closenessFind       bool
	highType            reflect.Type
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
	set := make([]Evaluated, m.countAlternatives)
	for i := 0; i < m.countAlternatives; i++ {
		set[i] = m.relativeCloseness[i]
	}

	for i := 0; i < m.countAlternatives; i++ {
		s += fmt.Sprint(i+1) + ": "
		indMax := 0
		if reflect.TypeOf(set[i]) == reflect.TypeOf(NumbersMin) {
			max := NumbersMin

			for j := range set {
				if set[j].ConvertToNumbers() > max {
					max = set[j].ConvertToNumbers()
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
		s += m.data[indMax].String() + "\t" + m.relativeCloseness[indMax].String() + "\n"
		set[indMax] = NumbersMin
	}
	return s
}

func NewMatrix(x, y int) *Matrix {
	m := Matrix{
		data:                make([]Alternative, x),
		countAlternatives:   x,
		countCriteria:       y,
		criteria:            make([]Criterion, y),
		distancesToNegative: make([]Evaluated, x),
		distancesToPositive: make([]Evaluated, x),
		relativeCloseness:   make([]Evaluated, x),
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

func TypingMatrices(matrices ...*Matrix) (reflect.Type, error) {
	hasInterval := false
	hasT1FSTriangle := false
	hasT1FSTrapezoid := false
	x, y := matrices[0].countAlternatives, matrices[0].countCriteria

	for k := range matrices {
		if matrices[k].countAlternatives != x || matrices[k].countCriteria != y {
			return nil, errors.New("incompatible sizes of matrices")
		}

		for i := range matrices[k].data {
			if hasT1FSTrapezoid {
				break
			}

			for _, eval := range matrices[k].data[i].grade {
				if reflect.TypeOf(eval) == reflect.TypeOf(Interval{}) {
					hasInterval = true
				}
				if reflect.TypeOf(eval) == reflect.TypeOf(&T1FS{}) {
					if eval.ConvertToT1FS(Default).form == Triangle {
						hasT1FSTriangle = true
					} else {
						hasT1FSTrapezoid = true
						break
					}
				}
			}
		}
	}

	ret := reflect.TypeOf(NumbersMin)
	if hasT1FSTrapezoid || hasT1FSTriangle {
		ret = reflect.TypeOf(&T1FS{})
	} else if hasInterval {
		ret = reflect.TypeOf(Interval{})
	}

	for k := range matrices {
		for i := range matrices[k].data {
			for j := range matrices[k].data[i].grade {
				if hasT1FSTrapezoid {
					matrices[k].data[i].grade[j] = matrices[k].data[i].grade[j].ConvertToT1FS(Trapezoid)
				} else if hasT1FSTriangle {
					matrices[k].data[i].grade[j] = matrices[k].data[i].grade[j].ConvertToT1FS(Triangle)
				} else if hasInterval {
					matrices[k].data[i].grade[j] = matrices[k].data[i].grade[j].ConvertToInterval()
				}
			}
		}
	}

	return ret, nil
}

func Aggregate(matrices []*Matrix, weights []Evaluated, mode Variants) (*Matrix, error) {
	if mode == AggregateRatings {
		result := NewMatrix(matrices[0].countAlternatives, matrices[0].countCriteria)
		if highType, err := TypingMatrices(matrices...); err != nil {
			return nil, errors.Join(err)
		} else {
			result.highType = highType
		}

		for k := range matrices {
			for i := range matrices[k].data {
				for j, eval := range matrices[k].data[i].grade {
					if result.data[i].grade[j] == nil {
						result.data[i].grade[j] = eval.Weighted(weights[k])
					} else {
						if ret, err := result.data[i].grade[j].Sum(eval.Weighted(weights[k])); err != nil {
							return nil, errors.Join(err)
						} else {
							result.data[i].grade[j] = ret
						}
					}
				}
			}
		}
		return result, nil
	} else if mode == AggregateDistances {
		x := matrices[0].countAlternatives
		result := NewMatrix(matrices[0].countAlternatives, matrices[0].countCriteria)
		for k := range matrices {
			if x != matrices[k].countAlternatives {
				return nil, errors.New("incompatible sizes of matrices")
			}

			for i, dist := range matrices[k].distancesToPositive {
				if tmp, err := result.distancesToPositive[i].Sum(dist.Weighted(weights[k]).ConvertToNumbers()); err != nil {
					return nil, errors.Join(err)
				} else {
					result.distancesToPositive[i] = tmp
				}
			}
		}
		return result, nil
	} else {
		return nil, errors.New("incomplete mode of aggregating")
	}
}
