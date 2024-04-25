package topsis

import (
	"context"
	"sync"
	"webApp/lib/eval"
	"webApp/lib/matrix"
	v "webApp/lib/variables"
)

func positiveIdealRateAIFS(alts []matrix.Alternative, Criteria []matrix.Criterion) (matrix.Alternative, error) {
	var wg sync.WaitGroup
	var err error = nil
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	positive := matrix.Alternative{Grade: make([]eval.Rating, len(alts[0].Grade)), CountOfCriteria: len(Criteria)}

	wg.Add(len(Criteria))
	for i, c := range Criteria {
		go func(i int, c matrix.Criterion) {
			defer wg.Done()

			for j := range alts {
				select {
				case <-ctx.Done():
					return
				default:
					if alts[j].Grade[i].GetType() != (&eval.AIFS{}).GetType() {
						cancel()
						err = v.IncompatibleTypes
						return
					}

					altsGrade := alts[j].Grade[i].ConvertToAIFS(v.Default)

					if positive.Grade[i].IsNil() {
						positive.Grade[i].Evaluated = eval.NewAIFS(altsGrade.Pi, altsGrade.Vert...)
					}

					if c.TypeOfCriteria == v.Benefit {
						positive.Grade[i] = eval.Max(positive.Grade[i], alts[j].Grade[i])
					} else {
						positive.Grade[i] = eval.Min(positive.Grade[i], alts[j].Grade[i])
					}
				}
			}
		}(i, c)
	}

	wg.Wait()
	return positive, err
}

func negativeIdealRateAIFS(alts []matrix.Alternative, Criteria []matrix.Criterion) (matrix.Alternative, error) {
	return positiveIdealRateAIFS(alts, matrix.ChangeTypes(Criteria))
}

func positiveIdealRateT1FS(alts []matrix.Alternative, Criteria []matrix.Criterion, Form v.Variants) (matrix.Alternative, error) {
	var wg sync.WaitGroup
	var err error = nil
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	positive := matrix.Alternative{Grade: make([]eval.Rating, len(alts[0].Grade)), CountOfCriteria: len(Criteria)}

	wg.Add(len(Criteria))
	for i, c := range Criteria {
		go func(i int, c matrix.Criterion) {
			defer wg.Done()

			for j := range alts {
				select {
				case <-ctx.Done():
					return
				default:
					if alts[j].Grade[i].GetType() != (&eval.T1FS{}).GetType() &&
						alts[j].Grade[i].GetType() != (&eval.IT2FS{}).GetType() {
						cancel()
						err = v.IncompatibleTypes
						return
					}

					if positive.Grade[i].IsNil() && c.TypeOfCriteria == v.Benefit {
						if Form == v.Triangle {
							positive.Grade[i].Evaluated = eval.NewT1FS(eval.NumbersMin, eval.NumbersMin, eval.NumbersMin)
						} else {
							positive.Grade[i].Evaluated = eval.NewT1FS(eval.NumbersMin, eval.NumbersMin, eval.NumbersMin, eval.NumbersMin)
						}
					} else if positive.Grade[i].IsNil() && c.TypeOfCriteria == v.Cost {
						if Form == v.Triangle {
							positive.Grade[i].Evaluated = eval.NewT1FS(eval.NumbersMax, eval.NumbersMax, eval.NumbersMax)
						} else {
							positive.Grade[i].Evaluated = eval.NewT1FS(eval.NumbersMax, eval.NumbersMax, eval.NumbersMax, eval.NumbersMax)
						}
					}

					if alts[j].Grade[i].GetType() == (&eval.T1FS{}).GetType() {
						if c.TypeOfCriteria == v.Benefit {
							positive.Grade[i] = eval.Max(positive.Grade[i], alts[j].Grade[i])
						} else {
							positive.Grade[i] = eval.Min(positive.Grade[i], alts[j].Grade[i])
						}
					} else {
						posGrade := positive.Grade[i].ConvertToT1FS(v.Default)
						altsGrade := alts[j].Grade[i].ConvertToIT2FS(v.Default)

						if c.TypeOfCriteria == v.Benefit {
							posGrade.Vert[0] =
								eval.Max(posGrade.Vert[0], altsGrade.Bottom[0].End).ConvertToNumber()
							posGrade.Vert[1] =
								eval.Max(posGrade.Vert[1], altsGrade.Upward[0]).ConvertToNumber()
							if Form == v.Triangle {
								posGrade.Vert[2] =
									eval.Max(posGrade.Vert[2], altsGrade.Bottom[1].End).ConvertToNumber()
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
			}
		}(i, c)
	}
	wg.Wait()
	return positive, err
}

func negativeIdealRateT1FS(alts []matrix.Alternative, Criteria []matrix.Criterion, Form v.Variants) (matrix.Alternative, error) {
	return positiveIdealRateT1FS(alts, matrix.ChangeTypes(Criteria), Form)
}

func positiveIdealRateInterval(alts []matrix.Alternative, Criteria []matrix.Criterion) (matrix.Alternative, error) {
	var wg sync.WaitGroup
	var err error = nil
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	positive := matrix.Alternative{Grade: make([]eval.Rating, len(alts[0].Grade)), CountOfCriteria: len(Criteria)}

	wg.Add(len(Criteria))
	for i, c := range Criteria {
		go func(i int, c matrix.Criterion) {
			defer wg.Done()

			for j := range alts {
				select {
				case <-ctx.Done():
					return
				default:
					if alts[j].Grade[i].GetType() != (eval.Interval{}).GetType() {
						cancel()
						err = v.IncompatibleTypes
						return
					}

					if positive.Grade[i].IsNil() {
						positive.Grade[i] = alts[j].Grade[i]
						continue
					}

					if c.TypeOfCriteria == v.Benefit {
						positive.Grade[i] = eval.Max(positive.Grade[i], alts[j].Grade[i])
					} else {
						positive.Grade[i] = eval.Min(positive.Grade[i], alts[j].Grade[i])
					}
				}
			}
		}(i, c)
	}
	wg.Wait()
	return positive, err
}

func negativeIdealRateInterval(alts []matrix.Alternative, Criteria []matrix.Criterion) (matrix.Alternative, error) {
	return positiveIdealRateInterval(alts, matrix.ChangeTypes(Criteria))
}

func positiveIdealRateNumber(alts []matrix.Alternative, Criteria []matrix.Criterion) (matrix.Alternative, error) {
	var wg sync.WaitGroup
	var err error = nil
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	positive := matrix.Alternative{Grade: make([]eval.Rating, len(alts[0].Grade)), CountOfCriteria: len(Criteria)}

	wg.Add(len(Criteria))
	for i, c := range Criteria {
		go func(i int, c matrix.Criterion) {
			defer wg.Done()

			for j := range alts {
				select {
				case <-ctx.Done():
					return
				default:
					if alts[j].Grade[i].GetType() != eval.NumbersMin.GetType() &&
						alts[j].Grade[i].GetType() != (eval.Interval{}).GetType() {
						err = v.IncompatibleTypes
						cancel()
						return
					}

					if positive.Grade[i].IsNil() && c.TypeOfCriteria == v.Benefit {
						positive.Grade[i].Evaluated = eval.NumbersMin
					} else if positive.Grade[i].IsNil() && c.TypeOfCriteria == v.Cost {
						positive.Grade[i].Evaluated = eval.NumbersMax
					}

					if alts[j].Grade[i].GetType() == (eval.Interval{}).GetType() {
						if c.TypeOfCriteria == v.Benefit {
							positive.Grade[i] = eval.Max(positive.Grade[i].ConvertToNumber(),
								alts[j].Grade[i].ConvertToInterval().End)
						} else {
							positive.Grade[i] = eval.Min(positive.Grade[i].ConvertToNumber(),
								alts[j].Grade[i].ConvertToInterval().Start)
						}
					} else {
						if c.TypeOfCriteria == v.Benefit {
							positive.Grade[i] = eval.Max(positive.Grade[i], alts[j].Grade[i])
						} else {
							positive.Grade[i] = eval.Min(positive.Grade[i], alts[j].Grade[i])
						}
					}
				}
			}
		}(i, c)
	}

	wg.Wait()
	return positive, err
}

func negativeIdealRateNumber(alts []matrix.Alternative, Criteria []matrix.Criterion) (matrix.Alternative, error) {
	return positiveIdealRateNumber(alts, matrix.ChangeTypes(Criteria))
}
