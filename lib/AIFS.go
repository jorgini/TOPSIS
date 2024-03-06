package lib

import (
	"context"
	"errors"
	"fmt"
	"math"
	"reflect"
	"sync"
)

type AIFS struct {
	*T1FS
	pi Number
}

func NewAIFS(pi Number, vert ...Number) *AIFS {
	if len(vert) != 3 && len(vert) != 4 {
		return nil
	}

	result := &AIFS{
		T1FS: NewT1FS(vert...),
		pi:   pi,
	}

	return result
}

func (a *AIFS) MemberFunction(alpha Number) Interval {
	if alpha > 1.0-a.pi {
		return Interval{0, 0}
	}

	if len(a.vert) == 3 {
		return Interval{a.vert[0] + (a.vert[1]-a.vert[0])*alpha/(1.0-a.pi),
			a.vert[2] - (a.vert[2]-a.vert[1])*alpha/(1.0-a.pi)}
	} else {
		return Interval{a.vert[0] + (a.vert[1]-a.vert[0])*alpha/(1-a.pi),
			a.vert[3] - (a.vert[3]-a.vert[2])*alpha/(1.0-a.pi)}
	}
}

func (a *AIFS) ConvertToNumber() Number {
	i := a.ConvertToInterval()
	return (i.Start + i.End) / 2
}

func (a *AIFS) ConvertToInterval() Interval {
	if a.decompose.Start == NumbersMin && a.decompose.End == NumbersMin {
		a.decompose = Interval{0, 0}
		for alpha := Number(0.0); alpha <= 1-a.pi; alpha += (1.0 - a.pi) / Number(CountOfAlfaSlices) {
			a.decompose = a.decompose.Sum(a.MemberFunction(alpha).Weighted(alpha)).ConvertToInterval()
		}
	}
	return a.decompose
}

func (a *AIFS) ConvertToT1FS(f Variants) *T1FS {
	if f == a.form || f == Default {
		return NewT1FS(a.vert...)
	} else if f == Triangle && a.form == Trapezoid {
		return NewT1FS(a.vert[0], (a.vert[1]+a.vert[2])/2, a.vert[3])
	} else {
		return NewT1FS(a.vert[0], a.vert[1], a.vert[1], a.vert[2])
	}
}

func (a *AIFS) ConvertToAIFS(f Variants) *AIFS {
	if f == a.form || f == Default {
		return a
	} else if f == Triangle && a.form == Trapezoid {
		return NewAIFS(a.pi, a.vert[0], (a.vert[1]+a.vert[2])/2, a.vert[3])
	} else {
		return NewAIFS(a.pi, a.vert[0], a.vert[1], a.vert[1], a.vert[2])
	}
}

func (a *AIFS) ConvertToIT2FS(f Variants) *IT2FS {
	delta1 := a.pi * (a.vert[1] - a.vert[0]) / (2 * (1 - a.pi))
	delta2 := a.pi * (a.vert[len(a.vert)-1] - a.vert[len(a.vert)-2]) / (2 * (1 - a.pi))

	if a.form == f || f == Default {
		if a.form == Triangle {
			return NewIT2FS([]Interval{{a.vert[0] - delta1, a.vert[0] + delta1},
				{a.vert[2] - delta2, a.vert[2] + delta2}}, []Number{a.vert[1]})
		} else {
			return NewIT2FS([]Interval{{a.vert[0] - delta1, a.vert[0] + delta1},
				{a.vert[3] - delta2, a.vert[3] + delta2}}, []Number{a.vert[1], a.vert[2]})
		}
	} else if a.form == Triangle && f == Trapezoid {
		return NewIT2FS([]Interval{{a.vert[0] - delta1, a.vert[0] + delta1},
			{a.vert[2] - delta2, a.vert[2] + delta2}}, []Number{a.vert[1], a.vert[1]})
	} else if a.form == Trapezoid && f == Triangle {
		return NewIT2FS([]Interval{{a.vert[0] - delta1, a.vert[0] + delta1},
			{a.vert[3] - delta2, a.vert[3] + delta2}}, []Number{(a.vert[1] + a.vert[2]) / 2})
	} else {
		return nil
	}
}

func (a *AIFS) GetForm() Variants {
	return a.form
}

func (a *AIFS) Weighted(weight Evaluated) Evaluated {
	wt := NewAIFS(a.pi, a.vert...)
	for i := range a.vert {
		wt.vert[i] = a.vert[i].Weighted(weight).ConvertToNumber()
	}

	return wt
}

func (a *AIFS) String() string {
	s := "["
	for number := range a.vert {
		s += fmt.Sprint(number) + " "
	}
	s += "pi=" + fmt.Sprint(a.pi) + "]"
	return s
}

func (a *AIFS) DiffNumber(other Evaluated, v Variants) (Number, error) {
	if reflect.TypeOf(other) == reflect.TypeOf(NumbersMin) {
		i := a.ConvertToInterval()
		return i.DiffNumber(other, v)
	} else if reflect.TypeOf(other) == reflect.TypeOf(a) {
		d := Number(0)
		if v == SqrtDistance {
			for i := range a.vert {
				d += Number(math.Pow(float64(a.vert[i]-other.ConvertToAIFS(Default).vert[i]), 2))
			}
			d += 0.01 * Number(math.Pow(float64(a.pi-other.ConvertToAIFS(Default).pi), 2))
		} else if v == CbrtDistance {
			for i := range a.vert {
				d += Number(math.Abs(math.Pow(float64(a.vert[i]-other.ConvertToAIFS(Default).vert[i]), 3)))
			}
			d += 0.001 * Number(math.Abs(math.Pow(float64(a.pi-other.ConvertToAIFS(Default).pi), 3)))
		} else {
			return 0, InvalidCaseOfOperation
		}
		return d / (Number(len(a.vert)) + 1), nil
	} else {
		return 0, IncompatibleTypes
	}
}

func (a *AIFS) DiffInterval(other Interval, typeOfCriterion bool, v Variants) (Interval, error) {
	i := a.ConvertToInterval()

	if d, err := i.DiffInterval(other, typeOfCriterion, v); err != nil {
		return Interval{}, errors.Join(err)
	} else {
		return d, nil
	}
}

func (a *AIFS) Sum(other Evaluated) Evaluated {
	ret := NewAIFS(a.pi, a.vert...)
	for i := range a.vert {
		ret.vert[i] = a.vert[i] + other.ConvertToAIFS(Default).vert[i]
	}
	ret.pi = a.pi + other.ConvertToAIFS(Default).pi
	return ret
}

func positiveIdealRateAIFS(alts []Alternative, criteria []Criterion) (*Alternative, error) {
	var wg sync.WaitGroup
	var err error = nil
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	positive := &Alternative{make([]Evaluated, len(alts[0].grade)), len(criteria)}

	wg.Add(len(criteria))
	for i, c := range criteria {
		go func(i int, c Criterion) {
			defer wg.Done()

			for j := range alts {
				select {
				case <-ctx.Done():
					return
				default:
					if reflect.TypeOf(alts[j].grade[i]) != reflect.TypeOf(&AIFS{}) {
						cancel()
						err = IncompatibleTypes
						return
					}

					altsGrade := alts[j].grade[i].ConvertToAIFS(Default)

					if positive.grade[i] == nil {
						positive.grade[i] = NewAIFS(altsGrade.pi, altsGrade.vert...)
					}

					if c.typeOfCriteria == Benefit {
						positive.grade[i] = Max(positive.grade[i], alts[j].grade[i])
					} else {
						positive.grade[i] = Min(positive.grade[i], alts[j].grade[i])
					}
				}
			}
		}(i, c)
	}

	wg.Wait()
	return positive, err
}

func negativeIdealRateAIFS(alts []Alternative, criteria []Criterion) (*Alternative, error) {
	return positiveIdealRateAIFS(alts, ChangeTypes(criteria))
}
