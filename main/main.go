package main

import (
	"fmt"
	pac "github.com/jorka/TOPSIS/lib"
	"math"
	"math/rand"
	"reflect"
	"sort"
)

func GenerateTest(valueType, weightType reflect.Type, seed int64) *pac.Matrix {
	gen := rand.New(rand.NewSource(seed))
	var ceil int
	if valueType == reflect.TypeOf(pac.Number(0)) {
		ceil = 1
	} else if valueType == reflect.TypeOf(pac.Interval{}) {
		ceil = 2
	} else {
		ceil = 3
	}

	n, m := gen.Intn(5)+1, gen.Intn(5)+1
	matrix := pac.NewMatrix(n, m)

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			typeEval := gen.Intn(ceil) + 1

			if typeEval == 1 {
				matrix.SetValue(pac.Number(gen.Float64()*10), i, j)
			} else if typeEval == 2 {
				a, b := gen.Float64()*10, gen.Float64()*10
				matrix.SetValue(pac.Interval{pac.Number(math.Min(a, b)), pac.Number(math.Max(a, b))}, i, j)
			} else {
				var vert []pac.Number
				if gen.Float64() > 0.5 {
					vert = make([]pac.Number, 3)
				} else {
					vert = make([]pac.Number, 4)
				}

				for i := range vert {
					vert[i] = pac.Number(gen.Float64() * 10)
				}

				sort.Slice(vert, func(i, j int) bool {
					return vert[i] < vert[j]
				})

				matrix.SetValue(pac.NewT1FS(vert...), i, j)
			}
		}
	}

	for i := 0; i < m; i++ {
		if weightType == reflect.TypeOf(pac.Interval{}) && gen.Float64() > 0.5 {
			a, b := gen.Float64()*10, gen.Float64()*10
			matrix.SetCriterion(pac.Interval{pac.Number(math.Min(a, b)), pac.Number(math.Max(a, b))},
				gen.Float64() > 0.3, i)
		} else {
			matrix.SetCriterion(pac.Number(gen.Float64()*10), gen.Float64() > 0.3, i)
		}
	}

	return matrix
}

func main() {
	m := GenerateTest(reflect.TypeOf(&pac.T1FS{}), reflect.TypeOf(pac.Interval{}), 100)
	fmt.Println(pac.ConvertToTopsisMatrix(m))
}
