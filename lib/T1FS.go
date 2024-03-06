package lib

import (
	"context"
	"errors"
	"fmt"
	"math"
	"reflect"
	"sync"
)

var CountOfAlfaSlices = 100

const (
	None      = 0
	Triangle  = 1
	Trapezoid = 2
)

type T1FS struct {
	decompose Interval
	vert      []Number
	form      Variants
}

func NewT1FS(vert ...Number) *T1FS {
	if len(vert) != 3 && len(vert) != 4 {
		return nil
	}

	result := &T1FS{
		decompose: Interval{NumbersMin, NumbersMin},
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

func (t *T1FS) ConvertToNumber() Number {
	i := t.ConvertToInterval()
	return (i.Start + i.End) / 2
}

func (t *T1FS) ConvertToInterval() Interval {
	if t.decompose.Start == NumbersMin && t.decompose.End == NumbersMin {
		t.decompose = Interval{0, 0}
		for alpha := Number(0.0); alpha <= 1; alpha += Number(1.0 / float64(CountOfAlfaSlices)) {
			t.decompose = t.decompose.Sum(t.MemberFunction(alpha).Weighted(alpha)).ConvertToInterval()
		}
	}
	return t.decompose
}

func (t *T1FS) ConvertToT1FS(v Variants) *T1FS {
	if v == t.form || v == Default {
		return t
	} else if v == Triangle && t.form == Trapezoid {
		return NewT1FS(t.vert[0], (t.vert[1]+t.vert[2])/2, t.vert[3])
	} else {
		return NewT1FS(t.vert[0], t.vert[1], t.vert[1], t.vert[2])
	}
}

func (t *T1FS) ConvertToAIFS(f Variants) *AIFS {
	if f == t.form || f == Default {
		return NewAIFS(0.0, t.vert...)
	} else if f == Triangle && t.form == Trapezoid {
		return NewAIFS(0.0, t.vert[0], (t.vert[1]+t.vert[2])/2, t.vert[3])
	} else {
		return NewAIFS(0.0, t.vert[0], t.vert[1], t.vert[1], t.vert[2])
	}
}

func (t *T1FS) ConvertToIT2FS(v Variants) *IT2FS {
	if (v == t.form || v == Default) && t.form == Triangle {
		return NewIT2FS([]Interval{{t.vert[0], t.vert[0]}, {t.vert[2], t.vert[2]}},
			[]Number{t.vert[1]})
	} else if (v == t.form || v == Default) && t.form == Trapezoid {
		return NewIT2FS([]Interval{{t.vert[0], t.vert[0]}, {t.vert[3], t.vert[3]}},
			[]Number{t.vert[1], t.vert[2]})
	} else if v == Triangle && t.form == Trapezoid {
		return NewIT2FS([]Interval{{t.vert[0], t.vert[0]}, {t.vert[3], t.vert[3]}},
			[]Number{(t.vert[1] + t.vert[2]) / 2})
	} else {
		return NewIT2FS([]Interval{{t.vert[0], t.vert[0]}, {t.vert[2], t.vert[2]}},
			[]Number{t.vert[1], t.vert[1]})
	}
}

func (t *T1FS) GetForm() Variants {
	return t.form
}

func (t *T1FS) Weighted(weight Evaluated) Evaluated {
	wt := NewT1FS(t.vert...)
	for i := range t.vert {
		wt.vert[i] = t.vert[i].Weighted(weight).ConvertToNumber()
	}

	return wt
}

func (t *T1FS) String() string {
	return fmt.Sprint(t.vert)
}

func (t *T1FS) DiffNumber(other Evaluated, v Variants) (Number, error) {
	if reflect.TypeOf(other) == reflect.TypeOf(NumbersMin) {
		i := t.ConvertToInterval()
		return i.DiffNumber(other, v)
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
			return 0, InvalidCaseOfOperation
		}
		return d / Number(len(t.vert)), nil
	} else {
		return 0, IncompatibleTypes
	}
}

func (t *T1FS) DiffInterval(other Interval, typeOfCriterion bool, v Variants) (Interval, error) {
	i := t.ConvertToInterval()

	if d, err := i.DiffInterval(other, typeOfCriterion, v); err != nil {
		return Interval{}, errors.Join(err)
	} else {
		return d, nil
	}
}

func (t *T1FS) Sum(other Evaluated) Evaluated {
	ret := NewT1FS(t.vert...)
	for i := range t.vert {
		ret.vert[i] = t.vert[i] + other.ConvertToT1FS(Default).vert[i]
	}
	return ret
}

func positiveIdealRateT1FS(alts []Alternative, criteria []Criterion, form Variants) (*Alternative, error) {
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
					if reflect.TypeOf(alts[j].grade[i]) != reflect.TypeOf(&T1FS{}) &&
						reflect.TypeOf(alts[j].grade[i]) != reflect.TypeOf(&IT2FS{}) {
						cancel()
						err = IncompatibleTypes
						return
					}

					if positive.grade[i] == nil && c.typeOfCriteria == Benefit {
						if form == Triangle {
							positive.grade[i] = NewT1FS(NumbersMin, NumbersMin, NumbersMin)
						} else {
							positive.grade[i] = NewT1FS(NumbersMin, NumbersMin, NumbersMin, NumbersMin)
						}
					} else if positive.grade[i] == nil && c.typeOfCriteria == Cost {
						if form == Triangle {
							positive.grade[i] = NewT1FS(NumbersMax, NumbersMax, NumbersMax)
						} else {
							positive.grade[i] = NewT1FS(NumbersMax, NumbersMax, NumbersMax, NumbersMax)
						}
					}

					if reflect.TypeOf(alts[j].grade[i]) == reflect.TypeOf(&T1FS{}) {
						if c.typeOfCriteria == Benefit {
							positive.grade[i] = Max(positive.grade[i], alts[j].grade[i])
						} else {
							positive.grade[i] = Min(positive.grade[i], alts[j].grade[i])
						}
					} else {
						posGrade := positive.grade[i].ConvertToT1FS(Default)
						altsGrade := alts[j].grade[i].ConvertToIT2FS(Default)

						if c.typeOfCriteria == Benefit {
							posGrade.vert[0] =
								Max(posGrade.vert[0], altsGrade.bottom[0].End).ConvertToNumber()
							posGrade.vert[1] =
								Max(posGrade.vert[1], altsGrade.upward[0]).ConvertToNumber()
							if form == Triangle {
								posGrade.vert[2] =
									Max(posGrade.vert[2], altsGrade.bottom[1].End).ConvertToNumber()
							} else {
								posGrade.vert[2] = Max(posGrade.vert[2], altsGrade.upward[1]).ConvertToNumber()
								posGrade.vert[3] = Max(posGrade.vert[3], altsGrade.bottom[1].End).ConvertToNumber()
							}
						} else {
							posGrade.vert[0] = Min(posGrade.vert[0], altsGrade.bottom[0].Start).ConvertToNumber()
							posGrade.vert[1] = Min(posGrade.vert[1], altsGrade.upward[0]).ConvertToNumber()
							if form == Triangle {
								posGrade.vert[2] = Min(posGrade.vert[2], altsGrade.bottom[1].Start).ConvertToNumber()
							} else {
								posGrade.vert[2] = Min(posGrade.vert[2], altsGrade.upward[1]).ConvertToNumber()
								posGrade.vert[3] = Min(posGrade.vert[3], altsGrade.bottom[1].Start).ConvertToNumber()
							}
						}
					}
				}
			}
		}(i, c)
	}
	wg.Wait()
	return positive, err
}

func negativeIdealRateT1FS(alts []Alternative, criteria []Criterion, form Variants) (*Alternative, error) {
	return positiveIdealRateT1FS(alts, ChangeTypes(criteria), form)
}
