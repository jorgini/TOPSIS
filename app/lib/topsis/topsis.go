package topsis

import (
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
			var inerr error
			tm.PositiveIdeal, inerr = PositiveIdealRateInterval(tm.Data, tm.Criteria)
			if inerr != nil {
				err = inerr
			}
		}()

		go func() {
			defer wg.Done()
			var inerr error
			tm.NegativeIdeal, inerr = NegativeIdealRateInterval(tm.Data, tm.Criteria)
			if inerr != nil {
				err = inerr
			}
		}()
	} else if tm.Data[0].Grade[0].GetType() == (&eval.T1FS{}).GetType() ||
		tm.Data[0].Grade[0].GetType() == (&eval.IT2FS{}).GetType() {
		go func() {
			defer wg.Done()
			var inerr error
			tm.PositiveIdeal, inerr = PositiveIdealRateT1FS(tm.Data, tm.Criteria, tm.FormFs)
			if inerr != nil {
				err = inerr
			}
		}()

		go func() {
			defer wg.Done()
			var inerr error
			tm.NegativeIdeal, inerr = NegativeIdealRateT1FS(tm.Data, tm.Criteria, tm.FormFs)
			if inerr != nil {
				err = inerr
			}
		}()
	} else if tm.Data[0].Grade[0].GetType() == (&eval.AIFS{}).GetType() {
		go func() {
			defer wg.Done()
			var inerr error
			tm.PositiveIdeal, inerr = PositiveIdealRateAIFS(tm.Data, tm.Criteria)
			if inerr != nil {
				err = inerr
			}
		}()

		go func() {
			defer wg.Done()
			var inerr error
			tm.NegativeIdeal, inerr = NegativeIdealRateAIFS(tm.Data, tm.Criteria)
			if inerr != nil {
				err = inerr
			}
		}()
	} else {
		go func() {
			defer wg.Done()
			var inerr error
			tm.PositiveIdeal, inerr = PositiveIdealRateNumber(tm.Data, tm.Criteria)
			if inerr != nil {
				err = inerr
			}
		}()

		go func() {
			defer wg.Done()
			var inerr error
			tm.NegativeIdeal, inerr = NegativeIdealRateNumber(tm.Data, tm.Criteria)
			if inerr != nil {
				err = inerr
			}
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

	wg.Add(tm.CountAlternatives)
	for i := 0; i < tm.CountAlternatives; i++ {
		go func(i int) {
			defer wg.Done()

			var inerr error
			if (tm.HighType != (&eval.T1FS{}).GetType() && tm.HighType != (&eval.IT2FS{}).GetType() &&
				tm.HighType != (&eval.AIFS{}).GetType()) || vt == v.AlphaSlices {
				if vi == v.Default {
					if tm.DistancesToPositive[i].Evaluated, inerr = tm.Data[i].NumberMetric(tm.PositiveIdeal, vn); inerr != nil {
						err = inerr
						return
					}

					if tm.DistancesToNegative[i].Evaluated, inerr = tm.Data[i].NumberMetric(tm.NegativeIdeal, vn); inerr != nil {
						err = inerr
						return
					}
				} else if vi == v.Sengupta {
					if tm.DistancesToPositive[i].Evaluated, inerr = tm.Data[i].IntervalMetric(tm.PositiveIdeal, tm.Criteria, vn); inerr != nil {
						err = inerr
						return
					}
					if tm.DistancesToNegative[i].Evaluated, inerr = tm.Data[i].IntervalMetric(tm.NegativeIdeal,
						matrix.ChangeTypes(tm.Criteria), vn); inerr != nil {
						err = inerr
						return
					}
				} else {
					err = v.InvalidCaseOfOperation
					return
				}
			} else {
				if tm.DistancesToPositive[i].Evaluated, inerr = tm.Data[i].FSMetric(tm.PositiveIdeal, vn); inerr != nil {
					err = inerr
					return
				}

				if tm.DistancesToNegative[i].Evaluated, inerr = tm.Data[i].FSMetric(tm.NegativeIdeal, vn); inerr != nil {
					err = inerr
					return
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
