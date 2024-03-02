package lib

import (
	"fmt"
	"reflect"
)

type SmartMatrix struct {
	*Matrix
	finalScores     []Evaluated
	finalScoresFind bool
}

func (sm *SmartMatrix) GetScores() []Evaluated {
	return sm.finalScores
}

func (sm *SmartMatrix) String() string {
	s := "Decision matrix:\n"
	for i := range sm.data {
		s += sm.data[i].String() + "\n"
	}

	if sm.criteriaSet {
		s += "\nWeights of criteria:\n"
		for i := 0; i < sm.countCriteria; i++ {
			s += sm.criteria[i].String() + " "
		}
	}

	if sm.finalScoresFind {
		s += "\nFinal Scores:\n"
		for i := 0; i < sm.countAlternatives; i++ {
			s += sm.finalScores[i].String() + " "
		}
		s += "\n"
	}
	return s
}

func (sm *SmartMatrix) Result() string {
	s := "Result:\n"
	set := make([]Evaluated, sm.countAlternatives)
	for i := 0; i < sm.countAlternatives; i++ {
		set[i] = sm.finalScores[i]
	}

	for i := 0; i < sm.countAlternatives; i++ {
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
		s += sm.data[indMax].String() + "\t" + sm.finalScores[indMax].String() + "\n"
		set[indMax] = NumbersMin
	}
	return s
}

func NewSmartMatrix(x, y int) *SmartMatrix {
	sm := SmartMatrix{
		Matrix:      NewMatrix(x, y),
		finalScores: make([]Evaluated, x),
	}

	for i := range sm.data {
		sm.data[i] = Alternative{make([]Evaluated, y), y}
	}
	return &sm
}

func AggregateScores(matrices []*SmartMatrix, weights []Evaluated) (*SmartMatrix, error) {
	x := matrices[0].countAlternatives
	result := NewSmartMatrix(matrices[0].countAlternatives, matrices[0].countCriteria)
	for k := range matrices {
		if x != matrices[k].countAlternatives {
			return nil, InvalidSize
		}
		for i, dist := range matrices[k].finalScores {
			result.finalScores[i] = result.finalScores[i].Sum(dist.Weighted(weights[k]).ConvertToNumber())
		}
	}
	return result, nil
}
