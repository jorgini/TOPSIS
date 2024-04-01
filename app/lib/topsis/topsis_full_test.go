package topsis

import (
	"fmt"
	"go.uber.org/goleak"
	"math"
	"reflect"
	"testing"
	"webApp/lib/eval"
	"webApp/lib/matrix"
	v "webApp/lib/variables"
)

func TopsisCalculating(topsisMatrix *TopsisMatrix, valueNorm, weightNorm, idelaAlg, fsDist,
	intDist, numDist v.Variants) ([]eval.Rating, error) {
	if err := matrix.TypingMatrices(*topsisMatrix.Matrix); err != nil {
		return nil, err
	}

	if err := topsisMatrix.Normalization(valueNorm, weightNorm); err != nil {
		return nil, err
	}

	topsisMatrix.CalcWeightedMatrix()

	err := topsisMatrix.FindIdeals(idelaAlg)
	if err != nil {
		return nil, err
	}

	if err := topsisMatrix.FindDistanceToIdeals(fsDist, intDist, numDist); err != nil {
		return nil, err
	}

	topsisMatrix.CalcCloseness()

	return topsisMatrix.GetCoefs(), nil
}

func TestTableDriven(t *testing.T) {
	var tests = []struct {
		initMat                                                   *TopsisMatrix
		valueNorm, weightNorm, idelaAlg, fsDist, intDist, numDist v.Variants
		resultRow                                                 []eval.Number
	}{
		{
			initMat:    ConvertToTopsisMatrix(matrix.GenerateMatrix(reflect.TypeOf(eval.Number(0)), reflect.TypeOf(eval.Number(0)), 100)),
			valueNorm:  v.NormalizeWithSum,
			weightNorm: v.NormalizeWithSum,
			idelaAlg:   v.Default,
			fsDist:     v.Default,
			intDist:    v.Default,
			numDist:    v.SqrtDistance,
			resultRow:  []eval.Number{0.676, 0.547, 0.266, 0.306},
		},
		{
			initMat:    ConvertToTopsisMatrix(matrix.GenerateMatrix(reflect.TypeOf(eval.Interval{}), reflect.TypeOf(eval.Interval{}), 150)),
			valueNorm:  v.NormalizeWithSum,
			weightNorm: v.NormalizeWithSum,
			idelaAlg:   v.Default,
			fsDist:     v.Default,
			intDist:    v.Default,
			numDist:    v.SqrtDistance,
			resultRow:  []eval.Number{0.61, 0.61, 0.49},
		},
		{
			initMat:    ConvertToTopsisMatrix(matrix.GenerateMatrix(reflect.TypeOf(&eval.T1FS{}), reflect.TypeOf(eval.Interval{}), 100)),
			valueNorm:  v.NormalizeValueWithMax,
			weightNorm: v.NormalizeWithSum,
			idelaAlg:   v.Default,
			fsDist:     v.Default,
			intDist:    v.Default,
			numDist:    v.CbrtDistance,
			resultRow:  []eval.Number{0.57, 0.63, 0.35, 0.63},
		},
		{
			initMat:    ConvertToTopsisMatrix(matrix.GenerateMatrix(reflect.TypeOf(&eval.AIFS{}), reflect.TypeOf(eval.Interval{}), 100)),
			valueNorm:  v.NormalizeValueWithMax,
			weightNorm: v.NormalizeWeightsByMidPoint,
			idelaAlg:   v.Default,
			fsDist:     v.AlphaSlices,
			intDist:    v.Default,
			numDist:    v.CbrtDistance,
			resultRow:  []eval.Number{0.604, 0.472, 0.445, 0.314},
		},
		{
			initMat:    ConvertToTopsisMatrix(matrix.GenerateMatrix(reflect.TypeOf(&eval.IT2FS{}), reflect.TypeOf(eval.Interval{}), 108)),
			valueNorm:  v.NormalizeValueWithMax,
			weightNorm: v.NormalizeWeightsByMidPoint,
			idelaAlg:   v.Default,
			fsDist:     v.AlphaSlices,
			intDist:    v.Default,
			numDist:    v.CbrtDistance,
			resultRow:  []eval.Number{0.301, 0.334, 0.292},
		},
	}
	fmt.Println(tests[4].initMat)
	for i, tt := range tests {
		testname := fmt.Sprintf("%d test", i)
		t.Run(testname, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			fmt.Println(tt.initMat.String())
			if res, err := TopsisCalculating(tt.initMat, tt.valueNorm, tt.weightNorm, tt.idelaAlg, tt.fsDist,
				tt.intDist, tt.numDist); err != nil {
				t.Errorf(err.Error())
			} else {
				for i, el := range res {
					if math.Abs(float64(el.ConvertToNumber()-tt.resultRow[i])) > 0.01 {
						t.Errorf("got %f, want %f\n", el, tt.resultRow[i])
					}
				}
			}
			tt.initMat.Result()
		})
	}
}
