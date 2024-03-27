package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
	"webApp/lib"
	"webApp/lib/eval"
	"webApp/lib/matrix"
	v "webApp/lib/variables"
)

type FinalModel struct {
	FID          int64                 `json:"fid" db:"fid"`
	Result       Coefficients          `json:"result" db:"result"`
	SensAnalysis lib.SensitivityResult `json:"sens_analysis" db:"sens_analysis"`
	Threshold    float64               `json:"threshold" db:"threshold"`
	LastChange   time.Time             `json:"last_change" db:"last_change"`
}

type Coefficients []eval.Rating

func (c Coefficients) Value() (driver.Value, error) {
	data, err := json.Marshal(c)
	return string(data), err
}

func (c *Coefficients) Scan(src interface{}) error {
	var tmp Coefficients
	var err error
	switch src.(type) {
	case string:
		err = json.Unmarshal([]byte(src.(string)), &tmp)
	case []byte:
		err = json.Unmarshal(src.([]byte), &tmp)
	case nil:
		return nil
	default:
		return errors.New("incompatible type for Weights")
	}
	if err != nil {
		return err
	}
	*c = tmp
	return nil
}

func CalcFinal(matrices []MatrixModel, task *TaskModel, threshold float64) (*FinalModel, error) {
	settings := lib.CalcSettings{}
	settings.Parse(task.CalcSettings)
	mxs := ConvertModelToMatrix(matrices, task.Criteria)
	if mxs == nil {
		return nil, errors.New("incompatible sizes of matrices and criteria")
	}
	weights := ConvertRatingsToEvaluated(task.ExpertsWeights)

	var err error
	var coeffs Coefficients
	if task.Method == v.TOPSIS {
		coeffs, err = lib.TopsisFullCalc(settings, mxs, weights)
		if err != nil {
			return nil, err
		}
	} else if task.Method == v.SMART {
		coeffs, err = lib.SmartFullCalc(settings, mxs, weights)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("invalid method of task")
	}

	sens, err := lib.SensAnalysis(v.Method(task.Method), task.CalcSettings, threshold, mxs, task.ExpertsWeights)
	if err != nil {
		return nil, err
	}

	result := FinalModel{
		FID:          task.SID,
		Result:       coeffs,
		SensAnalysis: *sens,
		Threshold:    threshold,
		LastChange:   time.Now(),
	}
	return &result, nil
}

func ConvertModelToMatrix(models []MatrixModel, criteria Criteria) []matrix.Matrix {
	mxs := make([]matrix.Matrix, len(models))
	for i := range mxs {
		mxs[i] = *models[i].Matrix
		for j := range criteria {
			if err := mxs[i].SetCriterion(criteria[j].Weight, criteria[j].TypeOfCriterion, j); err != nil {
				return nil
			}
		}
	}
	return mxs
}

func ConvertRatingsToEvaluated(ratings []eval.Rating) []eval.Evaluated {
	evaluated := make([]eval.Evaluated, len(ratings))
	for i := range evaluated {
		evaluated[i] = ratings[i]
	}
	return evaluated
}
