package eval

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	v "webApp/lib/variables"
)

type IT2FS struct {
	Decompose Interval   `json:"decompose"`
	Bottom    []Interval `json:"bottom"`
	Upward    []Number   `json:"upward"`
	Form      v.Variants `json:"form"`
}

func NewIT2FS(Bottom []Interval, Upward []Number) *IT2FS {
	if len(Bottom) != 2 || len(Upward) > 2 || len(Upward) == 0 {
		return nil
	}

	result := &IT2FS{
		Decompose: Interval{NumbersMin, NumbersMin},
		Bottom:    make([]Interval, len(Bottom)),
		Upward:    make([]Number, len(Upward)),
	}

	if len(Upward) == 1 {
		result.Form = v.Triangle
	} else {
		result.Form = v.Trapezoid
	}

	for i := range Bottom {
		result.Bottom[i] = Bottom[i]
	}

	for i := range Upward {
		result.Upward[i] = Upward[i]
	}

	return result
}

func (t *IT2FS) CopyEval() Rating {
	return Rating{NewIT2FS(t.Bottom, t.Upward)}
}

func (t *IT2FS) MemberFunction(alpha Number) (Interval, Interval) {
	if len(t.Upward) == 1 {
		return Interval{t.Bottom[0].Start + (t.Upward[0]-t.Bottom[0].Start)*alpha,
				t.Bottom[1].Start - (t.Bottom[1].Start-t.Upward[0])*alpha},
			Interval{t.Bottom[0].End - (t.Upward[0]-t.Bottom[0].End)*alpha,
				t.Bottom[1].End - (t.Bottom[1].End-t.Upward[0])*alpha}
	} else {
		return Interval{t.Bottom[0].Start + (t.Upward[0]-t.Bottom[0].Start)*alpha,
				t.Bottom[1].Start - (t.Bottom[1].Start-t.Upward[1])*alpha},
			Interval{t.Bottom[0].End - (t.Upward[0]-t.Bottom[0].End)*alpha,
				t.Bottom[1].End - (t.Bottom[1].End-t.Upward[1])*alpha}
	}
}

func (t *IT2FS) GetType() string {
	return reflect.TypeOf(t).String()
}

func (t *IT2FS) ConvertToNumber() Number {
	i := t.ConvertToInterval()
	return (i.Start + i.End) / 2
}

func (t *IT2FS) ConvertToInterval() Interval {
	if t.Decompose.Start == NumbersMin && t.Decompose.End == NumbersMin {
		t.Decompose = Interval{0, 0}
		for alpha := Number(0.0); alpha <= 1; alpha += Number(1.0 / float64(CountOfAlfaSlices)) {
			s, f := t.MemberFunction(alpha)
			s = s.Weighted(Rating{alpha * 0.5}).ConvertToInterval()
			f = f.Weighted(Rating{alpha * 0.5}).ConvertToInterval()
			t.Decompose = t.Decompose.Sum(s.Sum(f).Weighted(alpha)).ConvertToInterval()
		}
	}
	return t.Decompose
}

func (t *IT2FS) ConvertToT1FS(_ v.Variants) *T1FS {
	fmt.Println("Call deprecated method: ConvertToT1FS Form IT2FS")
	return nil
}

func (t *IT2FS) ConvertToAIFS(_ v.Variants) *AIFS {
	fmt.Println("Call deprecated method: ConvertToAIFS Form IT2FS")
	return nil
}

func (t *IT2FS) ConvertToIT2FS(f v.Variants) *IT2FS {
	if (f == v.Triangle && t.Form == v.Triangle) || (f == v.Trapezoid && t.Form == v.Trapezoid) || f == v.Default {
		return t
	} else if f == v.Triangle && t.Form == v.Trapezoid {
		return NewIT2FS(t.Bottom, []Number{(t.Upward[0] + t.Upward[1]) / 2})
	} else {
		return NewIT2FS(t.Bottom, []Number{t.Upward[0], t.Upward[0]})
	}
}

func (t *IT2FS) GetForm() v.Variants {
	return t.Form
}

func (t *IT2FS) Weighted(Weight Evaluated) Rating {
	wt := NewIT2FS(t.Bottom, t.Upward)
	for i := range t.Bottom {
		wt.Bottom[i] = t.Bottom[i].Weighted(Weight).ConvertToInterval()
	}

	for i := range t.Upward {
		wt.Upward[i] = t.Upward[i].Weighted(Weight).ConvertToNumber()
	}

	return Rating{wt}
}

func (t *IT2FS) String() string {
	if len(t.Upward) == 1 {
		return "{" + fmt.Sprint(t.Bottom[0]) + ", " + fmt.Sprint(t.Upward[0]) + ", " + fmt.Sprint(t.Bottom[1]) + "}"
	} else {
		return "{" + fmt.Sprint(t.Bottom[0]) + ", " + fmt.Sprint(t.Upward[0]) + ", " +
			fmt.Sprint(t.Upward[1]) + fmt.Sprint(t.Bottom[1]) + "}"
	}
}

func (t *IT2FS) DiffNumber(other Evaluated, variants v.Variants) (Number, error) {
	if other.GetType() == NumbersMin.GetType() {
		i := t.ConvertToInterval()
		return i.DiffNumber(other, variants)
	} else if other.GetType() == (&T1FS{}).GetType() {
		d := Number(0)
		var power float64
		if variants == v.SqrtDistance {
			power = 2
		} else if variants == v.CbrtDistance {
			power = 3
		} else {
			return 0, v.InvalidCaseOfOperation
		}
		d += Number(math.Pow(float64(t.Bottom[0].Start-other.ConvertToT1FS(v.Default).Vert[0]), power))
		d += Number(math.Pow(float64(t.Bottom[0].End-other.ConvertToT1FS(v.Default).Vert[0]), power))
		if t.Form == v.Triangle {
			d += Number(math.Pow(float64(t.Upward[0]-other.ConvertToT1FS(v.Default).Vert[1]), power))
			d += Number(math.Pow(float64(t.Bottom[1].Start-other.ConvertToT1FS(v.Default).Vert[2]), power))
			d += Number(math.Pow(float64(t.Bottom[1].End-other.ConvertToT1FS(v.Default).Vert[2]), power))
		} else {
			d += Number(math.Pow(float64(t.Upward[0]-other.ConvertToT1FS(v.Default).Vert[1]), power))
			d += Number(math.Pow(float64(t.Upward[1]-other.ConvertToT1FS(v.Default).Vert[2]), power))
			d += Number(math.Pow(float64(t.Bottom[1].Start-other.ConvertToT1FS(v.Default).Vert[3]), power))
			d += Number(math.Pow(float64(t.Bottom[1].End-other.ConvertToT1FS(v.Default).Vert[3]), power))
		}

		return d / Number(len(t.Upward)+4), nil
	} else {
		return 0, v.IncompatibleTypes
	}
}

func (t *IT2FS) DiffInterval(other Interval, typeOfCriterion bool, variants v.Variants) (Interval, error) {
	i := t.ConvertToInterval()

	if d, err := i.DiffInterval(other, typeOfCriterion, variants); err != nil {
		return Interval{}, errors.Join(err)
	} else {
		return d, nil
	}
}

func (t *IT2FS) Sum(other Evaluated) Rating {
	ret := NewIT2FS(t.Bottom, t.Upward)
	for i := range ret.Bottom {
		ret.Bottom[i] = t.Bottom[i].Sum(other.ConvertToIT2FS(v.Default).Bottom[i]).ConvertToInterval()
	}

	for i := range ret.Upward {
		ret.Upward[i] = t.Upward[i].Sum(other.ConvertToIT2FS(v.Default).Upward[i]).ConvertToNumber()
	}

	return Rating{ret}
}
