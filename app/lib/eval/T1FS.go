package eval

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	v "webApp/lib/variables"
)

var CountOfAlfaSlices = 100

type T1FS struct {
	Decompose Interval   `json:"decompose"`
	Vert      []Number   `json:"vert"`
	Form      v.Variants `json:"form"`
}

func NewT1FS(Vert ...Number) *T1FS {
	if len(Vert) != 3 && len(Vert) != 4 {
		return nil
	}

	result := &T1FS{
		Decompose: Interval{NumbersMin, NumbersMin},
		Vert:      make([]Number, len(Vert)),
	}

	if len(Vert) == 3 {
		result.Form = v.Triangle
	} else {
		result.Form = v.Trapezoid
	}

	for i := range Vert {
		result.Vert[i] = Vert[i]
	}

	return result
}

func (t *T1FS) CopyEval() Rating {
	return Rating{NewT1FS(t.Vert...)}
}

func (t *T1FS) MemberFunction(alpha Number) Interval {
	if len(t.Vert) == 3 {
		return Interval{t.Vert[0] + (t.Vert[1]-t.Vert[0])*alpha, t.Vert[2] - (t.Vert[2]-t.Vert[1])*alpha}
	} else {
		return Interval{t.Vert[0] + (t.Vert[1]-t.Vert[0])*alpha, t.Vert[3] - (t.Vert[3]-t.Vert[2])*alpha}
	}
}

func (t *T1FS) GetType() string {
	return reflect.TypeOf(t).String()
}

func (t *T1FS) ConvertToNumber() Number {
	i := t.ConvertToInterval()
	return (i.Start + i.End) / 2
}

func (t *T1FS) ConvertToInterval() Interval {
	if t.Decompose.Start == NumbersMin && t.Decompose.End == NumbersMin {
		t.Decompose = Interval{0, 0}
		for alpha := Number(0.0); alpha <= 1; alpha += Number(1.0 / float64(CountOfAlfaSlices)) {
			t.Decompose = t.Decompose.Sum(t.MemberFunction(alpha).Weighted(alpha)).ConvertToInterval()
		}
	}
	return t.Decompose
}

func (t *T1FS) ConvertToT1FS(f v.Variants) *T1FS {
	if f == t.Form || f == v.Default {
		return t
	} else if f == v.Triangle && t.Form == v.Trapezoid {
		return NewT1FS(t.Vert[0], (t.Vert[1]+t.Vert[2])/2, t.Vert[3])
	} else {
		return NewT1FS(t.Vert[0], t.Vert[1], t.Vert[1], t.Vert[2])
	}
}

func (t *T1FS) ConvertToAIFS(f v.Variants) *AIFS {
	if f == t.Form || f == v.Default {
		return NewAIFS(0.0, t.Vert...)
	} else if f == v.Triangle && t.Form == v.Trapezoid {
		return NewAIFS(0.0, t.Vert[0], (t.Vert[1]+t.Vert[2])/2, t.Vert[3])
	} else {
		return NewAIFS(0.0, t.Vert[0], t.Vert[1], t.Vert[1], t.Vert[2])
	}
}

func (t *T1FS) ConvertToIT2FS(f v.Variants) *IT2FS {
	if (f == t.Form || f == v.Default) && t.Form == v.Triangle {
		return NewIT2FS([]Interval{{t.Vert[0], t.Vert[0]}, {t.Vert[2], t.Vert[2]}},
			[]Number{t.Vert[1]})
	} else if (f == t.Form || f == v.Default) && t.Form == v.Trapezoid {
		return NewIT2FS([]Interval{{t.Vert[0], t.Vert[0]}, {t.Vert[3], t.Vert[3]}},
			[]Number{t.Vert[1], t.Vert[2]})
	} else if f == v.Triangle && t.Form == v.Trapezoid {
		return NewIT2FS([]Interval{{t.Vert[0], t.Vert[0]}, {t.Vert[3], t.Vert[3]}},
			[]Number{(t.Vert[1] + t.Vert[2]) / 2})
	} else {
		return NewIT2FS([]Interval{{t.Vert[0], t.Vert[0]}, {t.Vert[2], t.Vert[2]}},
			[]Number{t.Vert[1], t.Vert[1]})
	}
}

func (t *T1FS) GetForm() v.Variants {
	return t.Form
}

func (t *T1FS) Weighted(Weight Evaluated) Rating {
	wt := NewT1FS(t.Vert...)
	for i := range t.Vert {
		wt.Vert[i] = t.Vert[i].Weighted(Weight).ConvertToNumber()
	}

	return Rating{wt}
}

func (t *T1FS) String() string {
	return fmt.Sprint(t.Vert)
}

func (t *T1FS) DiffNumber(other Evaluated, variants v.Variants) (Number, error) {
	if other.GetType() == NumbersMin.GetType() {
		i := t.ConvertToInterval()
		return i.DiffNumber(other, variants)
	} else if other.GetType() == t.GetType() {
		d := Number(0)
		if variants == v.SqrtDistance {
			for i := range t.Vert {
				d += Number(math.Pow(float64(t.Vert[i]-other.ConvertToT1FS(v.Default).Vert[i]), 2))
			}
		} else if variants == v.CbrtDistance {
			for i := range t.Vert {
				d += Number(math.Abs(math.Pow(float64(t.Vert[i]-other.ConvertToT1FS(v.Default).Vert[i]), 3)))
			}
		} else {
			return 0, v.InvalidCaseOfOperation
		}
		return d / Number(len(t.Vert)), nil
	} else {
		return 0, v.IncompatibleTypes
	}
}

func (t *T1FS) DiffInterval(other Interval, typeOfCriterion bool, variants v.Variants) (Interval, error) {
	i := t.ConvertToInterval()

	if d, err := i.DiffInterval(other, typeOfCriterion, variants); err != nil {
		return Interval{}, errors.Join(err)
	} else {
		return d, nil
	}
}

func (t *T1FS) Sum(other Evaluated) Rating {
	ret := NewT1FS(t.Vert...)
	for i := range t.Vert {
		ret.Vert[i] = t.Vert[i] + other.ConvertToT1FS(v.Default).Vert[i]
	}
	return Rating{ret}
}
