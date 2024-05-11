package topsis

import (
	"context"
	"sync"
	"webApp/lib/eval"
	"webApp/lib/matrix"
	v "webApp/lib/variables"
)

func positiveIdeal(alts []matrix.Alternative, Criteria []matrix.Criterion, t string, f v.Variants, g int) (matrix.Alternative, error) {
	if g == 1 {
		if t == eval.Number(0).GetType() {
			return positiveIdealRateNumber0(alts, Criteria)
		} else if t == (eval.Interval{}).GetType() {
			return positiveIdealRateInterval0(alts, Criteria)
		} else if t == (&eval.T1FS{}).GetType() {
			return positiveIdealRateT1FS0(alts, Criteria, f)
		} else if t == (&eval.AIFS{}).GetType() {
			return positiveIdealRateAIFS0(alts, Criteria)
		} else {
			return matrix.Alternative{}, v.IncompatibleTypes
		}
	}

	positive := matrix.Alternative{Grade: make([]eval.Rating, len(alts[0].Grade)), CountOfCriteria: len(Criteria)}

	var wg sync.WaitGroup
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	batches := g
	if batches > len(Criteria) {
		batches = len(Criteria)
	}
	off := len(Criteria) / batches

	wg.Add(batches)
	for b := 0; b < batches; b++ {
		go func(b int) {
			defer wg.Done()

			start := off * b
			end := start + off
			if b == batches-1 {
				end = len(Criteria)
			}

			for i := start; i < end; i++ {
				select {
				case <-ctx.Done():
					return
				default:
					c := Criteria[i]
					var tmpErr error
					if t == eval.Number(0).GetType() {
						positive.Grade[i], tmpErr = positiveIdealRateNumber(alts, c, i)
					} else if t == (eval.Interval{}).GetType() {
						positive.Grade[i], tmpErr = positiveIdealRateInterval(alts, c, i)
					} else if t == (&eval.T1FS{}).GetType() {
						positive.Grade[i], tmpErr = positiveIdealRateT1FS(alts, c, i, f)
					} else if t == (&eval.AIFS{}).GetType() {
						positive.Grade[i], tmpErr = positiveIdealRateAIFS(alts, c, i)
					} else {
						cancel()
						err = v.IncompatibleTypes
						return
					}

					if tmpErr != nil {
						cancel()
						err = tmpErr
						return
					}
				}
			}
		}(b)
	}

	wg.Wait()
	return positive, err
}

func negativeIdeal(alts []matrix.Alternative, Criteria []matrix.Criterion, t string, f v.Variants, g int) (matrix.Alternative, error) {
	return positiveIdeal(alts, matrix.ChangeTypes(Criteria), t, f, g)
}

func positiveIdealRateAIFS(alts []matrix.Alternative, c matrix.Criterion, i int) (eval.Rating, error) {
	positive := eval.Rating{}
	for j := range alts {
		if alts[j].Grade[i].GetType() != (&eval.AIFS{}).GetType() {
			return eval.Rating{}, v.IncompatibleTypes
		}

		altsGrade := alts[j].Grade[i].ConvertToAIFS(v.Default)

		if positive.IsNil() {
			positive.Evaluated = eval.NewAIFS(altsGrade.Pi, altsGrade.Vert...)
			continue
		}

		if c.TypeOfCriteria == v.Benefit {
			positive = eval.Max(positive, alts[j].Grade[i])
		} else {
			positive = eval.Min(positive, alts[j].Grade[i])
		}
	}

	return positive, nil
}

func positiveIdealRateT1FS(alts []matrix.Alternative, c matrix.Criterion, i int, Form v.Variants) (eval.Rating, error) {
	positive := eval.Rating{}

	for j := range alts {
		if alts[j].Grade[i].GetType() != (&eval.T1FS{}).GetType() && alts[j].Grade[i].GetType() != (&eval.IT2FS{}).GetType() {
			return eval.Rating{}, v.IncompatibleTypes
		}

		if positive.IsNil() && c.TypeOfCriteria == v.Benefit {
			if Form == v.Triangle {
				positive.Evaluated = eval.NewT1FS(eval.NumbersMin, eval.NumbersMin, eval.NumbersMin)
			} else if Form == v.Trapezoid {
				positive.Evaluated = eval.NewT1FS(eval.NumbersMin, eval.NumbersMin, eval.NumbersMin, eval.NumbersMin)
			} else {
				return eval.Rating{}, v.IncompatibleTypes
			}
		} else if positive.IsNil() && c.TypeOfCriteria == v.Cost {
			if Form == v.Triangle {
				positive.Evaluated = eval.NewT1FS(eval.NumbersMax, eval.NumbersMax, eval.NumbersMax)
			} else if Form == v.Trapezoid {
				positive.Evaluated = eval.NewT1FS(eval.NumbersMax, eval.NumbersMax, eval.NumbersMax, eval.NumbersMax)
			} else {
				return eval.Rating{}, v.IncompatibleTypes
			}
		}

		if alts[j].Grade[i].GetType() == (&eval.T1FS{}).GetType() {
			if c.TypeOfCriteria == v.Benefit {
				positive = eval.Max(positive, alts[j].Grade[i])
			} else {
				positive = eval.Min(positive, alts[j].Grade[i])
			}
		} else {
			posGrade := positive.ConvertToT1FS(v.Default)
			altsGrade := alts[j].Grade[i].ConvertToIT2FS(v.Default)

			if c.TypeOfCriteria == v.Benefit {
				posGrade.Vert[0] = eval.Max(posGrade.Vert[0], altsGrade.Bottom[0].End).ConvertToNumber()
				posGrade.Vert[1] = eval.Max(posGrade.Vert[1], altsGrade.Upward[0]).ConvertToNumber()

				if Form == v.Triangle {
					posGrade.Vert[2] = eval.Max(posGrade.Vert[2], altsGrade.Bottom[1].End).ConvertToNumber()
				} else {
					posGrade.Vert[2] = eval.Max(posGrade.Vert[2], altsGrade.Upward[1]).ConvertToNumber()
					posGrade.Vert[3] = eval.Max(posGrade.Vert[3], altsGrade.Bottom[1].End).ConvertToNumber()
				}
			} else {
				posGrade.Vert[0] = eval.Min(posGrade.Vert[0], altsGrade.Bottom[0].Start).ConvertToNumber()
				posGrade.Vert[1] = eval.Min(posGrade.Vert[1], altsGrade.Upward[0]).ConvertToNumber()

				if Form == v.Triangle {
					posGrade.Vert[2] = eval.Min(posGrade.Vert[2], altsGrade.Bottom[1].Start).ConvertToNumber()
				} else {
					posGrade.Vert[2] = eval.Min(posGrade.Vert[2], altsGrade.Upward[1]).ConvertToNumber()
					posGrade.Vert[3] = eval.Min(posGrade.Vert[3], altsGrade.Bottom[1].Start).ConvertToNumber()
				}
			}
		}
	}

	return positive, nil
}

func positiveIdealRateInterval(alts []matrix.Alternative, c matrix.Criterion, i int) (eval.Rating, error) {
	positive := eval.Rating{}
	for j := range alts {
		if alts[j].Grade[i].GetType() != (eval.Interval{}).GetType() {
			return eval.Rating{}, v.IncompatibleTypes
		}

		if positive.IsNil() {
			positive = alts[j].Grade[i]
			continue
		}

		if c.TypeOfCriteria == v.Benefit {
			positive = eval.Max(positive, alts[j].Grade[i])
		} else {
			positive = eval.Min(positive, alts[j].Grade[i])
		}
	}

	return positive, nil
}

func positiveIdealRateNumber(alts []matrix.Alternative, c matrix.Criterion, i int) (eval.Rating, error) {
	positive := eval.Rating{}
	for j := range alts {
		if alts[j].Grade[i].GetType() != eval.NumbersMin.GetType() && alts[j].Grade[i].GetType() != (eval.Interval{}).GetType() {
			return eval.Rating{}, v.IncompatibleTypes
		}

		if positive.IsNil() && c.TypeOfCriteria == v.Benefit {
			positive.Evaluated = eval.NumbersMin
		} else if positive.IsNil() && c.TypeOfCriteria == v.Cost {
			positive.Evaluated = eval.NumbersMax
		}

		if alts[j].Grade[i].GetType() == (eval.Interval{}).GetType() {
			if c.TypeOfCriteria == v.Benefit {
				positive = eval.Max(positive.ConvertToNumber(), alts[j].Grade[i].ConvertToInterval().End)
			} else {
				positive = eval.Min(positive.ConvertToNumber(), alts[j].Grade[i].ConvertToInterval().Start)
			}
		} else {
			if c.TypeOfCriteria == v.Benefit {
				positive = eval.Max(positive, alts[j].Grade[i])
			} else {
				positive = eval.Min(positive, alts[j].Grade[i])
			}
		}
	}

	return positive, nil
}
