package lib

import (
	"errors"
	"fmt"
	"math"
	"reflect"
)

type IT2FS struct {
	decompose Interval
	bottom    []Interval
	upward    []Number
	form      Variants
}

func NewIT2FS(bottom []Interval, upward []Number) *IT2FS {
	if len(bottom) != 2 || len(upward) > 2 || len(upward) == 0 {
		return nil
	}

	result := &IT2FS{
		decompose: Interval{NumbersMin, NumbersMin},
		bottom:    make([]Interval, len(bottom)),
		upward:    make([]Number, len(upward)),
	}

	if len(upward) == 1 {
		result.form = Triangle
	} else {
		result.form = Trapezoid
	}

	for i := range bottom {
		result.bottom[i] = bottom[i]
	}

	for i := range upward {
		result.upward[i] = upward[i]
	}

	return result
}

func (t *IT2FS) MemberFunction(alpha Number) (Interval, Interval) {
	if len(t.upward) == 1 {
		return Interval{t.bottom[0].Start + (t.upward[0]-t.bottom[0].Start)*alpha,
				t.bottom[1].Start - (t.bottom[1].Start-t.upward[0])*alpha},
			Interval{t.bottom[0].End - (t.upward[0]-t.bottom[0].End)*alpha,
				t.bottom[1].End - (t.bottom[1].End-t.upward[0])*alpha}
	} else {
		return Interval{t.bottom[0].Start + (t.upward[0]-t.bottom[0].Start)*alpha,
				t.bottom[1].Start - (t.bottom[1].Start-t.upward[1])*alpha},
			Interval{t.bottom[0].End - (t.upward[0]-t.bottom[0].End)*alpha,
				t.bottom[1].End - (t.bottom[1].End-t.upward[1])*alpha}
	}
}

func (t *IT2FS) ConvertToNumber() Number {
	i := t.ConvertToInterval()
	return (i.Start + i.End) / 2
}

func (t *IT2FS) ConvertToInterval() Interval {
	if t.decompose.Start == NumbersMin && t.decompose.End == NumbersMin {
		t.decompose = Interval{0, 0}
		for alpha := Number(0.0); alpha <= 1; alpha += Number(1.0 / float64(CountOfAlfaSlices)) {
			s, f := t.MemberFunction(alpha)
			s = s.Weighted(alpha * 0.5).ConvertToInterval()
			f = f.Weighted(alpha * 0.5).ConvertToInterval()
			t.decompose = t.decompose.Sum(s.Sum(f).Weighted(alpha)).ConvertToInterval()
		}
	}
	return t.decompose
}

func (t *IT2FS) ConvertToT1FS(f Variants) *T1FS {
	fmt.Println("Call deprecated method")
	return nil
}

func (t *IT2FS) ConvertToAIFS(f Variants) *AIFS {
	fmt.Println("Call deprecated method")
	return nil
}

func (t *IT2FS) ConvertToIT2FS(v Variants) *IT2FS {
	if (v == Triangle && t.form == Triangle) || (v == Trapezoid && t.form == Trapezoid) || v == Default {
		return t
	} else if v == Triangle && t.form == Trapezoid {
		return NewIT2FS(t.bottom, []Number{(t.upward[0] + t.upward[1]) / 2})
	} else {
		return NewIT2FS(t.bottom, []Number{t.upward[0], t.upward[0]})
	}
}

func (t *IT2FS) Weighted(weight Evaluated) Evaluated {
	wt := NewIT2FS(t.bottom, t.upward)
	for i := range t.bottom {
		wt.bottom[i] = t.bottom[i].Weighted(weight).ConvertToInterval()
	}

	for i := range t.upward {
		wt.upward[i] = t.upward[i].Weighted(weight).ConvertToNumber()
	}

	return wt
}

func (t *IT2FS) String() string {
	if len(t.upward) == 1 {
		return "{" + fmt.Sprint(t.bottom[0]) + ", " + fmt.Sprint(t.upward[0]) + ", " + fmt.Sprint(t.bottom[1]) + "}"
	} else {
		return "{" + fmt.Sprint(t.bottom[0]) + ", " + fmt.Sprint(t.upward[0]) + ", " +
			fmt.Sprint(t.upward[1]) + fmt.Sprint(t.bottom[1]) + "}"
	}
}

func (t *IT2FS) DiffNumber(other Evaluated, v Variants) (Number, error) {
	if reflect.TypeOf(other) == reflect.TypeOf(NumbersMin) {
		i := t.ConvertToInterval()
		return i.DiffNumber(other, v)
	} else if reflect.TypeOf(other) == reflect.TypeOf(&T1FS{}) {
		d := Number(0)
		var power float64
		if v == SqrtDistance {
			power = 2
		} else if v == CbrtDistance {
			power = 3
		} else {
			return 0, InvalidCaseOfOperation
		}
		d += Number(math.Pow(float64(t.bottom[0].Start-other.ConvertToT1FS(Default).vert[0]), power))
		d += Number(math.Pow(float64(t.bottom[0].End-other.ConvertToT1FS(Default).vert[0]), power))
		if t.form == Triangle {
			d += Number(math.Pow(float64(t.upward[0]-other.ConvertToT1FS(Default).vert[1]), power))
			d += Number(math.Pow(float64(t.bottom[1].Start-other.ConvertToT1FS(Default).vert[2]), power))
			d += Number(math.Pow(float64(t.bottom[1].End-other.ConvertToT1FS(Default).vert[2]), power))
		} else {
			d += Number(math.Pow(float64(t.upward[0]-other.ConvertToT1FS(Default).vert[1]), power))
			d += Number(math.Pow(float64(t.upward[1]-other.ConvertToT1FS(Default).vert[2]), power))
			d += Number(math.Pow(float64(t.bottom[1].Start-other.ConvertToT1FS(Default).vert[3]), power))
			d += Number(math.Pow(float64(t.bottom[1].End-other.ConvertToT1FS(Default).vert[3]), power))
		}

		return d / Number(len(t.upward)+4), nil
	} else {
		return 0, IncompatibleTypes
	}
}

func (t *IT2FS) DiffInterval(other Interval, typeOfCriterion bool, v Variants) (Interval, error) {
	i := t.ConvertToInterval()

	if d, err := i.DiffInterval(other, typeOfCriterion, v); err != nil {
		return Interval{}, errors.Join(err)
	} else {
		return d, nil
	}
}

func (t *IT2FS) Sum(other Evaluated) Evaluated {
	ret := NewIT2FS(t.bottom, t.upward)
	for i := range ret.bottom {
		ret.bottom[i] = t.bottom[i].Sum(other.ConvertToIT2FS(Default).bottom[i]).ConvertToInterval()
	}

	for i := range ret.upward {
		ret.upward[i] = t.upward[i].Sum(other.ConvertToIT2FS(Default).upward[i]).ConvertToNumber()
	}

	return ret
}
