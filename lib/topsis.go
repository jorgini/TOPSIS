package lib

import (
	"errors"
	"reflect"
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
	if reflect.TypeOf(tm.data[0].grade[0]) == reflect.TypeOf(Interval{}) && v == Sengupta {
		tm.positiveIdeal, err = positiveIdealRateInterval(tm.data, tm.criteria)
		tm.negativeIdeal, err = negativeIdealRateInterval(tm.data, tm.criteria)

		if err != nil {
			return err
		}
	} else if reflect.TypeOf(tm.data[0].grade[0]) == reflect.TypeOf(&T1FS{}) ||
		reflect.TypeOf(tm.data[0].grade[0]) == reflect.TypeOf(&IT2FS{}) {
		tm.positiveIdeal, err = positiveIdealRateT1FS(tm.data, tm.criteria, tm.formFs)
		tm.negativeIdeal, err = negativeIdealRateT1FS(tm.data, tm.criteria, tm.formFs)

		if err != nil {
			return err
		}
	} else if reflect.TypeOf(tm.data[0].grade[0]) == reflect.TypeOf(&AIFS{}) {
		tm.positiveIdeal, err = positiveIdealRateAIFS(tm.data, tm.criteria)
		tm.negativeIdeal, err = negativeIdealRateAIFS(tm.data, tm.criteria)

		if err != nil {
			return err
		}
	} else {
		tm.positiveIdeal, err = positiveIdealRateNumber(tm.data, tm.criteria)
		tm.negativeIdeal, err = negativeIdealRateNumber(tm.data, tm.criteria)

		if err != nil {
			return err
		}
	}

	tm.idealsFind = true
	return nil
}

func (tm *TopsisMatrix) FindDistanceToIdeals(vt, vi, vn Variants) error {
	for i := 0; i < tm.countAlternatives; i++ {
		var err error
		if (tm.highType != reflect.TypeOf(&T1FS{}) && tm.highType != reflect.TypeOf(&IT2FS{}) &&
			tm.highType != reflect.TypeOf(&AIFS{})) || vt == AlphaSlices {
			if vi == Default {
				if tm.distancesToPositive[i], err = tm.data[i].NumberMetric(tm.positiveIdeal, vn); err != nil {
					return errors.Join(err)
				}

				if tm.distancesToNegative[i], err = tm.data[i].NumberMetric(tm.negativeIdeal, vn); err != nil {
					return errors.Join(err)
				}
			} else if vi == Sengupta {
				if tm.distancesToPositive[i], err = tm.data[i].IntervalMetric(tm.positiveIdeal,
					tm.criteria, Positive, vn); err != nil {
					return errors.Join(err)
				}

				if tm.distancesToNegative[i], err = tm.data[i].IntervalMetric(tm.negativeIdeal,
					tm.criteria, Negative, vn); err != nil {
					return errors.Join(err)
				}
			} else {
				return InvalidCaseOfOperation
			}
		} else {
			if tm.distancesToPositive[i], err = tm.data[i].FSMetric(tm.positiveIdeal, vn); err != nil {
				return errors.Join(err)
			}

			if tm.distancesToNegative[i], err = tm.data[i].FSMetric(tm.negativeIdeal, vn); err != nil {
				return errors.Join(err)
			}
		}
	}
	tm.distancesFind = true
	return nil
}

func (tm *TopsisMatrix) CalcCloseness() {
	for i := 0; i < tm.countAlternatives; i++ {
		if reflect.TypeOf(tm.distancesToPositive[i]) == reflect.TypeOf(Interval{}) {
			neg := tm.distancesToNegative[i].ConvertToInterval()
			pos := tm.distancesToPositive[i].ConvertToInterval()
			tm.relativeCloseness[i] = Interval{neg.Start / (neg.End + pos.End),
				neg.End / (neg.Start + pos.Start)}
		} else {
			tm.relativeCloseness[i] = tm.distancesToNegative[i].ConvertToNumber() /
				(tm.distancesToNegative[i].ConvertToNumber() + tm.distancesToPositive[i].ConvertToNumber())
		}
	}
	tm.closenessFind = true
}
