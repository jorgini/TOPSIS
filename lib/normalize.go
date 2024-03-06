package lib

import (
	"context"
	"math"
	"reflect"
	"sync"
)

func getBestValueWithCond(a, b Evaluated, typeOfCriterion bool) Number {
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

func getNormValueByMax(e Evaluated, min, max Number, typeOfCriterion bool) Evaluated {
	if reflect.TypeOf(e) == reflect.TypeOf(&T1FS{}) && typeOfCriterion == Cost {
		vertices := make([]Number, len(e.ConvertToT1FS(Default).vert))

		for k := range vertices {
			vertices[k] = min /
				e.ConvertToT1FS(Default).vert[len(vertices)-k-1]
		}
		e.ConvertToT1FS(Default).vert = vertices
	} else if reflect.TypeOf(e) == reflect.TypeOf(&AIFS{}) &&
		typeOfCriterion == Cost {
		vertices := make([]Number, len(e.ConvertToAIFS(Default).vert))

		for k := range vertices {
			vertices[k] = min /
				e.ConvertToAIFS(Default).vert[len(vertices)-k-1]
		}
		e.ConvertToAIFS(Default).vert = vertices
	} else if reflect.TypeOf(e) == reflect.TypeOf(&IT2FS{}) &&
		typeOfCriterion == Cost {
		grade := e.ConvertToIT2FS(Default)
		bottom := []Interval{{min / grade.bottom[1].End,
			min / grade.bottom[1].Start}, {min / grade.bottom[0].End,
			min / grade.bottom[0].Start}}

		var upward []Number
		if grade.form == Triangle {
			upward = []Number{min / grade.upward[0]}
		} else {
			upward = []Number{min / grade.upward[1], min / grade.upward[0]}
		}

		grade.bottom = bottom
		grade.upward = upward
	} else {
		e = e.Weighted(1 / max)
	}
	return e
}

func (m *Matrix) getMinMaxRecord(j int) (Number, Number) {
	maximum := NumbersMin
	minimum := NumbersMax

	for i := range m.data {
		eval := m.data[i].grade[j]

		if reflect.TypeOf(eval) == reflect.TypeOf(&T1FS{}) ||
			reflect.TypeOf(eval) == reflect.TypeOf(&AIFS{}) ||
			reflect.TypeOf(eval) == reflect.TypeOf(&IT2FS{}) {
			if m.criteria[j].typeOfCriteria == Benefit {
				maximum = getBestValueWithCond(maximum, eval, m.criteria[j].typeOfCriteria)
			} else {
				minimum = getBestValueWithCond(minimum, eval, m.criteria[j].typeOfCriteria)
			}
		} else {
			maximum = getBestValueWithCond(maximum, eval, m.criteria[j].typeOfCriteria)
		}
	}
	return minimum, maximum
}

func (m *Matrix) getSumForCriterion(j int) Number {
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
	return sum
}

func (m *Matrix) normalizationValue(v Variants, cancel context.CancelFunc) error {
	var wg sync.WaitGroup
	var err error = nil

	if v == NormalizeWithSum {
		wg.Add(m.countCriteria)
		for j := range m.criteria {
			go func(j int) {
				defer wg.Done()

				sum := m.getSumForCriterion(j)

				if sum == 0.0 {
					err = EmptyValues
					cancel()
					return
				}

				for i := range m.data {
					m.data[i].grade[j] = m.data[i].grade[j].Weighted(1 / sum)
				}
			}(j)
		}
		wg.Wait()
	} else if v == NormalizeWithMax {
		wg.Add(m.countCriteria)
		for j, c := range m.criteria {
			go func(j int, c Criterion) {
				defer wg.Done()
				minimum, maximum := m.getMinMaxRecord(j)

				if maximum.ConvertToNumber() == 0.0 {
					err = EmptyValues
					cancel()
					return
				}

				var wg2 sync.WaitGroup
				wg2.Add(m.countAlternatives)
				for i := range m.data {
					go func(i int) {
						defer wg2.Done()
						m.data[i].grade[j] = getNormValueByMax(m.data[i].grade[j], minimum, maximum, c.typeOfCriteria)
					}(i)
				}
				wg2.Wait()
			}(j, c)
		}

		wg.Wait()
	} else {
		err = InvalidCaseOfOperation
	}
	return err
}

func (m *Matrix) normalizationWeights(w Variants, cancel context.CancelFunc) error {
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
		cancel()
		return EmptyValues
	}

	var wg sync.WaitGroup
	wg.Add(m.countCriteria)
	for j, c := range m.criteria {
		go func(j int, c Criterion) {
			defer wg.Done()
			if reflect.TypeOf(c.weight) == reflect.TypeOf(Interval{}) && w != NormalizeWeightsByMidPoint {
				m.criteria[j].weight = Interval{c.weight.ConvertToInterval().Start / highSum,
					c.weight.ConvertToInterval().End / lowerSum}
			} else {
				m.criteria[j].weight = c.weight.Weighted(1 / highSum)
			}
		}(j, c)
	}

	wg.Wait()
	return nil
}

func (m *Matrix) Normalization(v Variants, w Variants) error {
	var wg sync.WaitGroup
	var err error = nil
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg.Add(2)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			return
		default:
			err = m.normalizationValue(v, cancel)
		}
	}()

	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			return
		default:
			err = m.normalizationWeights(w, cancel)
		}
	}()

	wg.Wait()

	return err
}

func (m *Matrix) CalcWeightedMatrix() {
	var wg sync.WaitGroup
	wg.Add(m.countCriteria)
	for j, c := range m.criteria {
		go func(j int, c Criterion) {
			defer wg.Done()

			for i := 0; i < m.countAlternatives; i++ {
				m.data[i].grade[j] = m.data[i].grade[j].Weighted(c.weight)
			}
		}(j, c)
	}
	wg.Wait()
}
