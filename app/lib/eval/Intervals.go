package eval

import (
	"fmt"
	"math"
	"reflect"
	v "webApp/lib/variables"
)

type Interval struct {
	Start Number `json:"start"`
	End   Number `json:"end"`
}

func (i Interval) CopyEval() Rating {
	return Rating{Interval{Start: i.Start, End: i.End}}
}

func (i Interval) String() string {
	return fmt.Sprintf("[%.3f, %.3f]", i.Start, i.End)
}

func (i Interval) GetType() string {
	return reflect.TypeOf(i).String()
}

func (i Interval) ConvertToNumber() Number {
	return (i.Start + i.End) / 2
}

func (i Interval) ConvertToInterval() Interval {
	return i
}

func (i Interval) ConvertToT1FS(f v.Variants) *T1FS {
	if f == v.Default || f == v.Triangle {
		return NewT1FS(i.Start, i.ConvertToNumber(), i.End)
	} else {
		return NewT1FS(i.Start, i.Start, i.End, i.End)
	}
}

func (i Interval) ConvertToAIFS(f v.Variants) *AIFS {
	if f == v.Default || f == v.Triangle {
		return NewAIFS(0.0, i.Start, i.ConvertToNumber(), i.End)
	} else {
		return NewAIFS(0.0, i.Start, i.Start, i.End, i.End)
	}
}

func (i Interval) ConvertToIT2FS(f v.Variants) *IT2FS {
	if f == v.Default || f == v.Triangle {
		return NewIT2FS([]Interval{{i.Start, i.Start}, {i.End, i.End}}, []Number{i.ConvertToNumber()})
	} else {
		return NewIT2FS([]Interval{{i.Start, i.Start}, {i.End, i.End}}, []Number{i.Start, i.End})
	}
}

func (i Interval) GetForm() v.Variants {
	return v.None
}

func (i Interval) Weighted(Weight Evaluated) Rating {
	if Weight.GetType() == NumbersMin.GetType() || Weight.GetType() == i.GetType() {
		return Rating{Interval{i.Start * Weight.ConvertToNumber(), i.End * Weight.ConvertToNumber()}}
	}
	return Rating{nil}
}

func (i Interval) DiffNumber(other Evaluated, variants v.Variants) (Number, error) {
	if other.GetType() != NumbersMax.GetType() {
		return 0, v.IncompatibleTypes
	}

	if variants == v.SqrtDistance {
		return (i.Start-other.ConvertToNumber())*(i.Start-other.ConvertToNumber()) +
			(i.End-other.ConvertToNumber())*(i.End-other.ConvertToNumber()), nil
	} else if variants == v.CbrtDistance {
		return Number(math.Abs(math.Pow(float64(i.Start-other.ConvertToNumber()), 3))) +
			Number(math.Abs(math.Pow(float64(i.End-other.ConvertToNumber()), 3))), nil
	} else {
		return NumbersMin, v.InvalidCaseOfOperation
	}
}

func (i Interval) DiffInterval(other Interval, typeOfCriterion bool, variants v.Variants) (Interval, error) {
	ret := Interval{}
	if typeOfCriterion == v.Benefit {
		if variants == v.SqrtDistance {
			ret.Start = (i.End - other.Start) * (i.End - other.Start)
			ret.End = (i.Start - other.End) * (i.Start - other.End)
		} else if variants == v.CbrtDistance {
			ret.Start = Number(math.Abs(math.Pow(float64(i.End-other.Start), 3)))
			ret.End = Number(math.Abs(math.Pow(float64(i.Start-other.End), 3)))
		} else {
			return Interval{}, v.InvalidCaseOfOperation
		}
	} else if typeOfCriterion == v.Cost {
		if variants == v.SqrtDistance {
			ret.Start = (i.Start - other.Start) * (i.Start - other.Start)
			ret.End = (i.End - other.End) * (i.End - other.End)
		} else if variants == v.CbrtDistance {
			ret.Start = Number(math.Abs(math.Pow(float64(i.Start-other.Start), 3)))
			ret.End = Number(math.Abs(math.Pow(float64(i.End-other.End), 3)))
		} else {
			return Interval{}, v.InvalidCaseOfOperation
		}
	} else {
		return Interval{}, v.InvalidCaseOfOperation
	}
	return ret, nil
}

func (i Interval) Sum(other Evaluated) Rating {
	return Rating{Interval{i.Start + other.ConvertToInterval().Start,
		i.End + other.ConvertToInterval().End}}
}

func (i Interval) SenguptaGeq(other Interval) bool {
	if math.Abs(float64(i.ConvertToNumber()-other.ConvertToNumber())) < 1e-3 {
		return i.End-i.Start < other.End-other.Start
	} else {
		return i.ConvertToNumber() > other.ConvertToNumber()
	}
}

func (i Interval) Equals(other Evaluated) bool {
	if other.GetType() != i.GetType() {
		return false
	}

	if i.Start.Equals(other.ConvertToInterval().Start) == false || i.End.Equals(other.ConvertToInterval().End) == false {
		return false
	}
	return true
}
