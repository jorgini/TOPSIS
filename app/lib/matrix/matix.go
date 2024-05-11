package matrix

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"sync"
	"webApp/lib/eval"
	v "webApp/lib/variables"
)

type Matrix struct {
	Data              []Alternative `json:"data"`
	CountAlternatives int           `json:"cnt_alt"`
	CountCriteria     int           `json:"cnt_crit"`
	Criteria          []Criterion   `json:"criteria"`
	CriteriaSet       bool          `json:"is_crit_set"`
	HighType          string        `json:"high_type"`
	FormFs            v.Variants    `json:"form_fs"`
}

func NewMatrix(x, y int) *Matrix {
	m := Matrix{
		Data:              make([]Alternative, x),
		CountAlternatives: x,
		CountCriteria:     y,
		Criteria:          NewCriteria(y),
		HighType:          eval.NumbersMin.GetType(),
		FormFs:            v.None,
	}

	for i := range m.Data {
		m.Data[i] = newAlternative(y)
	}
	return &m
}

func (m *Matrix) Value() (driver.Value, error) {
	data, err := json.Marshal(m)
	return string(data), err
}

func (m *Matrix) Scan(src interface{}) error {
	var matrix Matrix
	var err error
	switch src.(type) {
	case string:
		err = json.Unmarshal([]byte(src.(string)), &matrix)
	case []byte:
		err = json.Unmarshal(src.([]byte), &matrix)
	default:
		return errors.New("incompatible type for Matrix")
	}
	if err != nil {
		return err
	}
	*m = matrix
	return nil
}

func (m *Matrix) UpdateAlternativeRatings(i int, ratings []eval.Rating) error {
	if i < 0 || i > m.CountAlternatives || len(ratings) != m.CountCriteria {
		return v.InvalidSize
	}

	for j := 0; j < m.CountCriteria; j++ {
		_ = m.SetValue(ratings[j], i, j)
	}
	return nil
}

func (m *Matrix) SetValue(value eval.Evaluated, i, j int) error {
	if i < m.CountAlternatives && j < m.CountCriteria {
		m.Data[i].Grade[j] = value.CopyEval()
		m.HighType = eval.HighType(m.HighType, value.GetType())
		if m.FormFs < value.GetForm() {
			m.FormFs = value.GetForm()
		}
	} else {
		return v.OutOfBounds
	}
	return nil
}

func (m *Matrix) setRatings(data [][]eval.Evaluated) error {
	if m.CountAlternatives == 0 {
		return nil
	}

	if m.CountAlternatives != len(data) || m.CountCriteria != len(data[0]) {
		return v.InvalidSize
	}

	for i := 0; i < m.CountAlternatives; i++ {
		for j := 0; j < m.CountCriteria; j++ {
			m.Data[i].Grade[j] = data[i][j].CopyEval()
			m.HighType = eval.HighType(m.HighType, data[i][j].GetType())
			if m.FormFs < data[i][j].GetForm() {
				m.FormFs = data[i][j].GetForm()
			}
		}
	}
	return nil
}

func (m *Matrix) GetAlternativeRatings(i int) ([]eval.Rating, error) {
	if i < 0 || i > m.CountAlternatives {
		return nil, v.InvalidSize
	}

	return m.Data[i].Grade, nil
}

func (m *Matrix) SetCriteria(criteria []Criterion) error {
	if m.CountCriteria != len(criteria) {
		return v.InvalidSize
	}

	for i := 0; i < m.CountCriteria; i++ {
		m.Criteria[i].set(criteria[i].Weight, criteria[i].TypeOfCriteria)
	}
	m.CriteriaSet = true
	return nil
}

func (m *Matrix) SetCriterion(Weight eval.Evaluated, typeOf bool, i int) error {
	if i < m.CountCriteria {
		m.Criteria[i].set(Weight, typeOf)
	} else {
		return v.OutOfBounds
	}
	m.CriteriaSet = true
	return nil
}

func (m *Matrix) castToType(t string, f v.Variants) {
	for i := range m.Data {
		for c := range m.Data[i].Grade {
			if t == (&eval.IT2FS{}).GetType() {
				m.Data[i].Grade[c].Evaluated = m.Data[i].Grade[c].ConvertToIT2FS(f)
			} else if t == (&eval.AIFS{}).GetType() {
				m.Data[i].Grade[c].Evaluated = m.Data[i].Grade[c].ConvertToAIFS(f)
			} else if t == (&eval.T1FS{}).GetType() {
				m.Data[i].Grade[c].Evaluated = m.Data[i].Grade[c].ConvertToT1FS(f)
			} else if t == (eval.Interval{}).GetType() {
				m.Data[i].Grade[c].Evaluated = m.Data[i].Grade[c].ConvertToInterval()
			}
		}
	}
}

func TypingMatrices(g int, matrices ...Matrix) error {
	x, y := matrices[0].CountAlternatives, matrices[0].CountCriteria
	var highestType string
	var highestForm v.Variants

	var wg sync.WaitGroup

	for k := range matrices {
		if matrices[k].CountAlternatives != x || matrices[k].CountCriteria != y {
			return v.InvalidSize
		}

		highestType = eval.HighType(highestType, matrices[k].HighType)
		if highestForm < matrices[k].FormFs {
			highestForm = matrices[k].FormFs
		}
	}

	if g > len(matrices) {
		g = len(matrices)
	}
	off := len(matrices) / g

	wg.Add(g)
	for b := 0; b < g; b++ {
		go func(b int) {
			defer wg.Done()

			start := b * off
			end := (b + 1) * off
			if b == g-1 {
				end = len(matrices)
			}

			for k := start; k < end; k++ {
				matrices[k].castToType(highestType, highestForm)
			}
		}(b)
	}
	wg.Wait()

	weightType := eval.NumbersMin.GetType()

	for k := range matrices {
		weightType = eval.HighType(weightType, GetHighType(matrices[k].Criteria))
	}

	if weightType == (eval.Interval{}).GetType() {
		wg.Add(g)
		for b := 0; b < g; b++ {
			go func(b int) {
				defer wg.Done()

				start := b * off
				end := (b + 1) * off
				if b == g-1 {
					end = len(matrices)
				}

				for k := start; k < end; k++ {
					for i := range matrices[k].Criteria {
						matrices[k].Criteria[i].Weight.Evaluated = matrices[k].Criteria[i].Weight.ConvertToInterval()
					}
				}
			}(b)
		}
		wg.Wait()
	}

	return nil
}

func AggregateRatings(matrices []Matrix, weights []eval.Evaluated, g int) (*Matrix, error) {
	result := NewMatrix(matrices[0].CountAlternatives, matrices[0].CountCriteria)
	if err := TypingMatrices(g, matrices...); err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	if g > len(result.Data) {
		g = len(result.Data)
	}
	off := len(result.Data) / g

	wg.Add(g)
	for b := 0; b < g; b++ {
		go func(b int) {
			defer wg.Done()

			start := b * off
			end := (b + 1) * off
			if b == g-1 {
				end = len(result.Data)
			}

			for i := start; i < end; i++ {
				for j := range result.Data[i].Grade {
					for k := range matrices {
						if result.Data[i].Grade[j].IsNil() {
							_ = result.SetValue(matrices[k].Data[i].Grade[j].Weighted(weights[k]), i, j)
						} else {
							_ = result.SetValue(result.Data[i].Grade[j].Sum(matrices[k].Data[i].Grade[j].Weighted(weights[k])),
								i, j)
						}
					}
				}
			}
		}(b)
	}
	wg.Wait()

	for i := range result.Criteria {
		result.Criteria[i].set(matrices[0].Criteria[i].Weight, matrices[0].Criteria[i].TypeOfCriteria)
	}
	return result, nil
}

func CopyMatrix(matrix *Matrix) *Matrix {
	newMatrix := NewMatrix(matrix.CountAlternatives, matrix.CountCriteria)

	for i := range newMatrix.Data {
		_ = newMatrix.UpdateAlternativeRatings(i, matrix.Data[i].Grade)
	}

	for i := range newMatrix.Criteria {
		newMatrix.Criteria[i] = CopyCriterion(matrix.Criteria[i])
	}
	return newMatrix
}
