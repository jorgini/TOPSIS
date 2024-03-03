package lib

import (
	"errors"
	"reflect"
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
	}

	for i := range m.data {
		m.data[i] = Alternative{make([]Evaluated, y), y}
	}
	return &m
}

func (m *Matrix) SetValue(value Evaluated, i, j int) error {
	if i < m.countAlternatives && j < m.countCriteria {
		m.data[i].grade[j] = value
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

func TypingMatrices(matrices ...*Matrix) (reflect.Type, Variants, error) {
	hasInterval := false
	hasT1FS := false
	hasAIFS := false
	hasIT2FS := false
	hasTriangle := false
	hasTrapezoid := false
	x, y := matrices[0].countAlternatives, matrices[0].countCriteria

	for k := range matrices {
		if matrices[k].countAlternatives != x || matrices[k].countCriteria != y {
			return nil, 0, InvalidSize
		}

		for i := range matrices[k].data {
			if hasIT2FS && hasTrapezoid {
				break
			}

			for _, eval := range matrices[k].data[i].grade {
				if reflect.TypeOf(eval) == reflect.TypeOf(Interval{}) {
					hasInterval = true
				}
				if reflect.TypeOf(eval) == reflect.TypeOf(&T1FS{}) {
					hasT1FS = true
					if eval.ConvertToT1FS(Default).form == Triangle {
						hasTriangle = true
					} else {
						hasTrapezoid = true
						break
					}
				}
				if reflect.TypeOf(eval) == reflect.TypeOf(&AIFS{}) {
					hasAIFS = true
					if eval.ConvertToAIFS(Default).form == Triangle {
						hasTriangle = true
					} else {
						hasTrapezoid = true
						break
					}
				}
				if reflect.TypeOf(eval) == reflect.TypeOf(&IT2FS{}) {
					hasIT2FS = true
					if eval.ConvertToIT2FS(Default).form == Triangle {
						hasTriangle = true
					} else {
						hasTrapezoid = true
						break
					}
				}
			}
		}
	}

	ret := reflect.TypeOf(NumbersMin)
	if hasIT2FS {
		ret = reflect.TypeOf(&IT2FS{})
	} else if hasAIFS {
		ret = reflect.TypeOf(&AIFS{})
	} else if hasT1FS {
		ret = reflect.TypeOf(&T1FS{})
	} else if hasInterval {
		ret = reflect.TypeOf(Interval{})
	}

	var form Variants
	if hasTrapezoid {
		form = Trapezoid
	} else if hasTriangle {
		form = Triangle
	}

	for k := range matrices {
		for i := range matrices[k].data {
			for j := range matrices[k].data[i].grade {
				if hasIT2FS && hasTrapezoid {
					matrices[k].data[i].grade[j] = matrices[k].data[i].grade[j].ConvertToIT2FS(Trapezoid)
				} else if hasIT2FS && hasTriangle {
					matrices[k].data[i].grade[j] = matrices[k].data[i].grade[j].ConvertToIT2FS(Triangle)
				} else if hasAIFS && hasTrapezoid {
					matrices[k].data[i].grade[j] = matrices[k].data[i].grade[j].ConvertToAIFS(Trapezoid)
				} else if hasAIFS && hasTriangle {
					matrices[k].data[i].grade[j] = matrices[k].data[i].grade[j].ConvertToAIFS(Triangle)
				} else if hasT1FS && hasTrapezoid {
					matrices[k].data[i].grade[j] = matrices[k].data[i].grade[j].ConvertToT1FS(Trapezoid)
				} else if hasT1FS && hasTriangle {
					matrices[k].data[i].grade[j] = matrices[k].data[i].grade[j].ConvertToT1FS(Triangle)
				} else if hasInterval {
					matrices[k].data[i].grade[j] = matrices[k].data[i].grade[j].ConvertToInterval()
				}
			}
		}
	}

	hasIntervalWeight := false
	for k := range matrices {
		for _, c := range matrices[k].criteria {
			if reflect.TypeOf(c.weight) == reflect.TypeOf(Interval{}) {
				hasIntervalWeight = true
				break
			}
		}
		if hasIntervalWeight {
			break
		}
	}

	if hasIntervalWeight {
		for k := range matrices {
			for i := range matrices[k].criteria {
				matrices[k].criteria[i].weight = matrices[k].criteria[i].weight.ConvertToInterval()
			}
		}
	}

	return ret, form, nil
}

func AggregateRatings(matrices []*Matrix, weights []Evaluated) (*Matrix, error) {
	result := NewMatrix(matrices[0].countAlternatives, matrices[0].countCriteria)
	if highType, form, err := TypingMatrices(matrices...); err != nil {
		return nil, errors.Join(err)
	} else {
		result.highType = highType
		result.formFs = form
	}

	for k := range matrices {
		for i := range matrices[k].data {
			for j, eval := range matrices[k].data[i].grade {
				if result.data[i].grade[j] == nil {
					result.data[i].grade[j] = eval.Weighted(weights[k])
				} else {
					result.data[i].grade[j] = result.data[i].grade[j].Sum(eval.Weighted(weights[k]))
				}
			}
		}
	}
	return result, nil
}
