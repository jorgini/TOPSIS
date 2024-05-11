package topsis

import (
	"webApp/lib/eval"
	"webApp/lib/matrix"
	v "webApp/lib/variables"
)

func positiveIdealRateAIFS0(alts []matrix.Alternative, Criteria []matrix.Criterion) (matrix.Alternative, error) {
	positive := matrix.Alternative{Grade: make([]eval.Rating, len(alts[0].Grade)), CountOfCriteria: len(Criteria)}

	for i, c := range Criteria {
		for j := range alts {
			if alts[j].Grade[i].GetType() != (&eval.AIFS{}).GetType() {
				return matrix.Alternative{}, v.IncompatibleTypes
			}

			altsGrade := alts[j].Grade[i].ConvertToAIFS(v.Default)

			if positive.Grade[i].IsNil() {
				positive.Grade[i].Evaluated = eval.NewAIFS(altsGrade.Pi, altsGrade.Vert...)
				continue
			}

			if c.TypeOfCriteria == v.Benefit {
				positive.Grade[i] = eval.Max(positive.Grade[i], alts[j].Grade[i])

			} else {
				positive.Grade[i] = eval.Min(positive.Grade[i], alts[j].Grade[i])
			}
		}
	}

	return positive, nil
}

func positiveIdealRateT1FS0(alts []matrix.Alternative, Criteria []matrix.Criterion, Form v.Variants) (matrix.Alternative, error) {
	positive := matrix.Alternative{Grade: make([]eval.Rating, len(alts[0].Grade)), CountOfCriteria: len(Criteria)}

	for i, c := range Criteria {
		for j := range alts {
			if alts[j].Grade[i].GetType() != (&eval.T1FS{}).GetType() && alts[j].Grade[i].GetType() != (&eval.IT2FS{}).GetType() {
				return matrix.Alternative{}, v.IncompatibleTypes
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
	}

	return positive, nil
}

func positiveIdealRateInterval0(alts []matrix.Alternative, Criteria []matrix.Criterion) (matrix.Alternative, error) {
	positive := matrix.Alternative{Grade: make([]eval.Rating, len(alts[0].Grade)), CountOfCriteria: len(Criteria)}

	for i, c := range Criteria {
		for j := range alts {
			if alts[j].Grade[i].GetType() != (eval.Interval{}).GetType() {
				return matrix.Alternative{}, v.IncompatibleTypes
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
	return positive, nil
}

func positiveIdealRateNumber0(alts []matrix.Alternative, Criteria []matrix.Criterion) (matrix.Alternative, error) {
	positive := matrix.Alternative{Grade: make([]eval.Rating, len(alts[0].Grade)), CountOfCriteria: len(Criteria)}

	for i, c := range Criteria {
		for j := range alts {
			if alts[j].Grade[i].GetType() != eval.NumbersMin.GetType() && alts[j].Grade[i].GetType() != (eval.Interval{}).GetType() {
				return matrix.Alternative{}, v.IncompatibleTypes
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

	return positive, nil
}
