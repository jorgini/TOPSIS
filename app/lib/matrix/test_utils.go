package matrix

import (
	"math"
	"math/rand"
	"reflect"
	"sort"
	"webApp/lib/eval"
)

func GenerateMatrix(valueType, weightType reflect.Type, seed int64) *Matrix {
	gen := rand.New(rand.NewSource(seed))
	var ceil int
	if valueType == reflect.TypeOf(eval.Number(0)) {
		ceil = 1
	} else if valueType == reflect.TypeOf(eval.Interval{}) {
		ceil = 2
	} else if valueType == reflect.TypeOf(&eval.T1FS{}) {
		ceil = 3
	} else if valueType == reflect.TypeOf(&eval.AIFS{}) {
		ceil = 4
	} else if valueType == reflect.TypeOf(&eval.IT2FS{}) {
		ceil = 5
	}

	n, m := gen.Intn(5)+1, gen.Intn(5)+1
	newMatrix := NewMatrix(n, m)

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			typeEval := gen.Intn(ceil) + 1

			if typeEval == 1 {
				_ = newMatrix.SetValue(eval.Number(gen.Float64()*10), i, j)
			} else if typeEval == 2 {
				a, b := gen.Float64()*10, gen.Float64()*10
				_ = newMatrix.SetValue(eval.Interval{Start: eval.Number(math.Min(a, b)), End: eval.Number(math.Max(a, b))}, i, j)
			} else if typeEval == 3 {
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
			} else if typeEval == 4 {
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

				_ = newMatrix.SetValue(eval.NewAIFS(eval.Number(gen.Float64()), vert...), i, j)
			} else if typeEval == 5 {
				var vert []eval.Number
				if gen.Float64() > 0.5 {
					vert = make([]eval.Number, 5)
				} else {
					vert = make([]eval.Number, 6)
				}

				for i := range vert {
					vert[i] = eval.Number(gen.Float64() * 10)
				}

				sort.Slice(vert, func(i, j int) bool {
					return vert[i] < vert[j]
				})

				bottom := []eval.Interval{{vert[0], vert[1]}, {vert[len(vert)-2], vert[len(vert)-1]}}
				upward := make([]eval.Number, 0, 2)
				for i := 2; i < len(vert)-2; i++ {
					upward = append(upward, vert[i])
				}
				_ = newMatrix.SetValue(eval.NewIT2FS(bottom, upward), i, j)
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

//func AssertMatrix(first *Matrix, second *Matrix) bool {
//	for i := range first.Data {
//		for j := range first.Data[i].Grade {
//			if
//		}
//	}
//}
