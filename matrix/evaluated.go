package matrix

import (
	"reflect"
)

type Evaluated interface {
	ConvertToNumbers() Numbers
	ConvertToInterval() Interval
	Weighted(weight Evaluated) Evaluated
	String() string
	DistanceTo(other Evaluated) (Numbers, error)
}

func max(x, y Evaluated, typeOfCriterion bool) Evaluated {
	if reflect.TypeOf(x) == reflect.TypeOf(Interval{}) {
		tmp := x.ConvertToInterval()
		if typeOfCriterion == Benefit {
			x = tmp.End
		} else {
			x = tmp.Start
		}
	}

	if reflect.TypeOf(y) == reflect.TypeOf(Interval{}) {
		tmp := y.ConvertToInterval()
		if typeOfCriterion == Benefit {
			y = tmp.End
		} else {
			y = tmp.Start
		}
	}

	return IdealRate(x.ConvertToNumbers(), y.ConvertToNumbers(), typeOfCriterion)
}

func min(x, y Evaluated, typeOfCriterion bool) Evaluated {
	if reflect.TypeOf(x) == reflect.TypeOf(Interval{}) {
		tmp := x.ConvertToInterval()
		if typeOfCriterion == Benefit {
			x = tmp.Start
		} else {
			x = tmp.End
		}
	}

	if reflect.TypeOf(y) == reflect.TypeOf(Interval{}) {
		tmp := y.ConvertToInterval()
		if typeOfCriterion == Benefit {
			y = tmp.Start
		} else {
			y = tmp.End
		}
	}

	return IdealRate(x.ConvertToNumbers(), y.ConvertToNumbers(), !typeOfCriterion)
}
