package matrix

import (
	"errors"
	"fmt"
	"math"
	"reflect"
)

type Interval struct {
	Start, End Numbers
}

func (i Interval) String() string {
	return fmt.Sprintf("[%.2f, %.2f]", i.Start, i.End)
}

func (i Interval) ConvertToNumbers() Numbers {
	return (i.Start + i.End) / 2
}

func (i Interval) ConvertToInterval() Interval {
	return i
}

func (i Interval) Weighted(weight Evaluated) Evaluated {
	if reflect.TypeOf(weight) == reflect.TypeOf(Numbers(0)) {
		return Interval{i.Start * weight.ConvertToNumbers(),
			i.End * weight.ConvertToNumbers()}
	} else if reflect.TypeOf(weight) == reflect.TypeOf(i) {
		// придумать
	}
	return nil
}

func (i Interval) DistanceTo(other Evaluated) (Numbers, error) {
	if reflect.TypeOf(other) != reflect.TypeOf(NumbersMax) {
		return 0, errors.New("incomparable values")
	}

	return (i.Start-other.ConvertToNumbers())*(i.Start-other.ConvertToNumbers()) +
		(i.End-other.ConvertToNumbers())*(i.End-other.ConvertToNumbers()), nil
}

func NormalizeIntervals(data []Interval, v Variants) error {
	if v == 1 {
		sum := Numbers(0)

		for i := range data {
			sum += data[i].Start*data[i].Start + data[i].End*data[i].End
		}

		sum = Numbers(math.Sqrt(float64(sum)))

		if sum == 0.0 {
			return errors.New("can't normalize")
		}

		for i := range data {
			data[i].Start /= sum
			data[i].End /= sum
		}
	} else {
		return errors.New("invalid type of normalization for intervals")
	}
	return nil
}
