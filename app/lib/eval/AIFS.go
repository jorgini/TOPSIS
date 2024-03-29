package eval

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	v "webApp/lib/variables"
)

type AIFS struct {
	*T1FS `json:"t1fs"`
	Pi    Number `json:"pi"`
}

func NewAIFS(Pi Number, Vert ...Number) *AIFS {
	if len(Vert) != 3 && len(Vert) != 4 {
		return nil
	}

	result := &AIFS{
		T1FS: NewT1FS(Vert...),
		Pi:   Pi,
	}

	return result
}

func (a *AIFS) CopyEval() Rating {
	return Rating{NewAIFS(a.Pi, a.Vert...)}
}

func (a *AIFS) MemberFunction(alpha Number) Interval {
	if alpha > 1.0-a.Pi {
		return Interval{0, 0}
	}

	if len(a.Vert) == 3 {
		return Interval{a.Vert[0] + (a.Vert[1]-a.Vert[0])*alpha/(1.0-a.Pi),
			a.Vert[2] - (a.Vert[2]-a.Vert[1])*alpha/(1.0-a.Pi)}
	} else {
		return Interval{a.Vert[0] + (a.Vert[1]-a.Vert[0])*alpha/(1-a.Pi),
			a.Vert[3] - (a.Vert[3]-a.Vert[2])*alpha/(1.0-a.Pi)}
	}
}

func (a *AIFS) GetType() string {
	return reflect.TypeOf(a).String()
}

func (a *AIFS) ConvertToNumber() Number {
	i := a.ConvertToInterval()
	return (i.Start + i.End) / 2
}

func (a *AIFS) ConvertToInterval() Interval {
	if a.Decompose.Start == NumbersMin && a.Decompose.End == NumbersMin {
		a.Decompose = Interval{0, 0}
		for alpha := Number(0.0); alpha <= 1-a.Pi; alpha += (1.0 - a.Pi) / Number(CountOfAlfaSlices) {
			a.Decompose = a.Decompose.Sum(a.MemberFunction(alpha).Weighted(Rating{alpha})).ConvertToInterval()
		}
	}
	return a.Decompose
}

func (a *AIFS) ConvertToT1FS(f v.Variants) *T1FS {
	if f == a.Form || f == v.Default {
		return NewT1FS(a.Vert...)
	} else if f == v.Triangle && a.Form == v.Trapezoid {
		return NewT1FS(a.Vert[0], (a.Vert[1]+a.Vert[2])/2, a.Vert[3])
	} else {
		return NewT1FS(a.Vert[0], a.Vert[1], a.Vert[1], a.Vert[2])
	}
}

func (a *AIFS) ConvertToAIFS(f v.Variants) *AIFS {
	if f == a.Form || f == v.Default {
		return a
	} else if f == v.Triangle && a.Form == v.Trapezoid {
		return NewAIFS(a.Pi, a.Vert[0], (a.Vert[1]+a.Vert[2])/2, a.Vert[3])
	} else {
		return NewAIFS(a.Pi, a.Vert[0], a.Vert[1], a.Vert[1], a.Vert[2])
	}
}

func (a *AIFS) ConvertToIT2FS(f v.Variants) *IT2FS {
	delta1 := a.Pi * (a.Vert[1] - a.Vert[0]) / (2 * (1 - a.Pi))
	delta2 := a.Pi * (a.Vert[len(a.Vert)-1] - a.Vert[len(a.Vert)-2]) / (2 * (1 - a.Pi))

	if a.Form == f || f == v.Default {
		if a.Form == v.Triangle {
			return NewIT2FS([]Interval{{a.Vert[0] - delta1, a.Vert[0] + delta1},
				{a.Vert[2] - delta2, a.Vert[2] + delta2}}, []Number{a.Vert[1]})
		} else {
			return NewIT2FS([]Interval{{a.Vert[0] - delta1, a.Vert[0] + delta1},
				{a.Vert[3] - delta2, a.Vert[3] + delta2}}, []Number{a.Vert[1], a.Vert[2]})
		}
	} else if a.Form == v.Triangle && f == v.Trapezoid {
		return NewIT2FS([]Interval{{a.Vert[0] - delta1, a.Vert[0] + delta1},
			{a.Vert[2] - delta2, a.Vert[2] + delta2}}, []Number{a.Vert[1], a.Vert[1]})
	} else if a.Form == v.Trapezoid && f == v.Triangle {
		return NewIT2FS([]Interval{{a.Vert[0] - delta1, a.Vert[0] + delta1},
			{a.Vert[3] - delta2, a.Vert[3] + delta2}}, []Number{(a.Vert[1] + a.Vert[2]) / 2})
	} else {
		return nil
	}
}

func (a *AIFS) GetForm() v.Variants {
	return a.Form
}

func (a *AIFS) Weighted(Weight Evaluated) Rating {
	wt := NewAIFS(a.Pi, a.Vert...)
	for i := range a.Vert {
		wt.Vert[i] = a.Vert[i].Weighted(Weight).ConvertToNumber()
	}

	return Rating{wt}
}

func (a *AIFS) String() string {
	s := "["
	for _, num := range a.Vert {
		s += fmt.Sprint(num) + " "
	}
	s += "Pi=" + fmt.Sprint(a.Pi) + "]"
	return s
}

func (a *AIFS) DiffNumber(other Evaluated, variants v.Variants) (Number, error) {
	if other.GetType() == NumbersMin.GetType() {
		i := a.ConvertToInterval()
		return i.DiffNumber(other, variants)
	} else if other.GetType() == a.GetType() {
		d := Number(0)
		if variants == v.SqrtDistance {
			for i := range a.Vert {
				d += Number(math.Pow(float64(a.Vert[i]-other.ConvertToAIFS(v.Default).Vert[i]), 2))
			}
			d += 0.01 * Number(math.Pow(float64(a.Pi-other.ConvertToAIFS(v.Default).Pi), 2))
		} else if variants == v.CbrtDistance {
			for i := range a.Vert {
				d += Number(math.Abs(math.Pow(float64(a.Vert[i]-other.ConvertToAIFS(v.Default).Vert[i]), 3)))
			}
			d += 0.001 * Number(math.Abs(math.Pow(float64(a.Pi-other.ConvertToAIFS(v.Default).Pi), 3)))
		} else {
			return 0, v.InvalidCaseOfOperation
		}
		return d / (Number(len(a.Vert)) + 1), nil
	} else {
		return 0, v.IncompatibleTypes
	}
}

func (a *AIFS) DiffInterval(other Interval, typeOfCriterion bool, variants v.Variants) (Interval, error) {
	i := a.ConvertToInterval()

	if d, err := i.DiffInterval(other, typeOfCriterion, variants); err != nil {
		return Interval{}, errors.Join(err)
	} else {
		return d, nil
	}
}

func (a *AIFS) Sum(other Evaluated) Rating {
	ret := NewAIFS(a.Pi, a.Vert...)
	for i := range a.Vert {
		ret.Vert[i] = a.Vert[i] + other.ConvertToAIFS(v.Default).Vert[i]
	}
	ret.Pi = a.Pi + other.ConvertToAIFS(v.Default).Pi
	return Rating{ret}
}
