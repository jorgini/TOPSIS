package matrix

import (
	"errors"
	"reflect"
)

type Variants int

func (m *Matrix) normalizationValue(v Variants) error {
	if reflect.TypeOf(m.data[0].grade[0]) == reflect.TypeOf(NumbersMax) {
		for i := range m.data {
			newdata := make([]Numbers, m.countAlternatives)
			for j := range m.data[i].grade {
				newdata[j] = m.data[i].grade[j].ConvertToNumbers()
			}

			if err := NormalizeNumbers(newdata, v); err != nil {
				return errors.Join(err)
			}

			for j := range m.data[i].grade {
				m.data[i].grade[j] = newdata[j]
			}
		}
	} else if reflect.TypeOf(m.data[0].grade[0]) == reflect.TypeOf(Interval{}) {
		for i := range m.data {
			newdata := make([]Interval, m.countAlternatives)
			for j := range m.data[i].grade {
				newdata[j] = m.data[i].grade[j].ConvertToInterval()
			}
			if err := NormalizeIntervals(newdata, v); err != nil {
				return errors.Join(err)
			}

			for j := range m.data[i].grade {
				m.data[i].grade[j] = newdata[j]
			}
		}
	}
	return nil
}

func (m *Matrix) normalizationWeights(w Variants) error {
	if reflect.TypeOf(m.criteria[0].weight) == reflect.TypeOf(NumbersMax) {
		sum := Numbers(0.0)
		for j := range m.criteria {
			sum += m.criteria[j].weight.ConvertToNumbers()
		}

		if sum == Numbers(0.0) {
			return errors.New("can't normalize weights")
		}

		for j := range m.criteria {
			m.criteria[j].weight = m.criteria[j].weight.ConvertToNumbers() / sum
		}
	} else if reflect.TypeOf(m.criteria[0].weight) == reflect.TypeOf(Interval{}) {
		// придумать
	}
	return nil
}

func (m *Matrix) Normalization(v Variants, w Variants) error {
	if err := m.normalizationValue(v); err != nil {
		return errors.Join(err)
	}

	if err := m.normalizationWeights(w); err != nil {
		return errors.Join(err)
	}

	return nil
}
