package matrix

import (
	"errors"
	"math"
	"reflect"
)

func (m *Matrix) normalizationValue(v Variants) error {
	if v == NormalizeWithSum {
		for j := range m.criteria {
			sum := Number(0)
			for i := range m.data {
				eval := m.data[i].grade[j]

				if reflect.TypeOf(eval) == reflect.TypeOf(NumbersMax) {
					sum += eval.ConvertToNumbers() * eval.ConvertToNumbers()
				} else if reflect.TypeOf(eval) == reflect.TypeOf(Interval{}) {
					sum += eval.ConvertToInterval().Start*eval.ConvertToInterval().Start +
						eval.ConvertToInterval().End*eval.ConvertToInterval().End
				} else if reflect.TypeOf(eval) == reflect.TypeOf(&T1FS{}) {
					tmp := eval.ConvertToT1FS(Default)
					sum += tmp.vert[0]*tmp.vert[0] + tmp.vert[len(tmp.vert)-1]*tmp.vert[len(tmp.vert)-1]
				}
			}

			sum = Number(math.Sqrt(float64(sum)))

			if sum == 0.0 {
				return errors.New("empty values in alternative")
			}

			for i := range m.data {
				m.data[i].grade[j] = m.data[i].grade[j].Weighted(1 / sum)
			}
		}
	} else if v == NormalizeWithMax {
		for j, criterion := range m.criteria {
			maximum := NumbersMin
			minimum := NumbersMax

			for i := range m.data {
				eval := m.data[i].grade[j]

				if reflect.TypeOf(eval) == reflect.TypeOf(&T1FS{}) {
					if criterion.typeOfCriteria == Benefit {
						maximum = positiveIdealRateNumber(maximum, eval, criterion.typeOfCriteria)
					} else {
						minimum = positiveIdealRateNumber(minimum, eval, criterion.typeOfCriteria)
					}
				} else {
					maximum = positiveIdealRateNumber(maximum, eval, Benefit)
				}
			}

			if maximum.ConvertToNumbers() == 0.0 {
				return errors.New("empty values in alternative")
			}

			for i := range m.data {
				if reflect.TypeOf(m.data[i].grade[j]) == reflect.TypeOf(&T1FS{}) && criterion.typeOfCriteria == Cost {
					vertices := m.data[i].grade[j].ConvertToT1FS(Default).vert

					for j := range vertices {
						vertices[j] = minimum / vertices[len(vertices)-j-1]
					}
				} else {
					m.data[i].grade[j] = m.data[i].grade[j].Weighted(1 / maximum)
				}
			}
		}
	} else {
		return errors.New("incomplete case of normalization")
	}
	return nil
}

func (m *Matrix) normalizationWeights(w Variants) error {
	highSum := Number(0.0)
	lowerSum := Number(0.0)
	for _, c := range m.criteria {
		highSum += c.weight.ConvertToNumbers()
		if reflect.TypeOf(c.weight) == reflect.TypeOf(Interval{}) && w != NormalizeWeightsByMidPoint {
			highSum += c.weight.ConvertToInterval().End
			lowerSum += c.weight.ConvertToInterval().Start
		}
	}

	if highSum == 0 {
		return errors.New("can't normalize weights")
	}

	for j, c := range m.criteria {
		if reflect.TypeOf(c.weight) == reflect.TypeOf(Interval{}) && w != NormalizeWeightsByMidPoint {
			m.criteria[j].weight = Interval{c.weight.ConvertToInterval().Start / highSum,
				c.weight.ConvertToInterval().End / lowerSum}
		} else {
			m.criteria[j].weight = c.weight.Weighted(1 / highSum)
		}
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
