package lib

import (
	"context"
	"math"
	"reflect"
	"sync"
)

type Alternative struct {
	grade           []Evaluated
	countOfCriteria int
}

func (a *Alternative) String() string {
	s := "[ "
	for i := 0; i < a.countOfCriteria; i++ {
		s += a.grade[i].String() + " "
	}
	return s + " ]"
}

func (a *Alternative) NumberMetric(to *Alternative, v Variants) (Number, error) {
	if a.countOfCriteria != to.countOfCriteria {
		return 0, InvalidSize
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result := Number(0)
	wg.Add(a.countOfCriteria)
	for i := 0; i < a.countOfCriteria; i++ {
		go func(i int) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				if tmp, innerError := a.grade[i].DiffNumber(to.grade[i].ConvertToNumber(), v); innerError != nil {
					err = innerError
					cancel()
					return
				} else {
					mu.Lock()
					result += tmp
					mu.Unlock()
				}
			}
		}(i)
	}
	wg.Wait()

	if err != nil {
		return 0, err
	}

	if v == SqrtDistance {
		result = Number(math.Sqrt(float64(result)))
	} else {
		result = Number(math.Cbrt(float64(result)))
	}

	return result, nil
}

func (a *Alternative) IntervalMetric(to *Alternative, criteria []Criterion, v Variants) (Interval, error) {
	if a.countOfCriteria != to.countOfCriteria {
		return Interval{}, InvalidSize
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result := Interval{Number(0), Number(0)}
	wg.Add(len(criteria))
	for i, c := range criteria {
		go func(i int, c Criterion) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				if tmp, innerError := a.grade[i].DiffInterval(to.grade[i].ConvertToInterval(), c.typeOfCriteria, v); innerError != nil {
					err = innerError
					cancel()
					return
				} else {
					mu.Lock()
					result = result.Sum(tmp).ConvertToInterval()
					mu.Unlock()
				}
			}
		}(i, c)
	}

	wg.Wait()

	if err != nil {
		return Interval{}, err
	}

	if v == SqrtDistance {
		result = Interval{Number(math.Sqrt(float64(result.ConvertToInterval().Start))),
			Number(math.Sqrt(float64(result.ConvertToInterval().End)))}
	} else {
		result = Interval{Number(math.Cbrt(float64(result.ConvertToInterval().Start))),
			Number(math.Cbrt(float64(result.ConvertToInterval().End)))}
	}

	return result.ConvertToInterval(), nil
}

func (a *Alternative) FSMetric(to *Alternative, v Variants) (Number, error) {
	if a.countOfCriteria != to.countOfCriteria {
		return 0, InvalidSize
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result := Number(0)
	wg.Add(a.countOfCriteria)
	for i := 0; i < a.countOfCriteria; i++ {
		go func(i int) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				if tmp, innerError := a.grade[i].DiffNumber(to.grade[i], v); innerError != nil {
					err = innerError
					cancel()
					return
				} else {
					mu.Lock()
					result += tmp
					mu.Unlock()
				}
			}
		}(i)
	}

	wg.Wait()

	if err != nil {
		return 0, err
	}

	if v == SqrtDistance {
		result = Number(math.Sqrt(float64(result)))
	} else {
		result = Number(math.Cbrt(float64(result)))
	}

	return result, nil
}

func (a *Alternative) Sum() Evaluated {
	var sum Evaluated

	if reflect.TypeOf(a.grade[0]) == reflect.TypeOf(NumbersMin) {
		sum = Number(0.0)
	} else if reflect.TypeOf(a.grade[0]) == reflect.TypeOf(Interval{}) {
		sum = Interval{0.0, 0.0}
	} else if reflect.TypeOf(a.grade[0]) == reflect.TypeOf(&T1FS{}) {
		if a.grade[0].ConvertToT1FS(Default).form == Triangle {
			sum = NewT1FS(0.0, 0.0, 0.0)
		} else {
			sum = NewT1FS(0.0, 0.0, 0.0, 0.0)
		}
	} else if reflect.TypeOf(a.grade[0]) == reflect.TypeOf(&AIFS{}) {
		if a.grade[0].ConvertToAIFS(Default).form == Triangle {
			sum = NewAIFS(0.0, 0.0, 0.0, 0.0)
		} else {
			sum = NewAIFS(0.0, 0.0, 0.0, 0.0, 0.0)
		}
	} else if reflect.TypeOf(a.grade[0]) == reflect.TypeOf(&IT2FS{}) {
		if a.grade[0].ConvertToIT2FS(Default).form == Triangle {
			sum = NewIT2FS([]Interval{{0.0, 0.0}, {0.0, 0.0}}, []Number{0.0})
		} else {
			sum = NewIT2FS([]Interval{{0.0, 0.0}, {0.0, 0.0}}, []Number{0.0, 0.0})
		}
	}

	for _, el := range a.grade {
		sum = sum.Sum(el)
	}
	return sum
}
