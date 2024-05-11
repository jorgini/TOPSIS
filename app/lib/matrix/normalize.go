package matrix

import (
	"context"
	"math"
	"sync"
	"webApp/lib/eval"
	v "webApp/lib/variables"
)

func getBestValueWithCond(a, b eval.Evaluated, typeOfCriterion bool) eval.Number {
	val := []eval.Evaluated{a, b}
	for i := range val {
		if val[i].GetType() == (eval.Interval{}).GetType() {
			if typeOfCriterion == v.Benefit {
				val[i] = val[i].ConvertToInterval().End
			} else {
				val[i] = val[i].ConvertToInterval().Start
			}
		}

		if val[i].GetType() == (&eval.T1FS{}).GetType() {
			if typeOfCriterion == v.Benefit {
				val[i] = val[i].ConvertToT1FS(v.Default).Vert[len(val[i].ConvertToT1FS(v.Default).Vert)-1]
			} else {
				val[i] = val[i].ConvertToT1FS(v.Default).Vert[0]
			}
		}

		if val[i].GetType() == (&eval.AIFS{}).GetType() {
			if typeOfCriterion == v.Benefit {
				val[i] = val[i].ConvertToAIFS(v.Default).Vert[len(val[i].ConvertToAIFS(v.Default).Vert)-1]
			} else {
				val[i] = val[i].ConvertToAIFS(v.Default).Vert[0]
			}
		}

		if val[i].GetType() == (&eval.IT2FS{}).GetType() {
			if typeOfCriterion == v.Benefit {
				val[i] = val[i].ConvertToIT2FS(v.Default).Bottom[1].End
			} else {
				val[i] = val[i].ConvertToIT2FS(v.Default).Bottom[0].Start
			}
		}
	}

	if typeOfCriterion == v.Benefit {
		return eval.Number(math.Max(float64(val[0].ConvertToNumber()), float64(val[1].ConvertToNumber())))
	} else {
		return eval.Number(math.Min(float64(val[0].ConvertToNumber()), float64(val[1].ConvertToNumber())))
	}
}

func getNormValueByMax(e eval.Evaluated, min, max eval.Number, typeOfCriterion bool) eval.Rating {
	if e.GetType() == (&eval.T1FS{}).GetType() && typeOfCriterion == v.Cost {
		vertices := make([]eval.Number, len(e.ConvertToT1FS(v.Default).Vert))

		for k := range vertices {
			vertices[k] = min /
				e.ConvertToT1FS(v.Default).Vert[len(vertices)-k-1]
		}
		e.ConvertToT1FS(v.Default).Vert = vertices
	} else if e.GetType() == (&eval.AIFS{}).GetType() &&
		typeOfCriterion == v.Cost {
		vertices := make([]eval.Number, len(e.ConvertToAIFS(v.Default).Vert))

		for k := range vertices {
			vertices[k] = min /
				e.ConvertToAIFS(v.Default).Vert[len(vertices)-k-1]
		}
		e.ConvertToAIFS(v.Default).Vert = vertices
	} else if e.GetType() == (&eval.IT2FS{}).GetType() &&
		typeOfCriterion == v.Cost {
		grade := e.ConvertToIT2FS(v.Default)
		Bottom := []eval.Interval{{min / grade.Bottom[1].End,
			min / grade.Bottom[1].Start}, {min / grade.Bottom[0].End,
			min / grade.Bottom[0].Start}}

		var Upward []eval.Number
		if grade.Form == v.Triangle {
			Upward = []eval.Number{min / grade.Upward[0]}
		} else {
			Upward = []eval.Number{min / grade.Upward[1], min / grade.Upward[0]}
		}

		grade.Bottom = Bottom
		grade.Upward = Upward
	} else {
		e = e.Weighted(1 / max)
	}
	return eval.Rating{Evaluated: e}
}

func (m *Matrix) getMinMaxRecord(j int) (eval.Number, eval.Number) {
	maximum := eval.NumbersMin
	minimum := eval.NumbersMax

	for i := range m.Data {
		rating := m.Data[i].Grade[j]

		if rating.GetType() == (&eval.T1FS{}).GetType() || rating.GetType() == (&eval.AIFS{}).GetType() ||
			rating.GetType() == (&eval.IT2FS{}).GetType() {
			if m.Criteria[j].TypeOfCriteria == v.Benefit {
				maximum = getBestValueWithCond(maximum, rating, m.Criteria[j].TypeOfCriteria)
			} else {
				minimum = getBestValueWithCond(minimum, rating, m.Criteria[j].TypeOfCriteria)
			}
		} else {
			maximum = getBestValueWithCond(maximum, rating, m.Criteria[j].TypeOfCriteria)
		}
	}
	return minimum, maximum
}

func (m *Matrix) getSumForCriterion(j int) eval.Number {
	sum := eval.Number(0.0)
	for i := range m.Data {
		rating := m.Data[i].Grade[j]

		if rating.GetType() == eval.NumbersMax.GetType() {
			sum += rating.ConvertToNumber() * rating.ConvertToNumber()
		} else if rating.GetType() == (eval.Interval{}).GetType() {
			sum += rating.ConvertToInterval().Start*rating.ConvertToInterval().Start +
				rating.ConvertToInterval().End*rating.ConvertToInterval().End
		} else if rating.GetType() == (&eval.T1FS{}).GetType() {
			tmp := rating.ConvertToT1FS(v.Default)
			sum += tmp.Vert[0]*tmp.Vert[0] + tmp.Vert[len(tmp.Vert)-1]*tmp.Vert[len(tmp.Vert)-1]
		} else if rating.GetType() == (&eval.AIFS{}).GetType() {
			tmp := rating.ConvertToAIFS(v.Default)
			sum += tmp.Vert[0]*tmp.Vert[0] + tmp.Vert[len(tmp.Vert)-1]*tmp.Vert[len(tmp.Vert)-1]
		} else if rating.GetType() == (&eval.IT2FS{}).GetType() {
			tmp := rating.ConvertToIT2FS(v.Default)
			sum += tmp.Bottom[0].Start*tmp.Bottom[0].Start + tmp.Bottom[1].End*tmp.Bottom[1].End
		}
	}

	sum = eval.Number(math.Sqrt(float64(sum)))
	return sum
}

func (m *Matrix) normalizationValue(variants v.Variants, g int) error {
	var wg sync.WaitGroup
	var err error = nil
	if g > m.CountCriteria {
		g = m.CountCriteria
	}
	off := m.CountCriteria / g

	if variants == v.NormalizeWithSum {
		wg.Add(g)
		for j := 0; j < g; j++ {
			go func(j int) {
				defer wg.Done()

				start := j * off
				end := (j + 1) * off
				if j == g-1 {
					end = m.CountCriteria
				}

				for c := start; c < end; c++ {
					sum := m.getSumForCriterion(c)

					if sum == 0.0 {
						err = v.EmptyValues
						return
					}

					for i := range m.Data {
						m.Data[i].Grade[c] = m.Data[i].Grade[c].Weighted(1 / sum)
					}
				}
			}(j)
		}
		wg.Wait()
	} else if variants == v.NormalizeValueWithMax {
		wg.Add(g)
		for j := 0; j < g; j++ {
			go func(j int) {
				defer wg.Done()

				start := j * off
				end := (j + 1) * off
				if j == g-1 {
					end = m.CountCriteria
				}

				for c := start; c < end; c++ {
					criterion := m.Criteria[c]
					minimum, maximum := m.getMinMaxRecord(c)

					if maximum.ConvertToNumber() == 0.0 {
						err = v.EmptyValues
						return
					}

					for i := range m.Data {
						m.Data[i].Grade[c] = getNormValueByMax(m.Data[i].Grade[c], minimum, maximum, criterion.TypeOfCriteria)
					}
				}
			}(j)
		}

		wg.Wait()
	} else {
		err = v.InvalidCaseOfOperation
	}
	return err
}

func (m *Matrix) normalizationWeights(variants v.Variants) error {
	highSum := eval.Number(0.0)
	lowerSum := eval.Number(0.0)

	for _, c := range m.Criteria {
		if c.Weight.GetType() == (eval.Interval{}).GetType() && variants != v.NormalizeWeightsByMidPoint {
			highSum += c.Weight.ConvertToInterval().End
			lowerSum += c.Weight.ConvertToInterval().Start
		} else {
			highSum += c.Weight.ConvertToNumber()
		}
	}

	if highSum == 0 {
		return v.EmptyValues
	}

	for j, c := range m.Criteria {
		if c.Weight.GetType() == (eval.Interval{}).GetType() && variants != v.NormalizeWeightsByMidPoint {
			m.Criteria[j].Weight.Evaluated = eval.Interval{Start: c.Weight.ConvertToInterval().Start / highSum,
				End: c.Weight.ConvertToInterval().End / lowerSum}
		} else {
			m.Criteria[j].Weight = c.Weight.Weighted(1 / highSum)
		}
	}

	return nil
}

func (m *Matrix) Normalization(values v.Variants, weights v.Variants, g int) error {
	var wg sync.WaitGroup
	var err error = nil
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if g > 1 {
		g--
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			return
		default:
			if inerr := m.normalizationValue(values, g); inerr != nil {
				cancel()
				err = inerr
			}
		}
	}()

	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			return
		default:
			if inerr := m.normalizationWeights(weights); inerr != nil {
				cancel()
				err = inerr
			}
		}
	}()

	wg.Wait()
	return err
}

func (m *Matrix) CalcWeightedMatrix(g int) {
	var wg sync.WaitGroup
	if g > m.CountCriteria {
		g = m.CountCriteria
	}
	off := m.CountCriteria / g
	wg.Add(g)
	for j := 0; j < g; j++ {
		go func(j int) {
			defer wg.Done()

			start := j * off
			end := (j + 1) * off
			if j == g-1 {
				end = m.CountCriteria
			}

			for c := start; c < end; c++ {
				criterion := m.Criteria[c]
				for i := 0; i < m.CountAlternatives; i++ {
					m.Data[i].Grade[c] = m.Data[i].Grade[c].Weighted(criterion.Weight)
				}
			}
		}(j)
	}
	wg.Wait()
}
