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
	} else {
		ceil = 3
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
