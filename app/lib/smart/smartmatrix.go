package smart

import (
	"webApp/lib/eval"
	"webApp/lib/matrix"
	v "webApp/lib/variables"
)

type SmartMatrix struct {
	*matrix.Matrix
	FinalScores     []eval.Rating `json:"final_scores"`
	FinalScoresFind bool          `json:"is_final_find"`
}

func (sm *SmartMatrix) GetScores() []eval.Rating {
	return sm.FinalScores
}

func (sm *SmartMatrix) String() string {
	s := "Decision matrix:\n"
	for i := range sm.Data {
		s += sm.Data[i].String() + "\n"
	}

	if sm.CriteriaSet {
		s += "\nWeights of Criteria:\n"
		for i := 0; i < sm.CountCriteria; i++ {
			s += sm.Criteria[i].String() + " "
		}
	}

	if sm.FinalScoresFind {
		s += "\nFinal Scores:\n"
		for i := 0; i < sm.CountAlternatives; i++ {
			s += sm.FinalScores[i].String() + " "
		}
		s += "\n"
	}
	return s
}

func NewSmartMatrix(x, y int) *SmartMatrix {
	sm := SmartMatrix{
		Matrix:      matrix.NewMatrix(x, y),
		FinalScores: make([]eval.Rating, x),
	}

	for i := range sm.Data {
		sm.Data[i] = matrix.Alternative{Grade: make([]eval.Rating, y), CountOfCriteria: y}
	}
	return &sm
}

func ConvertToSmartMatrix(m *matrix.Matrix) *SmartMatrix {
	return &SmartMatrix{
		Matrix:          matrix.CopyMatrix(m),
		FinalScores:     make([]eval.Rating, m.CountAlternatives),
		FinalScoresFind: false,
	}
}

func AggregateScores(matrices []SmartMatrix, weights []eval.Evaluated) (*SmartMatrix, error) {
	x := matrices[0].CountAlternatives
	result := NewSmartMatrix(matrices[0].CountAlternatives, matrices[0].CountCriteria)
	for k := range matrices {
		if x != matrices[k].CountAlternatives {
			return nil, v.InvalidSize
		}
		for i, score := range matrices[k].FinalScores {
			if k == 0 {
				result.FinalScores[i] = score.Weighted(weights[k])
				continue
			}
			result.FinalScores[i] = result.FinalScores[i].Sum(score.Weighted(weights[k]))
		}
	}
	return result, nil
}
