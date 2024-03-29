package smart

import (
	"webApp/lib/eval"
	"webApp/lib/matrix"
	v "webApp/lib/variables"
)

func SmartCalculating(smartMatrix *SmartMatrix) ([]eval.Rating, error) {
	if err := matrix.TypingMatrices(*smartMatrix.Matrix); err != nil {
		return nil, err
	}

	if err := smartMatrix.Normalization(v.NormalizeWithSum, v.NormalizeWithSum); err != nil {
		return nil, err
	}

	smartMatrix.CalcWeightedMatrix()

	smartMatrix.CalcFinalScore()

	return smartMatrix.GetScores(), nil
}
