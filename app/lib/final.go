package lib

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"math/rand"
	"sync"
	"time"
	"webApp/lib/eval"
	"webApp/lib/matrix"
	"webApp/lib/smart"
	"webApp/lib/topsis"
	v "webApp/lib/variables"
)

type SensitivityResult struct {
	Results   []matrix.RankedList
	Threshold float64
}

func (s SensitivityResult) Value() (driver.Value, error) {
	data, err := json.Marshal(s)
	return string(data), err
}

func (s *SensitivityResult) Scan(src interface{}) error {
	var sens SensitivityResult
	var err error
	switch src.(type) {
	case string:
		err = json.Unmarshal([]byte(src.(string)), &sens)
	case []byte:
		err = json.Unmarshal(src.([]byte), &sens)
	default:
		return errors.New("incompatible type for SensitivityResult")
	}
	if err != nil {
		return err
	}
	*s = sens
	return nil
}

type CalcSettings struct {
	ValueNorm   v.Variants
	WeighNorm   v.Variants
	RankingAlg  v.Variants
	FsDist      v.Variants
	IntDist     v.Variants
	NumDist     v.Variants
	Aggregating v.Variants
}

func (c *CalcSettings) Comprise() int64 {
	result := c.ValueNorm | (c.WeighNorm << 4) | (c.RankingAlg << 8) | (c.FsDist << 12) |
		(c.IntDist << 16) | (c.NumDist << 20) | (c.Aggregating << 24)
	return int64(result)
}

func (c *CalcSettings) Parse(settings int64) {
	c.ValueNorm = v.Variants(settings & 0b1111)
	c.WeighNorm = v.Variants((settings >> 4) & 0b1111)
	c.RankingAlg = v.Variants((settings >> 8) & 0b1111)
	c.FsDist = v.Variants((settings >> 12) & 0b1111)
	c.IntDist = v.Variants((settings >> 16) & 0b1111)
	c.NumDist = v.Variants((settings >> 20) & 0b1111)
	c.Aggregating = v.Variants((settings >> 24) & 0b1111)
}

func SensAnalysis(method v.Method, calcSettings int64, threshold float64, mxs []matrix.Matrix, w []eval.Rating) (*SensitivityResult, error) {
	if len(mxs) != len(w) {
		return nil, v.InvalidSize
	}

	settings := CalcSettings{}
	settings.Parse(calcSettings)

	result := SensitivityResult{Results: make([]matrix.RankedList, 10), Threshold: threshold}
	gen := rand.New(rand.NewSource(time.Now().Unix()))
	var err error
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()

			changeMatrices := make([]matrix.Matrix, len(mxs))
			weights := make([]eval.Evaluated, len(w))
			for i := range changeMatrices {
				weights[i] = w[i].Evaluated
				changeMatrices[i] = randomChange(&mxs[i], threshold, gen)
			}

			var inerr error
			if method == v.TOPSIS {
				result.Results[i], inerr = TopsisFullCalc(settings, changeMatrices, weights)
				if inerr != nil {
					err = inerr
					return
				}
			} else if method == v.SMART {
				result.Results[i], inerr = SmartFullCalc(settings, changeMatrices, weights)
				if inerr != nil {
					err = inerr
					return
				}
			} else {
				err = errors.New("invalid method")
				return
			}
		}(i)
	}
	wg.Wait()

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func randomChange(m *matrix.Matrix, threshold float64, gen *rand.Rand) matrix.Matrix {
	newMatrix := matrix.NewMatrix(m.CountAlternatives, m.CountCriteria)

	for i := range newMatrix.Data {
		for j := range newMatrix.Data[i].Grade {
			gap := gen.Float64() * threshold
			if gen.Int()%2 == 1 {
				gap = -gap
			}
			_ = newMatrix.SetValue(m.Data[i].Grade[j].Weighted(eval.Number(1+gap)), i, j)
		}
	}
	_ = newMatrix.SetCriteria(m.Criteria)
	return *newMatrix
}

func TopsisFullCalc(settings CalcSettings, mxs []matrix.Matrix, weights []eval.Evaluated) (matrix.RankedList, error) {
	var err error
	if settings.Aggregating == v.AggregateMatrix {
		aggMatrix, err := matrix.AggregateRatings(mxs, weights)
		if err != nil {
			return matrix.RankedList{}, err
		}

		resultMatrix := topsis.ConvertToTopsisMatrix(aggMatrix)
		if err = resultMatrix.Normalization(settings.ValueNorm, settings.WeighNorm); err != nil {
			return matrix.RankedList{}, err
		}

		resultMatrix.CalcWeightedMatrix()

		if err = resultMatrix.FindIdeals(settings.RankingAlg); err != nil {
			return matrix.RankedList{}, err
		}

		if err = resultMatrix.FindDistanceToIdeals(settings.FsDist, settings.IntDist, settings.NumDist); err != nil {
			return matrix.RankedList{}, err
		}

		resultMatrix.CalcCloseness()
		return resultMatrix.RankedList(settings.RankingAlg), nil
	} else if settings.Aggregating == v.AggregateFinals {
		matrices := make([]topsis.TopsisMatrix, len(mxs))
		var wg sync.WaitGroup
		wg.Add(len(mxs))
		for i := range mxs {
			go func(i int) {
				defer wg.Done()
				matrices[i] = *topsis.ConvertToTopsisMatrix(&mxs[i])
				if inerr := matrices[i].Normalization(settings.ValueNorm, settings.WeighNorm); inerr != nil {
					err = inerr
					return
				}

				matrices[i].CalcWeightedMatrix()

				if inerr := matrices[i].FindIdeals(settings.RankingAlg); inerr != nil {
					err = inerr
					return
				}

				if inerr := matrices[i].FindDistanceToIdeals(settings.FsDist, settings.IntDist, settings.NumDist); inerr != nil {
					err = inerr
					return
				}
			}(i)
		}
		wg.Wait()
		if err != nil {
			return matrix.RankedList{}, err
		}

		aggMatrix, err := topsis.AggregateDistances(matrices, weights)
		if err != nil {
			return matrix.RankedList{}, err
		}

		aggMatrix.CalcCloseness()
		return aggMatrix.RankedList(settings.RankingAlg), nil
	} else {
		return matrix.RankedList{}, v.InvalidCaseOfOperation
	}
}

func SmartFullCalc(settings CalcSettings, mxs []matrix.Matrix, weights []eval.Evaluated) (matrix.RankedList, error) {
	var err error
	if settings.Aggregating == v.AggregateMatrix {
		aggMatrix, err := matrix.AggregateRatings(mxs, weights)
		if err != nil {
			return matrix.RankedList{}, err
		}

		resultMatrix := smart.ConvertToSmartMatrix(aggMatrix)

		if err := resultMatrix.Normalization(settings.ValueNorm, settings.WeighNorm); err != nil {
			return matrix.RankedList{}, err
		}

		resultMatrix.CalcWeightedMatrix()

		resultMatrix.CalcFinalScore()

		return resultMatrix.RankedList(settings.RankingAlg), nil
	} else if settings.Aggregating == v.AggregateFinals {
		matrices := make([]smart.SmartMatrix, len(mxs))

		var wg sync.WaitGroup
		wg.Add(len(mxs))
		for i := range mxs {
			go func(i int) {
				defer wg.Done()
				matrices[i] = *smart.ConvertToSmartMatrix(&mxs[i])
				if inerr := matrices[i].Normalization(settings.ValueNorm, settings.WeighNorm); inerr != nil {
					err = inerr
					return
				}

				matrices[i].CalcWeightedMatrix()

				matrices[i].CalcFinalScore()
			}(i)
		}
		wg.Wait()
		if err != nil {
			return matrix.RankedList{}, err
		}

		result, err := smart.AggregateScores(matrices, weights)
		if err != nil {
			return matrix.RankedList{}, err
		}

		return result.RankedList(settings.RankingAlg), nil
	} else {
		return matrix.RankedList{}, v.InvalidCaseOfOperation
	}
}
