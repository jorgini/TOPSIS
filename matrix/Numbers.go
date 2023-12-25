package matrix

import (
	"errors"
	"fmt"
	"math"
	"reflect"
)

type Numbers float64

const NumbersMax Numbers = 2147483647

const NumbersMin Numbers = -2147483647

func (n Numbers) String() string {
	return fmt.Sprintf("%.2f", float64(n))
}

func (n Numbers) ConvertToNumbers() Numbers {
	return n
}

func (n Numbers) ConvertToInterval() Interval {
	return Interval{n, n}
}

func (n Numbers) Weighted(weight Evaluated) Evaluated {
	if reflect.TypeOf(weight) == reflect.TypeOf(Numbers(0)) {
		return n * weight.ConvertToNumbers()
	} else if reflect.TypeOf(weight) == reflect.TypeOf(Interval{}) {
		// придумать
	}
	return nil
}

func (n Numbers) DistanceTo(other Evaluated) (Numbers, error) {
	if reflect.TypeOf(other) != reflect.TypeOf(n) {
		return 0, errors.New("incomparable values")
	}

	return (n - other.ConvertToNumbers()) * (n - other.ConvertToNumbers()), nil
}

func NormalizeNumbers(data []Numbers, v Variants) error {
	if v == 1 {
		sum := Numbers(0)

		for i := range data {
			sum += data[i] * data[i]
		}

		sum = Numbers(math.Sqrt(float64(sum)))

		if sum == 0.0 {
			return errors.New("can't normalize")
		}

		for i := range data {
			data[i] /= sum
		}
	} else if v == 2 {
		max := NumbersMin

		for i := range data {
			if max < data[i] {
				max = data[i]
			}
		}

		if max == 0.0 {
			return errors.New("can't normalize")
		}

		for i := range data {
			data[i] /= max
		}
	} else {
		return errors.New("invalid type of normalization for numbers")
	}

	return nil
}

func IdealRate(x, y Numbers, typeOfCriterion bool) Numbers {
	if typeOfCriterion == Benefit {
		return Numbers(math.Max(float64(x), float64(y)))
	}
	return Numbers(math.Min(float64(x), float64(y)))
}
