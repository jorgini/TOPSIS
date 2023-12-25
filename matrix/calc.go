package matrix

import "errors"

func (m *Matrix) CalcWeightedMatrix() {
	for j := 0; j < m.countCriteria; j++ {
		for i := 0; i < m.countAlternatives; i++ {
			m.data[i].grade[j] = m.data[i].grade[j].Weighted(m.criteria[j].weight)
		}
	}
}

func (m *Matrix) FindIdeals() {
	positive := Alternative{make([]Evaluated, m.countCriteria), m.countCriteria}
	negative := Alternative{make([]Evaluated, m.countCriteria), m.countCriteria}

	for i := 0; i < m.countAlternatives; i++ {
		for j := 0; j < m.countCriteria; j++ {
			if positive.grade[j] == nil {
				positive.grade[j] = m.data[i].grade[j]
			} else {
				positive.grade[j] = max(positive.grade[j], m.data[i].grade[j], m.criteria[j].typeOfCriteria)
			}

			if negative.grade[j] == nil {
				negative.grade[j] = m.data[i].grade[j]
			} else {
				negative.grade[j] = min(negative.grade[j], m.data[i].grade[j], m.criteria[j].typeOfCriteria)
			}
		}
	}

	m.positiveIdeal = &positive
	m.negativeIdeal = &negative
	m.idealsFind = true
}

func (m *Matrix) FindDistanceToIdeals() error {
	for i := 0; i < m.countAlternatives; i++ {
		var err error
		if m.distancesToPositive[i], err = m.data[i].FindDistance(m.positiveIdeal); err != nil {
			return errors.Join(err)
		}

		if m.distancesToNegative[i], err = m.data[i].FindDistance(m.negativeIdeal); err != nil {
			return errors.Join(err)
		}
	}
	m.distancesFind = true
	return nil
}

func (m *Matrix) CalcCloseness() error {
	for i := 0; i < m.countAlternatives; i++ {
		if m.distancesToNegative[i]+m.distancesToPositive[i] == Numbers(0) {
			return errors.New("can't calc relative closeness")
		}

		m.relativeCloseness[i] = m.distancesToNegative[i] /
			(m.distancesToNegative[i] + m.distancesToPositive[i])
	}
	m.closenessFind = true
	return nil
}
