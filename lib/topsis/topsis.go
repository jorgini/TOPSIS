package topsis

import (
	"context"
	"sync"
	"webApp/lib/eval"
	"webApp/lib/matrix"
	v "webApp/lib/variables"
)

func (tm *TopsisMatrix) FindIdeals(variants v.Variants) error {
	var err error
	var wg sync.WaitGroup

	wg.Add(2)
	if tm.Data[0].Grade[0].GetType() == (eval.Interval{}).GetType() && variants == v.Sengupta {
		go func() {
			defer wg.Done()
			tm.PositiveIdeal, err = PositiveIdealRateInterval(tm.Data, tm.Criteria)
		}()

		go func() {
			defer wg.Done()
			tm.NegativeIdeal, err = NegativeIdealRateInterval(tm.Data, tm.Criteria)
		}()
	} else if tm.Data[0].Grade[0].GetType() == (&eval.T1FS{}).GetType() ||
		tm.Data[0].Grade[0].GetType() == (&eval.IT2FS{}).GetType() {
		go func() {
			defer wg.Done()
			tm.PositiveIdeal, err = PositiveIdealRateT1FS(tm.Data, tm.Criteria, tm.FormFs)
		}()

		go func() {
			defer wg.Done()
			tm.NegativeIdeal, err = NegativeIdealRateT1FS(tm.Data, tm.Criteria, tm.FormFs)
		}()
	} else if tm.Data[0].Grade[0].GetType() == (&eval.AIFS{}).GetType() {
		go func() {
			defer wg.Done()
			tm.PositiveIdeal, err = PositiveIdealRateAIFS(tm.Data, tm.Criteria)
		}()

		go func() {
			defer wg.Done()
			tm.NegativeIdeal, err = NegativeIdealRateAIFS(tm.Data, tm.Criteria)
		}()
	} else {
		go func() {
			defer wg.Done()
			tm.PositiveIdeal, err = PositiveIdealRateNumber(tm.Data, tm.Criteria)
		}()

		go func() {
			defer wg.Done()
			tm.NegativeIdeal, err = NegativeIdealRateNumber(tm.Data, tm.Criteria)
		}()
	}

	wg.Wait()
	if err == nil {
		tm.IdealsFind = true
	}
	return err
}

func (tm *TopsisMatrix) FindDistanceToIdeals(vt, vi, vn v.Variants) error {
	var wg sync.WaitGroup
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(tm.CountAlternatives)
	for i := 0; i < tm.CountAlternatives; i++ {
		go func(i int) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			default:
				if (tm.HighType != (&eval.T1FS{}).GetType() && tm.HighType != (&eval.IT2FS{}).GetType() &&
					tm.HighType != (&eval.AIFS{}).GetType()) || vt == v.AlphaSlices {
					if vi == v.Default {
						if tm.DistancesToPositive[i].Evaluated, err = tm.Data[i].NumberMetric(tm.PositiveIdeal, vn); err != nil {
							cancel()
							return
						}

						if tm.DistancesToNegative[i].Evaluated, err = tm.Data[i].NumberMetric(tm.NegativeIdeal, vn); err != nil {
							cancel()
							return
						}
					} else if vi == v.Sengupta {
						if tm.DistancesToPositive[i].Evaluated, err = tm.Data[i].IntervalMetric(tm.PositiveIdeal,
							tm.Criteria, vn); err != nil {
							cancel()
							return
						}

						if tm.DistancesToNegative[i].Evaluated, err = tm.Data[i].IntervalMetric(tm.NegativeIdeal,
							matrix.ChangeTypes(tm.Criteria), vn); err != nil {
							cancel()
							return
						}
					} else {
						err = v.InvalidCaseOfOperation
						cancel()
						return
					}
				} else {
					if tm.DistancesToPositive[i].Evaluated, err = tm.Data[i].FSMetric(tm.PositiveIdeal, vn); err != nil {
						cancel()
						return
					}

					if tm.DistancesToNegative[i].Evaluated, err = tm.Data[i].FSMetric(tm.NegativeIdeal, vn); err != nil {
						cancel()
						return
					}
				}
			}
		}(i)
	}

	wg.Wait()
	if err == nil {
		tm.DistancesFind = true
	}
	return err
}

func (tm *TopsisMatrix) CalcCloseness() {
	var wg sync.WaitGroup

	wg.Add(tm.CountAlternatives)
	for i := 0; i < tm.CountAlternatives; i++ {
		go func(i int) {
			defer wg.Done()
			if tm.DistancesToPositive[i].GetType() == (eval.Interval{}).GetType() {
				neg := tm.DistancesToNegative[i].ConvertToInterval()
				pos := tm.DistancesToPositive[i].ConvertToInterval()
				tm.RelativeCloseness[i].Evaluated = eval.Interval{Start: neg.Start / (neg.End + pos.End),
					End: neg.End / (neg.Start + pos.Start)}
			} else {
				tm.RelativeCloseness[i].Evaluated = tm.DistancesToNegative[i].ConvertToNumber() /
					(tm.DistancesToNegative[i].ConvertToNumber() + tm.DistancesToPositive[i].ConvertToNumber())
			}
		}(i)
	}
	wg.Wait()
	tm.ClosenessFind = true
}
