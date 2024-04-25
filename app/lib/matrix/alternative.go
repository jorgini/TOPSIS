package matrix

import (
	"math"
	"webApp/lib/eval"
	v "webApp/lib/variables"
)

type Alternative struct {
	Grade           []eval.Rating `json:"grade"`
	CountOfCriteria int           `json:"cnt_of_crit"`
	_               [60]byte
}

func newAlternative(size int) Alternative {
	alt := Alternative{Grade: make([]eval.Rating, size), CountOfCriteria: size}
	for i := range alt.Grade {
		alt.Grade[i] = eval.Rating{Evaluated: nil}
	}
	return alt
}

func (a *Alternative) String() string {
	s := "[ "
	for i := 0; i < a.CountOfCriteria; i++ {
		s += a.Grade[i].String() + " "
	}
	return s + " ]"
}

func (a *Alternative) NumberMetric(to Alternative, vn v.Variants) (eval.Number, error) {
	if a.CountOfCriteria != to.CountOfCriteria {
		return 0, v.InvalidSize
	}

	result := eval.Number(0)
	for i := 0; i < a.CountOfCriteria; i++ {
		if tmp, err := a.Grade[i].DiffNumber(to.Grade[i].ConvertToNumber(), vn); err != nil {
			return 0, err
		} else {
			result += tmp
		}
	}

	if vn == v.SqrtDistance {
		result = eval.Number(math.Sqrt(float64(result)))
	} else {
		result = eval.Number(math.Cbrt(float64(result)))
	}

	return result, nil
}

func (a *Alternative) IntervalMetric(to Alternative, Criteria []Criterion, variants v.Variants) (eval.Interval, error) {
	if a.CountOfCriteria != to.CountOfCriteria {
		return eval.Interval{}, v.InvalidSize
	}

	result := eval.Interval{}
	for i, c := range Criteria {
		if tmp, err := a.Grade[i].DiffInterval(to.Grade[i].ConvertToInterval(), c.TypeOfCriteria, variants); err != nil {
			return eval.Interval{}, err
		} else {
			result = result.Sum(tmp).ConvertToInterval()
		}
	}

	if variants == v.SqrtDistance {
		result = eval.Interval{Start: eval.Number(math.Sqrt(float64(result.ConvertToInterval().Start))),
			End: eval.Number(math.Sqrt(float64(result.ConvertToInterval().End)))}
	} else {
		result = eval.Interval{Start: eval.Number(math.Cbrt(float64(result.ConvertToInterval().Start))),
			End: eval.Number(math.Cbrt(float64(result.ConvertToInterval().End)))}
	}

	return result.ConvertToInterval(), nil
}

func (a *Alternative) FSMetric(to Alternative, vn v.Variants) (eval.Number, error) {
	if a.CountOfCriteria != to.CountOfCriteria {
		return 0, v.InvalidSize
	}

	result := eval.Number(0)
	for i := 0; i < a.CountOfCriteria; i++ {
		if tmp, err := a.Grade[i].DiffNumber(to.Grade[i], vn); err != nil {
			return 0, err
		} else {
			result += tmp
		}
	}

	if vn == v.SqrtDistance {
		result = eval.Number(math.Sqrt(float64(result)))
	} else {
		result = eval.Number(math.Cbrt(float64(result)))
	}

	return result, nil
}

func (a *Alternative) Sum() eval.Rating {
	var sum eval.Evaluated

	if a.Grade[0].GetType() == eval.NumbersMin.GetType() {
		sum = eval.Number(0.0)
	} else if a.Grade[0].GetType() == (eval.Interval{}).GetType() {
		sum = eval.Interval{}
	} else if a.Grade[0].GetType() == (&eval.T1FS{}).GetType() {
		if a.Grade[0].ConvertToT1FS(v.Default).Form == v.Triangle {
			sum = eval.NewT1FS(0.0, 0.0, 0.0)
		} else {
			sum = eval.NewT1FS(0.0, 0.0, 0.0, 0.0)
		}
	} else if a.Grade[0].GetType() == (&eval.AIFS{}).GetType() {
		if a.Grade[0].ConvertToAIFS(v.Default).Form == v.Triangle {
			sum = eval.NewAIFS(0.0, 0.0, 0.0, 0.0)
		} else {
			sum = eval.NewAIFS(0.0, 0.0, 0.0, 0.0, 0.0)
		}
	} else if a.Grade[0].GetType() == (&eval.IT2FS{}).GetType() {
		if a.Grade[0].ConvertToIT2FS(v.Default).Form == v.Triangle {
			sum = eval.NewIT2FS([]eval.Interval{{0.0, 0.0}, {0.0, 0.0}}, []eval.Number{0.0})
		} else {
			sum = eval.NewIT2FS([]eval.Interval{{0.0, 0.0}, {0.0, 0.0}}, []eval.Number{0.0, 0.0})
		}
	}

	for _, el := range a.Grade {
		sum = sum.Sum(el)
	}
	return eval.Rating{Evaluated: sum}
}
