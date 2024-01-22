package matrix

import (
	"errors"
	"fmt"
	"math"
	"reflect"
)

var CountOfAlfaSlices = 100

const (
	Triangle  = 1
	Trapezoid = 2
	Default   = 3
)

//type FuzzySets interface {
//	Evaluated
//	MemberFunction(alpha Number) Interval
//}

type T1FS struct {
	decompose Interval
	alphas    []Number
	vert      []Number
	form      Variants
}

func NewT1FS(vert ...Number) *T1FS {
	if len(vert) != 3 && len(vert) != 4 {
		return nil
	}

	result := &T1FS{
		decompose: Interval{NumbersMin, NumbersMin},
		alphas:    make([]Number, 0),
		vert:      make([]Number, len(vert)),
	}

	if len(vert) == 3 {
		result.form = Triangle
	} else {
		result.form = Trapezoid
	}

	for i := range vert {
		result.vert[i] = vert[i]
	}

	return result
}

func (t *T1FS) MemberFunction(alpha Number) Interval {
	if len(t.vert) == 3 {
		return Interval{t.vert[0] + (t.vert[1]-t.vert[0])*alpha, t.vert[2] - (t.vert[2]-t.vert[1])*alpha}
	} else {
		return Interval{t.vert[0] + (t.vert[1]-t.vert[0])*alpha, t.vert[3] - (t.vert[3]-t.vert[2])*alpha}
	}
}

func (t *T1FS) ConvertToNumbers() Number {
	i := t.ConvertToInterval()
	return (i.Start + i.End) / 2
}

func (t *T1FS) ConvertToInterval() Interval {
	if t.decompose.Start == NumbersMin && t.decompose.End == NumbersMin {
		t.decompose = Interval{0, 0}
		for alpha := Number(0.0); alpha <= 1; alpha += Number(1.0 / float64(CountOfAlfaSlices)) {
			if tmp, err := t.decompose.Sum(t.MemberFunction(alpha).Weighted(alpha)); err != nil {
				return Interval{}
			} else {
				t.decompose = tmp.ConvertToInterval()
			}
		}
	}
	return t.decompose
}

func (t *T1FS) ConvertToT1FS(v Variants) *T1FS {
	if (v == Triangle && t.form == Triangle) || (v == Trapezoid && t.form == Trapezoid) || v == Default {
		return t
	} else if v == Triangle && t.form == Trapezoid {
		return NewT1FS(t.vert[0], (t.vert[1]+t.vert[2])/2, t.vert[3])
	} else {
		return NewT1FS(t.vert[0], t.vert[1], t.vert[1], t.vert[2])
	}
}

func (t *T1FS) Weighted(weight Evaluated) Evaluated {
	for i := range t.vert {
		t.vert[i] = t.vert[i].Weighted(weight).ConvertToNumbers()
	}

	return t
}

func (t *T1FS) String() string {
	return fmt.Sprint(t.vert)
}

func (t *T1FS) DistanceNumber(other Evaluated, v Variants) (Number, error) {
	if reflect.TypeOf(other) == reflect.TypeOf(NumbersMin) {
		i := t.ConvertToInterval()
		return i.DistanceNumber(other, v)
	} else if reflect.TypeOf(other) == reflect.TypeOf(t) {
		d := Number(0)
		if v == SqrtDistance {
			for i := range t.vert {
				d += Number(math.Pow(float64(t.vert[i]-other.ConvertToT1FS(Default).vert[i]), 2))
			}
		} else if v == CbrtDistance {
			for i := range t.vert {
				d += Number(math.Abs(math.Pow(float64(t.vert[i]-other.ConvertToT1FS(Default).vert[i]), 3)))
			}
		} else {
			return 0, errors.New("incomplete case of distance")
		}
		return d / Number(len(t.vert)), nil
	} else {
		return 0, errors.New("incomparable values")
	}
}

func (t *T1FS) DistanceInterval(other Interval, typeOfCriterion bool, v Variants) (Interval, error) {
	i := t.ConvertToInterval()

	if d, err := i.DistanceInterval(other, typeOfCriterion, v); err != nil {
		return Interval{}, errors.Join(err)
	} else {
		return d, nil
	}
}

func (t *T1FS) Sum(other Evaluated) (Evaluated, error) {
	if reflect.TypeOf(other) != reflect.TypeOf(t) {
		return nil, errors.New("can't sum difficult types of evaluated")
	} else {
		for i := range t.vert {
			t.vert[i] += other.ConvertToT1FS(Default).vert[i]
		}
		return t, nil
	}
}

func positiveIdealRateT1FS(x, y *T1FS, typeOfCriterion bool) *T1FS {
	ideal := NewT1FS(x.vert...)

	if typeOfCriterion == Benefit {
		for i := range ideal.vert {
			ideal.vert[i] = Number(math.Max(float64(x.vert[i]), float64(y.vert[i])))
		}
	} else {
		for i := range ideal.vert {
			ideal.vert[i] = Number(math.Min(float64(x.vert[i]), float64(y.vert[i])))
		}
	}

	return ideal
}

func negativeIdealRateT1FS(x, y *T1FS, typeOfCriterion bool) *T1FS {
	return positiveIdealRateT1FS(x, y, !typeOfCriterion)
}
