package tests

import (
	"fmt"
	topsis "github.com/jorka/TOPSIS/lib"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func GenerateMatrix(valueType, weightType reflect.Type, seed int64) *topsis.Matrix {
	gen := rand.New(rand.NewSource(seed))
	var ceil int
	if valueType == reflect.TypeOf(topsis.Number(0)) {
		ceil = 1
	} else if valueType == reflect.TypeOf(topsis.Interval{}) {
		ceil = 2
	} else {
		ceil = 3
	}

	n, m := gen.Intn(5)+1, gen.Intn(5)+1
	matrix := topsis.NewMatrix(n, m)

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			typeEval := gen.Intn(ceil) + 1

			if typeEval == 1 {
				matrix.SetValue(topsis.Number(gen.Float64()*10), i, j)
			} else if typeEval == 2 {
				a, b := gen.Float64()*10, gen.Float64()*10
				matrix.SetValue(topsis.Interval{topsis.Number(math.Min(a, b)), topsis.Number(math.Max(a, b))}, i, j)
			} else {
				var vert []topsis.Number
				if gen.Float64() > 0.5 {
					vert = make([]topsis.Number, 3)
				} else {
					vert = make([]topsis.Number, 4)
				}

				for i := range vert {
					vert[i] = topsis.Number(gen.Float64() * 10)
				}

				sort.Slice(vert, func(i, j int) bool {
					return vert[i] < vert[j]
				})

				matrix.SetValue(topsis.NewT1FS(vert...), i, j)
			}
		}
	}

	for i := 0; i < m; i++ {
		if weightType == reflect.TypeOf(topsis.Interval{}) && gen.Float64() > 0.5 {
			a, b := gen.Float64()*10, gen.Float64()*10
			matrix.SetCriterion(topsis.Interval{topsis.Number(math.Min(a, b)), topsis.Number(math.Max(a, b))},
				gen.Float64() > 0.3, i)
		} else {
			matrix.SetCriterion(topsis.Number(gen.Float64()*10), gen.Float64() > 0.3, i)
		}
	}

	return matrix
}

func TopsisCalculating(matrix *topsis.TopsisMatrix, valueNorm, weightNorm, idelaAlg, fsDist,
	intDist, numDist topsis.Variants) ([]topsis.Evaluated, error) {
	if t, f, err := topsis.TypingMatrices(matrix.Matrix); err != nil {
		return nil, err
	} else {
		matrix.SetTypeAndForm(t, f)
	}

	if err := matrix.Normalization(valueNorm, weightNorm); err != nil {
		return nil, err
	}

	matrix.CalcWeightedMatrix()

	matrix.FindIdeals(idelaAlg)

	if err := matrix.FindDistanceToIdeals(fsDist, intDist, numDist); err != nil {
		return nil, err
	}

	matrix.CalcCloseness()

	return matrix.GetCoefs(), nil
}

func SmartCalculating(matrix *topsis.SmartMatrix) ([]topsis.Evaluated, error) {
	if _, _, err := topsis.TypingMatrices(matrix.Matrix); err != nil {
		return nil, err
	}

	if err := matrix.Normalization(topsis.NormalizeWithSum, topsis.NormalizeWithSum); err != nil {
		return nil, err
	}

	matrix.CalcWeightedMatrix()

	matrix.CalcFinalScore()

	return matrix.GetScores(), nil
}

func TestTableDriven(t *testing.T) {
	var tests = []struct {
		initMat                                                   *topsis.TopsisMatrix
		valueNorm, weightNorm, idelaAlg, fsDist, intDist, numDist topsis.Variants
		resultRow                                                 []topsis.Number
	}{
		{
			initMat:    topsis.ConvertToTopsisMatrix(GenerateMatrix(reflect.TypeOf(topsis.Number(0)), reflect.TypeOf(topsis.Number(0)), 100)),
			valueNorm:  topsis.NormalizeWithSum,
			weightNorm: topsis.NormalizeWithSum,
			idelaAlg:   topsis.Default,
			fsDist:     topsis.Default,
			intDist:    topsis.Default,
			numDist:    topsis.SqrtDistance,
			resultRow:  []topsis.Number{0.676, 0.547, 0.266, 0.306},
		},
		{
			initMat:    topsis.ConvertToTopsisMatrix(GenerateMatrix(reflect.TypeOf(topsis.Interval{}), reflect.TypeOf(topsis.Interval{}), 150)),
			valueNorm:  topsis.NormalizeWithSum,
			weightNorm: topsis.NormalizeWithSum,
			idelaAlg:   topsis.Default,
			fsDist:     topsis.Default,
			intDist:    topsis.Default,
			numDist:    topsis.SqrtDistance,
			resultRow:  []topsis.Number{0.61, 0.61, 0.49},
		},
		{
			initMat:    topsis.ConvertToTopsisMatrix(GenerateMatrix(reflect.TypeOf(&topsis.T1FS{}), reflect.TypeOf(topsis.Interval{}), 100)),
			valueNorm:  topsis.NormalizeWithMax,
			weightNorm: topsis.NormalizeWithSum,
			idelaAlg:   topsis.Default,
			fsDist:     topsis.Default,
			intDist:    topsis.Default,
			numDist:    topsis.CbrtDistance,
			resultRow:  []topsis.Number{0.57, 0.63, 0.35, 0.63},
		},
	}

	for i, tt := range tests {
		testname := fmt.Sprintf("%d test", i)
		t.Run(testname, func(t *testing.T) {
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
