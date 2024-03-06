package lib

import (
	"context"
	"reflect"
	"sync"
)

type Variants int

const (
	NormalizeWeightsByMidPoint = 0
	NormalizeWithSum           = 1
	NormalizeWithMax           = 2
	Sengupta                   = 3
	AlphaSlices                = 4
	SqrtDistance               = 8
	CbrtDistance               = 9
	Default                    = 10
	Positive                   = true
	Negative                   = false
)

func (tm *TopsisMatrix) FindIdeals(v Variants) error {
	var err error
	var wg sync.WaitGroup

	wg.Add(2)
	if reflect.TypeOf(tm.data[0].grade[0]) == reflect.TypeOf(Interval{}) && v == Sengupta {
		go func() {
			defer wg.Done()
			tm.positiveIdeal, err = positiveIdealRateInterval(tm.data, tm.criteria)
		}()

		go func() {
			defer wg.Done()
			tm.negativeIdeal, err = negativeIdealRateInterval(tm.data, tm.criteria)
		}()
	} else if reflect.TypeOf(tm.data[0].grade[0]) == reflect.TypeOf(&T1FS{}) ||
		reflect.TypeOf(tm.data[0].grade[0]) == reflect.TypeOf(&IT2FS{}) {
		go func() {
			defer wg.Done()
			tm.positiveIdeal, err = positiveIdealRateT1FS(tm.data, tm.criteria, tm.formFs)
		}()

		go func() {
			defer wg.Done()
			tm.negativeIdeal, err = negativeIdealRateT1FS(tm.data, tm.criteria, tm.formFs)
		}()
	} else if reflect.TypeOf(tm.data[0].grade[0]) == reflect.TypeOf(&AIFS{}) {
		go func() {
			defer wg.Done()
			tm.positiveIdeal, err = positiveIdealRateAIFS(tm.data, tm.criteria)
		}()

		go func() {
			defer wg.Done()
			tm.negativeIdeal, err = negativeIdealRateAIFS(tm.data, tm.criteria)
		}()
	} else {
		go func() {
			defer wg.Done()
			tm.positiveIdeal, err = positiveIdealRateNumber(tm.data, tm.criteria)
		}()

		go func() {
			defer wg.Done()
			tm.negativeIdeal, err = negativeIdealRateNumber(tm.data, tm.criteria)
		}()
	}

	wg.Wait()
	if err == nil {
		tm.idealsFind = true
	}
	return err
}

func (tm *TopsisMatrix) FindDistanceToIdeals(vt, vi, vn Variants) error {
	var wg sync.WaitGroup
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(tm.countAlternatives)
	for i := 0; i < tm.countAlternatives; i++ {
		go func(i int) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			default:
				if (tm.highType != reflect.TypeOf(&T1FS{}) && tm.highType != reflect.TypeOf(&IT2FS{}) &&
					tm.highType != reflect.TypeOf(&AIFS{})) || vt == AlphaSlices {
					if vi == Default {
						if tm.distancesToPositive[i], err = tm.data[i].NumberMetric(tm.positiveIdeal, vn); err != nil {
							cancel()
							return
						}

						if tm.distancesToNegative[i], err = tm.data[i].NumberMetric(tm.negativeIdeal, vn); err != nil {
							cancel()
							return
						}
					} else if vi == Sengupta {
						if tm.distancesToPositive[i], err = tm.data[i].IntervalMetric(tm.positiveIdeal,
							tm.criteria, vn); err != nil {
							cancel()
							return
						}

						if tm.distancesToNegative[i], err = tm.data[i].IntervalMetric(tm.negativeIdeal,
							ChangeTypes(tm.criteria), vn); err != nil {
							cancel()
							return
						}
					} else {
						err = InvalidCaseOfOperation
						cancel()
						return
					}
				} else {
					if tm.distancesToPositive[i], err = tm.data[i].FSMetric(tm.positiveIdeal, vn); err != nil {
						cancel()
						return
					}

					if tm.distancesToNegative[i], err = tm.data[i].FSMetric(tm.negativeIdeal, vn); err != nil {
						cancel()
						return
					}
				}
			}
		}(i)
	}

	wg.Wait()
	if err == nil {
		tm.distancesFind = true
	}
	return err
}

func (tm *TopsisMatrix) CalcCloseness() {
	var wg sync.WaitGroup

	wg.Add(tm.countAlternatives)
	for i := 0; i < tm.countAlternatives; i++ {
		go func(i int) {
			defer wg.Done()
			if reflect.TypeOf(tm.distancesToPositive[i]) == reflect.TypeOf(Interval{}) {
				neg := tm.distancesToNegative[i].ConvertToInterval()
				pos := tm.distancesToPositive[i].ConvertToInterval()
				tm.relativeCloseness[i] = Interval{neg.Start / (neg.End + pos.End),
					neg.End / (neg.Start + pos.Start)}
			} else {
				tm.relativeCloseness[i] = tm.distancesToNegative[i].ConvertToNumber() /
					(tm.distancesToNegative[i].ConvertToNumber() + tm.distancesToPositive[i].ConvertToNumber())
			}
		}(i)
	}
	wg.Wait()
	tm.closenessFind = true
}
