package lib

import (
	"errors"
	"math"
	"reflect"
)

func GetBestValueWithCond(a, b Evaluated, typeOfCriterion bool) Number {
	val := []Evaluated{a, b}

	for i := range val {
		if reflect.TypeOf(val[i]) == reflect.TypeOf(Interval{}) {
			if typeOfCriterion == Benefit {
				val[i] = val[i].ConvertToInterval().End
			} else {
				val[i] = val[i].ConvertToInterval().Start
			}
		}

		if reflect.TypeOf(val[i]) == reflect.TypeOf(&T1FS{}) {
			if typeOfCriterion == Benefit {
				val[i] = val[i].ConvertToT1FS(Default).vert[len(val[i].ConvertToT1FS(Default).vert)-1]
			} else {
				val[i] = val[i].ConvertToT1FS(Default).vert[0]
			}
		}

		if reflect.TypeOf(val[i]) == reflect.TypeOf(&AIFS{}) {
			if typeOfCriterion == Benefit {
				val[i] = val[i].ConvertToAIFS(Default).vert[len(val[i].ConvertToAIFS(Default).vert)-1]
			} else {
				val[i] = val[i].ConvertToAIFS(Default).vert[0]
			}
		}

		if reflect.TypeOf(val[i]) == reflect.TypeOf(&IT2FS{}) {
			if typeOfCriterion == Benefit {
				val[i] = val[i].ConvertToIT2FS(Default).bottom[1].End
			} else {
				val[i] = val[i].ConvertToIT2FS(Default).bottom[0].Start
			}
		}
	}

	if typeOfCriterion == Benefit {
		return Number(math.Max(float64(val[0].ConvertToNumber()), float64(val[1].ConvertToNumber())))
	} else {
		return Number(math.Min(float64(val[0].ConvertToNumber()), float64(val[1].ConvertToNumber())))
	}
}

func (m *Matrix) normalizationValue(v Variants) error {
	if v == NormalizeWithSum {
		for j := range m.criteria {
			sum := Number(0.0)
			for i := range m.data {
				eval := m.data[i].grade[j]

				if reflect.TypeOf(eval) == reflect.TypeOf(NumbersMax) {
					sum += eval.ConvertToNumber() * eval.ConvertToNumber()
				} else if reflect.TypeOf(eval) == reflect.TypeOf(Interval{}) {
					sum += eval.ConvertToInterval().Start*eval.ConvertToInterval().Start +
						eval.ConvertToInterval().End*eval.ConvertToInterval().End
				} else if reflect.TypeOf(eval) == reflect.TypeOf(&T1FS{}) {
					tmp := eval.ConvertToT1FS(Default)
					sum += tmp.vert[0]*tmp.vert[0] + tmp.vert[len(tmp.vert)-1]*tmp.vert[len(tmp.vert)-1]
				} else if reflect.TypeOf(eval) == reflect.TypeOf(&AIFS{}) {
					tmp := eval.ConvertToAIFS(Default)
					sum += tmp.vert[0]*tmp.vert[0] + tmp.vert[len(tmp.vert)-1]*tmp.vert[len(tmp.vert)-1]
				} else if reflect.TypeOf(eval) == reflect.TypeOf(&IT2FS{}) {
					tmp := eval.ConvertToIT2FS(Default)
					sum += tmp.bottom[0].Start*tmp.bottom[0].Start + tmp.bottom[1].End*tmp.bottom[1].End
				}
			}

			sum = Number(math.Sqrt(float64(sum)))

			if sum == 0.0 {
				return EmptyValues
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

				if reflect.TypeOf(eval) == reflect.TypeOf(&T1FS{}) ||
					reflect.TypeOf(eval) == reflect.TypeOf(&AIFS{}) ||
					reflect.TypeOf(eval) == reflect.TypeOf(&IT2FS{}) {
					if criterion.typeOfCriteria == Benefit {
						maximum = GetBestValueWithCond(maximum, eval, criterion.typeOfCriteria)
					} else {
						minimum = GetBestValueWithCond(minimum, eval, criterion.typeOfCriteria)
					}
				} else {
					maximum = GetBestValueWithCond(maximum, eval, criterion.typeOfCriteria)
				}
			}

			if maximum.ConvertToNumber() == 0.0 {
				return EmptyValues
			}

			for i := range m.data {
				if reflect.TypeOf(m.data[i].grade[j]) == reflect.TypeOf(&T1FS{}) && criterion.typeOfCriteria == Cost {
					vertices := make([]Number, len(m.data[i].grade[j].ConvertToT1FS(Default).vert))

					for k := range vertices {
						vertices[k] = minimum / m.data[i].grade[j].ConvertToT1FS(Default).vert[len(vertices)-k-1]
					}
					m.data[i].grade[j].ConvertToT1FS(Default).vert = vertices
				} else if reflect.TypeOf(m.data[i].grade[j]) == reflect.TypeOf(&AIFS{}) && criterion.typeOfCriteria == Cost {
					vertices := make([]Number, len(m.data[i].grade[j].ConvertToAIFS(Default).vert))

					for k := range vertices {
						vertices[k] = minimum / m.data[i].grade[j].ConvertToAIFS(Default).vert[len(vertices)-k-1]
					}
					m.data[i].grade[j].ConvertToAIFS(Default).vert = vertices
				} else if reflect.TypeOf(m.data[i].grade[j]) == reflect.TypeOf(&IT2FS{}) && criterion.typeOfCriteria == Cost {
					grade := m.data[i].grade[j].ConvertToIT2FS(Default)
					bottom := []Interval{{minimum / grade.bottom[1].End, minimum / grade.bottom[1].Start},
						{minimum / grade.bottom[0].End, minimum / grade.bottom[0].Start}}

					var upward []Number
					if grade.form == Triangle {
						upward = []Number{minimum / grade.upward[0]}
					} else {
						upward = []Number{minimum / grade.upward[1], minimum / grade.upward[0]}
					}

					grade.bottom = bottom
					grade.upward = upward
				} else {
					m.data[i].grade[j] = m.data[i].grade[j].Weighted(1 / maximum)
				}
			}
		}
	} else {
		return InvalidCaseOfOperation
	}
	return nil
}

func (m *Matrix) normalizationWeights(w Variants) error {
	highSum := Number(0.0)
	lowerSum := Number(0.0)
	for _, c := range m.criteria {
		if reflect.TypeOf(c.weight) == reflect.TypeOf(Interval{}) && w != NormalizeWeightsByMidPoint {
			highSum += c.weight.ConvertToInterval().End
			lowerSum += c.weight.ConvertToInterval().Start
		} else {
			highSum += c.weight.ConvertToNumber()
		}
	}

	if highSum == 0 {
		return EmptyValues
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

func (m *Matrix) CalcWeightedMatrix() {
	for j := 0; j < m.countCriteria; j++ {
		for i := 0; i < m.countAlternatives; i++ {
			m.data[i].grade[j] = m.data[i].grade[j].Weighted(m.criteria[j].weight)
		}
	}
}
