package smart

import (
	"fmt"
	"go.uber.org/goleak"
	"reflect"
	"testing"
	"webApp/lib/eval"
	"webApp/lib/matrix"
	v "webApp/lib/variables"
)

func SmartCalculating(smartMatrix *SmartMatrix, normValue, NormWeights v.Variants) ([]eval.Rating, error) {
	if err := matrix.TypingMatrices(5, *smartMatrix.Matrix); err != nil {
		return nil, err
	}

	if err := smartMatrix.Normalization(normValue, NormWeights, 5); err != nil {
		return nil, err
	}

	smartMatrix.CalcWeightedMatrix(3)

	smartMatrix.CalcFinalScore(5)

	return smartMatrix.GetScores(), nil
}

func TestTableDriven(t *testing.T) {
	var testCases = []struct {
		initMat    *SmartMatrix
		valueNorm  v.Variants
		weightNorm v.Variants
		resultRow  []eval.Evaluated
	}{
		{
			initMat:    ConvertToSmartMatrix(matrix.GenerateMatrix(reflect.TypeOf(eval.Number(0)), reflect.TypeOf(eval.Number(0)), 100)),
			valueNorm:  v.NormalizeWithSum,
			weightNorm: v.NormalizeWithSum,
			resultRow:  []eval.Evaluated{eval.Number(0.275), eval.Number(0.464), eval.Number(0.493), eval.Number(0.625)},
		},
		{
			initMat:    ConvertToSmartMatrix(matrix.GenerateMatrix(reflect.TypeOf(eval.Interval{}), reflect.TypeOf(eval.Interval{}), 150)),
			valueNorm:  v.NormalizeWithSum,
			weightNorm: v.NormalizeWithSum,
			resultRow: []eval.Evaluated{
				eval.Interval{0.490, 0.556},
				eval.Interval{0.221, 0.251},
				eval.Interval{0.320, 0.504}},
		},
		{
			initMat:    ConvertToSmartMatrix(matrix.GenerateMatrix(reflect.TypeOf(&eval.T1FS{}), reflect.TypeOf(eval.Interval{}), 100)),
			valueNorm:  v.NormalizeValueWithMax,
			weightNorm: v.NormalizeWithSum,
			resultRow: []eval.Evaluated{
				eval.NewT1FS([]eval.Number{0.175, 0.188, 0.275, 0.280}...),
				eval.NewT1FS([]eval.Number{0.208, 0.315, 0.315, 0.425}...),
				eval.NewT1FS([]eval.Number{0.270, 0.290, 0.713, 0.775}...),
				eval.NewT1FS([]eval.Number{0.502, 0.502, 0.537, 0.571}...),
			},
		},
		{
			initMat:    ConvertToSmartMatrix(matrix.GenerateMatrix(reflect.TypeOf(&eval.AIFS{}), reflect.TypeOf(eval.Interval{}), 100)),
			valueNorm:  v.NormalizeValueWithMax,
			weightNorm: v.NormalizeWeightsByMidPoint,
			resultRow: []eval.Evaluated{
				eval.NewAIFS(0, []eval.Number{0.231, 0.231, 0.656, 0.657}...),
				eval.NewAIFS(0, []eval.Number{0.179, 0.349, 0.349, 0.519}...),
				eval.NewAIFS(0.681, []eval.Number{0.260, 0.261, 0.261, 0.285}...),
				eval.NewAIFS(0.801, []eval.Number{0.220, 0.412, 0.442, 0.967}...),
			},
		},
		{
			initMat:    ConvertToSmartMatrix(matrix.GenerateMatrix(reflect.TypeOf(&eval.IT2FS{}), reflect.TypeOf(eval.Interval{}), 108)),
			valueNorm:  v.NormalizeWithSum,
			weightNorm: v.NormalizeWeightsByMidPoint,
			resultRow: []eval.Evaluated{
				eval.NewIT2FS([]eval.Interval{{0.094, 0.143}, {0.475, 0.510}}, []eval.Number{0.182, 0.375}),
				eval.NewIT2FS([]eval.Interval{{0.281, 0.288}, {0.475, 0.480}}, []eval.Number{0.350, 0.430}),
				eval.NewIT2FS([]eval.Interval{{0.217, 0.241}, {0.325, 0.483}}, []eval.Number{0.261, 0.319}),
			},
		},
	}

	for i, tt := range testCases {
		testname := fmt.Sprintf("%d test", i)
		t.Run(testname, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			fmt.Println(tt.initMat.String())
			fmt.Println(tt.initMat.Criteria)
			if res, err := SmartCalculating(tt.initMat, tt.valueNorm, tt.weightNorm); err != nil {
				t.Errorf(err.Error())
			} else {
				fmt.Println(tt.initMat.RankedList(v.Default))
				fmt.Println(tt.initMat.FinalScores)
				for i, el := range res {
					if !el.Equals(tt.resultRow[i]) {
						t.Errorf("got %s, want %s\n", el.String(), tt.resultRow[i].String())
					}
				}
			}
		})
	}
}
