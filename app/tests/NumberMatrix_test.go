package tests

import (
	"fmt"
	"go.uber.org/goleak"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"webApp/lib/eval"
	"webApp/lib/matrix"
	"webApp/lib/smart"
	"webApp/lib/topsis"
	v "webApp/lib/variables"
)

func GenerateMatrix(valueType, weightType reflect.Type, seed int64) *matrix.Matrix {
	gen := rand.New(rand.NewSource(seed))
	var ceil int
	if valueType == reflect.TypeOf(eval.Number(0)) {
		ceil = 1
	} else if valueType == reflect.TypeOf(eval.Interval{}) {
		ceil = 2
	} else {
		ceil = 3
	}

	n, m := gen.Intn(5)+1, gen.Intn(5)+1
	newMatrix := matrix.NewMatrix(n, m)

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			typeEval := gen.Intn(ceil) + 1

			if typeEval == 1 {
				_ = newMatrix.SetValue(eval.Number(gen.Float64()*10), i, j)
			} else if typeEval == 2 {
				a, b := gen.Float64()*10, gen.Float64()*10
				_ = newMatrix.SetValue(eval.Interval{Start: eval.Number(math.Min(a, b)), End: eval.Number(math.Max(a, b))}, i, j)
			} else {
				var vert []eval.Number
				if gen.Float64() > 0.5 {
					vert = make([]eval.Number, 3)
				} else {
					vert = make([]eval.Number, 4)
				}

				for i := range vert {
					vert[i] = eval.Number(gen.Float64() * 10)
				}

				sort.Slice(vert, func(i, j int) bool {
					return vert[i] < vert[j]
				})

				_ = newMatrix.SetValue(eval.NewT1FS(vert...), i, j)
			}
		}
	}

	for i := 0; i < m; i++ {
		if weightType == reflect.TypeOf(eval.Interval{}) && gen.Float64() > 0.5 {
			a, b := gen.Float64()*10, gen.Float64()*10
			_ = newMatrix.SetCriterion(eval.Interval{Start: eval.Number(math.Min(a, b)), End: eval.Number(math.Max(a, b))},
				gen.Float64() > 0.3, i)
		} else {
			_ = newMatrix.SetCriterion(eval.Number(gen.Float64()*10), gen.Float64() > 0.3, i)
		}
	}

	return newMatrix
}

func TopsisCalculating(topsisMatrix *topsis.TopsisMatrix, valueNorm, weightNorm, idelaAlg, fsDist,
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

func SmartCalculating(smartMatrix *smart.SmartMatrix) ([]eval.Rating, error) {
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

func TestTableDriven(t *testing.T) {
	var tests = []struct {
		initMat                                                   *topsis.TopsisMatrix
		valueNorm, weightNorm, idelaAlg, fsDist, intDist, numDist v.Variants
		resultRow                                                 []eval.Number
	}{
		{
			initMat:    topsis.ConvertToTopsisMatrix(GenerateMatrix(reflect.TypeOf(eval.Number(0)), reflect.TypeOf(eval.Number(0)), 100)),
			valueNorm:  v.NormalizeWithSum,
			weightNorm: v.NormalizeWithSum,
			idelaAlg:   v.Default,
			fsDist:     v.Default,
			intDist:    v.Default,
			numDist:    v.SqrtDistance,
			resultRow:  []eval.Number{0.676, 0.547, 0.266, 0.306},
		},
		{
			initMat:    topsis.ConvertToTopsisMatrix(GenerateMatrix(reflect.TypeOf(eval.Interval{}), reflect.TypeOf(eval.Interval{}), 150)),
			valueNorm:  v.NormalizeWithSum,
			weightNorm: v.NormalizeWithSum,
			idelaAlg:   v.Default,
			fsDist:     v.Default,
			intDist:    v.Default,
			numDist:    v.SqrtDistance,
			resultRow:  []eval.Number{0.61, 0.61, 0.49},
		},
		{
			initMat:    topsis.ConvertToTopsisMatrix(GenerateMatrix(reflect.TypeOf(&eval.T1FS{}), reflect.TypeOf(eval.Interval{}), 100)),
			valueNorm:  v.NormalizeValueWithMax,
			weightNorm: v.NormalizeWithSum,
			idelaAlg:   v.Default,
			fsDist:     v.Default,
			intDist:    v.Default,
			numDist:    v.CbrtDistance,
			resultRow:  []eval.Number{0.57, 0.63, 0.35, 0.63},
		},
	}

	for i, tt := range tests {
		testname := fmt.Sprintf("%d test", i)
		t.Run(testname, func(t *testing.T) {
			defer goleak.VerifyNone(t)
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
		})
	}
}
