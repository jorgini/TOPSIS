package main

import (
	"fmt"
	pac "github.com/jorka/TOPSIS/matrix"
)

func main() {
	m := pac.NewMatrix(3, 3)

	m.SetValue(pac.Interval{3, 5}, 0, 0)
	m.SetValue(pac.Interval{7, 9}, 0, 1)
	m.SetValue(pac.Interval{2, 4}, 0, 2)
	m.SetValue(pac.Interval{0, 5}, 1, 0)
	m.SetValue(pac.Interval{5, 6}, 1, 1)
	m.SetValue(pac.Interval{4, 4}, 1, 2)
	m.SetValue(pac.Interval{3.4, 3.8}, 2, 0)
	m.SetValue(pac.Interval{5.5, 7.6}, 2, 1)
	m.SetValue(pac.Interval{1, 4}, 2, 2)

	m.SetCriterion(pac.Numbers(5), pac.Benefit, 0)
	m.SetCriterion(pac.Numbers(8), pac.Cost, 1)
	m.SetCriterion(pac.Numbers(2), pac.Benefit, 2)

	m.Normalization(1, 1)

	fmt.Println(m)

	m.CalcWeightedMatrix()

	fmt.Println(m)

	m.FindIdeals()

	fmt.Println(m)

	m.FindDistanceToIdeals()

	m.CalcCloseness()

	fmt.Println(m)

	fmt.Println(m.Result())
}
