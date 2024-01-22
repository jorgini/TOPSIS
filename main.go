package main

import (
	"fmt"
	pac "github.com/jorka/TOPSIS/matrix"
)

func main() {
	m := pac.NewMatrix(3, 3)

	m.SetValue(pac.Interval{3, 5}, 0, 0)
	m.SetValue(pac.Number(5), 0, 1)
	m.SetValue(pac.Interval{2, 4}, 0, 2)
	m.SetValue(pac.Interval{0, 5}, 1, 0)
	m.SetValue(pac.NewT1FS(3.0, 4.5, 5.3), 1, 1)
	m.SetValue(pac.Interval{4, 4}, 1, 2)
	m.SetValue(pac.Number(3.5), 2, 0)
	m.SetValue(pac.Interval{5.5, 7.6}, 2, 1)
	m.SetValue(pac.NewT1FS(6.7, 7.3, 8.0), 2, 2)

	m.SetCriterion(pac.Number(5), pac.Benefit, 0)
	m.SetCriterion(pac.Number(8), pac.Cost, 1)
	m.SetCriterion(pac.Number(2), pac.Benefit, 2)

	pac.TypingMatrices(m)

	m.Normalization(1, 1)

	fmt.Println(m)

	m.CalcWeightedMatrix()

	fmt.Println(m)

	m.FindIdeals(pac.Default)

	fmt.Println(m)

	m.FindDistanceToIdeals(pac.Default, pac.Default, pac.SqrtDistance)

	m.CalcCloseness()

	fmt.Println(m)

	fmt.Println(m.Result())

}
