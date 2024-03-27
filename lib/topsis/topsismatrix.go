package topsis

import (
	"fmt"
	"webApp/lib/eval"
	"webApp/lib/matrix"
	v "webApp/lib/variables"
)

type TopsisMatrix struct {
	*matrix.Matrix
	PositiveIdeal       *matrix.Alternative `json:"pos_ideal"`
	NegativeIdeal       *matrix.Alternative `json:"neg_ideal"`
	IdealsFind          bool                `json:"is_ideal_find"`
	DistancesToPositive []eval.Rating       `json:"dist_to_pos"`
	DistancesToNegative []eval.Rating       `json:"dist_to_neg"`
	DistancesFind       bool                `json:"is_dist_find"`
	RelativeCloseness   []eval.Rating       `json:"relate_close"`
	ClosenessFind       bool                `json:"is_close_find"`
}

func (tm *TopsisMatrix) GetCoefs() []eval.Rating {
	return tm.RelativeCloseness
}

func (tm *TopsisMatrix) String() string {
	s := "Decision matrix:\n"
	for i := range tm.Data {
		s += tm.Data[i].String() + "\n"
	}

	if tm.CriteriaSet {
		s += "\nWeights of Criteria:\n"
		for i := 0; i < tm.CountCriteria; i++ {
			s += tm.Criteria[i].String() + " "
		}
	}

	if tm.IdealsFind {
		s += "\nPositive-ideal alternative:\n" + tm.PositiveIdeal.String() +
			"\nNegative-ideal alternative:\n" + tm.NegativeIdeal.String()
	}

	if tm.DistancesFind {
		s += "\nDistances to positive:\n"
		for i := 0; i < tm.CountAlternatives; i++ {
			s += tm.DistancesToPositive[i].String() + " "
		}
		s += "\nDistance to negative:\n"
		for i := 0; i < tm.CountAlternatives; i++ {
			s += tm.DistancesToNegative[i].String() + " "
		}
	}

	if tm.ClosenessFind {
		s += "\nRelative closeness:\n"
		for i := 0; i < tm.CountAlternatives; i++ {
			s += tm.RelativeCloseness[i].String() + " "
		}
		s += "\n"
	}
	return s
}

func (tm *TopsisMatrix) Result() string {
	s := "Result:\n"
	set := make([]eval.Rating, tm.CountAlternatives)
	for i := 0; i < tm.CountAlternatives; i++ {
		set[i] = tm.RelativeCloseness[i]
	}

	for i := 0; i < tm.CountAlternatives; i++ {
		s += fmt.Sprint(i+1) + ": "
		indMax := 0
		if set[i].GetType() == eval.NumbersMin.GetType() {
			max := eval.NumbersMin

			for j := range set {
				if set[j].ConvertToNumber() > max {
					max = set[j].ConvertToNumber()
					indMax = j
				}
			}
		} else {
			max := eval.Interval{Start: eval.NumbersMin, End: eval.NumbersMin}

			for j := range set {
				if set[j].ConvertToInterval().SenguptaGeq(max) {
					max = set[j].ConvertToInterval()
					indMax = j
				}
			}
		}
		s += tm.Data[indMax].String() + "\t" + tm.RelativeCloseness[indMax].String() + "\n"
		set[indMax] = eval.Rating{Evaluated: eval.NumbersMin}
	}
	return s
}

func NewTopsisMatrix(x, y int) *TopsisMatrix {
	return &TopsisMatrix{
		Matrix:              matrix.NewMatrix(x, y),
		DistancesToNegative: make([]eval.Rating, x),
		DistancesToPositive: make([]eval.Rating, x),
		RelativeCloseness:   make([]eval.Rating, x),
	}
}

func ConvertToTopsisMatrix(m *matrix.Matrix) *TopsisMatrix {
	return &TopsisMatrix{
		Matrix:              matrix.CopyMatrix(m),
		DistancesToNegative: make([]eval.Rating, m.CountAlternatives),
		DistancesToPositive: make([]eval.Rating, m.CountAlternatives),
		RelativeCloseness:   make([]eval.Rating, m.CountAlternatives),
	}
}

func AggregateDistances(matrices []TopsisMatrix, weights []eval.Evaluated) (*TopsisMatrix, error) {
	x := matrices[0].CountAlternatives
	result := NewTopsisMatrix(matrices[0].CountAlternatives, matrices[0].CountCriteria)

	for k := range matrices {
		if x != matrices[k].CountAlternatives {
			return nil, v.InvalidSize
		}
		for i, dist := range matrices[k].DistancesToPositive {
			if k == 0 {
				result.DistancesToPositive[i] = dist.Weighted(weights[k])
				continue
			}
			result.DistancesToPositive[i] = result.DistancesToPositive[i].Sum(dist.Weighted(weights[k]))
		}
		for i, dist := range matrices[k].DistancesToNegative {
			if k == 0 {
				result.DistancesToNegative[i] = dist.Weighted(weights[k])
				continue
			}
			result.DistancesToNegative[i] = result.DistancesToNegative[i].Sum(dist.Weighted(weights[k]))
		}
	}
	return result, nil
}
