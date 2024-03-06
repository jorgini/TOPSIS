package lib

import (
	"context"
	"errors"
	"reflect"
	"sync"
)

var (
	InvalidCaseOfOperation = errors.New("invalid case of operation")
	OutOfBounds            = errors.New("invalid index")
	InvalidSize            = errors.New("invalid size")
	EmptyValues            = errors.New("no values specified")
	IncompatibleTypes      = errors.New("incompatible types for operation")
	NoUsageMethod          = errors.New("this method shouldn't have usage")
)

type Matrix struct {
	data              []Alternative
	countAlternatives int
	countCriteria     int
	criteria          []Criterion
	criteriaSet       bool
	finalScores       []Evaluated
	highType          reflect.Type
	formFs            Variants
}

func NewMatrix(x, y int) *Matrix {
	m := Matrix{
		data:              make([]Alternative, x),
		countAlternatives: x,
		countCriteria:     y,
		criteria:          make([]Criterion, y),
		highType:          reflect.TypeOf(NumbersMin),
		formFs:            None,
	}

	for i := range m.data {
		m.data[i] = Alternative{make([]Evaluated, y), y}
	}
	return &m
}

func (m *Matrix) SetValue(value Evaluated, i, j int) error {
	if i < m.countAlternatives && j < m.countCriteria {
		m.data[i].grade[j] = value
		m.highType = HighType(m.highType, reflect.TypeOf(value))
		if m.formFs < value.GetForm() {
			m.formFs = value.GetForm()
		}
	} else {
		return OutOfBounds
	}
	return nil
}

func (m *Matrix) SetMatrix(data [][]Evaluated) error {
	if m.countAlternatives == 0 {
		return nil
	}

	if m.countAlternatives != len(data) || m.countCriteria != len(data[0]) {
		return InvalidSize
	}

	for i := 0; i < m.countAlternatives; i++ {
		for j := 0; j < m.countCriteria; j++ {
			m.data[i].grade[j] = data[i][j]
			m.highType = HighType(m.highType, reflect.TypeOf(data[i][j]))
			if m.formFs < data[i][j].GetForm() {
				m.formFs = data[i][j].GetForm()
			}
		}
	}
	return nil
}

func (m *Matrix) SetCriterion(weight Evaluated, typeOf bool, i int) error {
	if i < m.countCriteria {
		m.criteria[i].set(weight, typeOf)
	} else {
		return OutOfBounds
	}
	m.criteriaSet = true
	return nil
}

func (m *Matrix) SetCriteria(weights []Evaluated, typeOf []bool) error {
	if len(weights) != m.countCriteria || len(typeOf) != m.countCriteria {
		return InvalidSize
	} else {
		for i, w := range weights {
			m.criteria[i].set(w, typeOf[i])
		}
	}
	m.criteriaSet = true
	return nil
}

func (m *Matrix) SetTypeAndForm(t reflect.Type, f Variants) {
	m.formFs = f
	m.highType = t
}

func (m *Matrix) CastToType(t reflect.Type, f Variants) {
	var wg sync.WaitGroup
	wg.Add(len(m.data))
	for i := range m.data {
		go func(i int) {
			defer wg.Done()
			for j := range m.data[i].grade {
				if t == reflect.TypeOf(&IT2FS{}) {
					m.data[i].grade[j] = m.data[i].grade[j].ConvertToIT2FS(f)
				} else if t == reflect.TypeOf(&AIFS{}) {
					m.data[i].grade[j] = m.data[i].grade[j].ConvertToAIFS(f)
				} else if t == reflect.TypeOf(&T1FS{}) {
					m.data[i].grade[j] = m.data[i].grade[j].ConvertToT1FS(f)
				} else if t == reflect.TypeOf(Interval{}) {
					m.data[i].grade[j] = m.data[i].grade[j].ConvertToInterval()
				}
			}
		}(i)
	}
	wg.Wait()
}

func TypingMatrices(matrices ...*Matrix) error {
	x, y := matrices[0].countAlternatives, matrices[0].countCriteria
	var highestType reflect.Type
	var highestForm Variants
	var wg sync.WaitGroup
	var mu sync.Mutex
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(len(matrices))
	for k := range matrices {
		if matrices[k].countAlternatives != x || matrices[k].countCriteria != y {
			cancel()
			return InvalidSize
		}

		go func(k int, ctx context.Context) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				mu.Lock()
				highestType = HighType(highestType, matrices[k].highType)
				if highestForm < matrices[k].formFs {
					highestForm = matrices[k].formFs
				}
				mu.Unlock()
			}
		}(k, ctx)
	}

	wg.Wait()

	wg.Add(len(matrices))
	for k := range matrices {
		go func(k int) {
			defer wg.Done()
			matrices[k].CastToType(highestType, highestForm)
		}(k)
	}

	wg.Wait()

	weightType := reflect.TypeOf(NumbersMin)
	wg.Add(len(matrices))
	for k := range matrices {
		go func(k int) {
			defer wg.Done()
			mu.Lock()
			weightType = HighType(weightType, GetHighType(matrices[k].criteria))
			mu.Unlock()
		}(k)
	}

	wg.Wait()
	if weightType == reflect.TypeOf(Interval{}) {
		wg.Add(len(matrices))
		for k := range matrices {
			go func(k int) {
				defer wg.Done()
				for i := range matrices[k].criteria {
					matrices[k].criteria[i].weight = matrices[k].criteria[i].weight.ConvertToInterval()
				}
			}(k)
		}
	}

	wg.Wait()
	return nil
}

func AggregateRatings(matrices []*Matrix, weights []Evaluated) (*Matrix, error) {
	result := NewMatrix(matrices[0].countAlternatives, matrices[0].countCriteria)
	if err := TypingMatrices(matrices...); err != nil {
		return nil, err
	}
	result.highType = matrices[0].highType
	result.formFs = matrices[0].formFs

	var wg1 sync.WaitGroup
	var mu sync.Mutex

	wg1.Add(len(matrices))
	for k := range matrices {
		go func(k int) {
			defer wg1.Done()

			for i := range matrices[k].data {
				for j, eval := range matrices[k].data[i].grade {
					mu.Lock()
					if result.data[i].grade[j] == nil {
						result.data[i].grade[j] = eval.Weighted(weights[k])
					} else {
						result.data[i].grade[j] = result.data[i].grade[j].Sum(eval.Weighted(weights[k]))
					}
					mu.Unlock()
				}
			}
		}(k)
	}
	wg1.Wait()
	return result, nil
}
