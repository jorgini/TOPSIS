package topsis

import (
	"sort"
	"sync"
	"webApp/lib/eval"
	"webApp/lib/matrix"
	v "webApp/lib/variables"
)

func (tm *TopsisMatrix) FindIdeals(variants v.Variants, g int) error {
	var err error
	var wg sync.WaitGroup
	if g > 1 {
		g /= 2
	}

	wg.Add(2)
	if tm.Data[0].Grade[0].GetType() == (eval.Interval{}).GetType() && variants == v.Sengupta {
		go func() {
			defer wg.Done()
			var inerr error
			tm.PositiveIdeal, inerr = positiveIdeal(tm.Data, tm.Criteria, eval.Interval{}.GetType(), v.None, g)
			if inerr != nil {
				err = inerr
			}
		}()

		go func() {
			defer wg.Done()
			var inerr error
			tm.NegativeIdeal, inerr = negativeIdeal(tm.Data, tm.Criteria, eval.Interval{}.GetType(), v.None, g)
			if inerr != nil {
				err = inerr
			}
		}()
	} else if tm.Data[0].Grade[0].GetType() == (&eval.T1FS{}).GetType() ||
		tm.Data[0].Grade[0].GetType() == (&eval.IT2FS{}).GetType() {
		go func() {
			defer wg.Done()
			var inerr error
			tm.PositiveIdeal, inerr = positiveIdeal(tm.Data, tm.Criteria, (&eval.T1FS{}).GetType(), tm.FormFs, g)
			if inerr != nil {
				err = inerr
			}
		}()

		go func() {
			defer wg.Done()
			var inerr error
			tm.NegativeIdeal, inerr = negativeIdeal(tm.Data, tm.Criteria, (&eval.T1FS{}).GetType(), tm.FormFs, g)
			if inerr != nil {
				err = inerr
			}
		}()
	} else if tm.Data[0].Grade[0].GetType() == (&eval.AIFS{}).GetType() {
		go func() {
			defer wg.Done()
			var inerr error
			tm.PositiveIdeal, inerr = positiveIdeal(tm.Data, tm.Criteria, (&eval.AIFS{}).GetType(), v.None, g)
			if inerr != nil {
				err = inerr
			}
		}()

		go func() {
			defer wg.Done()
			var inerr error
			tm.NegativeIdeal, inerr = negativeIdeal(tm.Data, tm.Criteria, (&eval.AIFS{}).GetType(), v.None, g)
			if inerr != nil {
				err = inerr
			}
		}()
	} else {
		go func() {
			defer wg.Done()
			var inerr error
			tm.PositiveIdeal, inerr = positiveIdeal(tm.Data, tm.Criteria, eval.NumbersMin.GetType(), v.None, g)
			if inerr != nil {
				err = inerr
			}
		}()

		go func() {
			defer wg.Done()
			var inerr error
			tm.NegativeIdeal, inerr = negativeIdeal(tm.Data, tm.Criteria, eval.NumbersMin.GetType(), v.None, g)
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

func (tm *TopsisMatrix) FindDistanceToIdeals(vt, vi, vn v.Variants, g int) error {
	var wg sync.WaitGroup
	var err error
	if g > tm.CountAlternatives {
		g = tm.CountAlternatives
	}
	off := tm.CountAlternatives / g
	wg.Add(g)
	for b := 0; b < g; b++ {
		go func(b int) {
			defer wg.Done()

			start := b * off
			end := (b + 1) * off
			if b == g-1 {
				end = tm.CountAlternatives
			}

			var inerr error

			for i := start; i < end; i++ {
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
			}
		}(b)
	}
	wg.Wait()

	if err == nil {
		tm.DistancesFind = true
	}
	return err
}

func (tm *TopsisMatrix) CalcCloseness(g int) {
	var wg sync.WaitGroup
	if g > tm.CountAlternatives {
		g = tm.CountAlternatives
	}
	off := tm.CountAlternatives / g
	wg.Add(g)
	for b := 0; b < g; b++ {
		go func(b int) {
			defer wg.Done()

			start := b * off
			end := (b + 1) * off
			if b == g-1 {
				end = tm.CountAlternatives
			}

			for i := start; i < end; i++ {
				if tm.DistancesToPositive[i].GetType() == (eval.Interval{}).GetType() {
					neg := tm.DistancesToNegative[i].ConvertToInterval()
					pos := tm.DistancesToPositive[i].ConvertToInterval()
					tm.RelativeCloseness[i].Evaluated = eval.Interval{Start: neg.Start / (neg.End + pos.End),
						End: neg.End / (neg.Start + pos.Start)}
				} else {
					tm.RelativeCloseness[i].Evaluated = tm.DistancesToNegative[i].ConvertToNumber() /
						(tm.DistancesToNegative[i].ConvertToNumber() + tm.DistancesToPositive[i].ConvertToNumber())
				}
			}
		}(b)
	}
	wg.Wait()
	tm.ClosenessFind = true
}

func (tm *TopsisMatrix) RankedList(ranking v.Variants) matrix.RankedList {
	set := make([]eval.Rating, len(tm.RelativeCloseness))
	ind := make([]int, len(tm.RelativeCloseness))
	for i := range set {
		ind[i] = i
		set[i] = tm.RelativeCloseness[i].CopyEval()
	}

	if ranking == v.Sengupta {
		sort.Slice(ind, func(i, j int) bool {
			l := set[ind[i]].ConvertToInterval()
			r := set[ind[j]].ConvertToInterval()
			return l.SenguptaGeq(r)
		})
		sort.Slice(set, func(i, j int) bool {
			l := set[i].ConvertToInterval()
			r := set[j].ConvertToInterval()
			return l.SenguptaGeq(r)
		})
	} else {
		sort.Slice(ind, func(i, j int) bool {
			return set[ind[i]].ConvertToNumber() > set[ind[j]].ConvertToNumber()
		})
		sort.Slice(set, func(i, j int) bool {
			return set[i].ConvertToNumber() > set[j].ConvertToNumber()
		})
	}
	return matrix.RankedList{Coeffs: set, Order: ind}
}
