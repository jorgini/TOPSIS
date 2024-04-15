package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"time"
	"webApp/lib"
	"webApp/lib/eval"
	v "webApp/lib/variables"
)

const (
	Draft      = false
	Complete   = true
	Individual = "individual"
	Group      = "group"
)

type TaskModel struct {
	SID            int64                `json:"sid,omitempty" db:"sid"`
	MaintainerID   int64                `json:"uid,omitempty" db:"maintainer"`
	Password       string               `json:"password,omitempty" db:"password"`
	Title          string               `json:"title" db:"title"`
	Description    string               `json:"description,omitempty" db:"description"`
	LastChangesAt  time.Time            `json:"last_change" db:"last_change"`
	TaskType       string               `json:"task_type" db:"task_type"`
	Method         string               `json:"method" db:"method"`
	CalcSettings   int64                `json:"calc_settings" db:"calc_settings"`
	LingScale      eval.LinguisticScale `json:"ling_scale" db:"ling_scale"`
	Alternatives   Alts                 `json:"alternatives" db:"alternatives"`
	Criteria       Criteria             `json:"criteria" db:"criteria"`
	ExpertsWeights Weights              `json:"experts_weights" db:"experts_weights"`
	Status         bool                 `json:"status" db:"status"`
}

func GetDefaultTask(title string, uid int64) TaskModel {
	settings := lib.CalcSettings{
		ValueNorm:   v.NormalizeWithSum,
		WeighNorm:   v.NormalizeWithSum,
		RankingAlg:  v.Default,
		FsDist:      v.Default,
		IntDist:     v.Default,
		NumDist:     v.SqrtDistance,
		Aggregating: v.AggregateMatrix,
	}
	logrus.Info(settings.Comprise())
	task := TaskModel{
		Title:        title,
		MaintainerID: uid,
		Description:  "",
		Password:     "qwerty",
		TaskType:     Individual,
		Method:       v.SMART,
		CalcSettings: settings.Comprise(),
		LingScale:    *eval.DefaultNumberScale,
		Status:       Draft,
	}
	return task
}

type TaskShortCard struct {
	SID         int64     `json:"sid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Method      string    `json:"method"`
	TaskType    string    `json:"taskType"`
	LastChange  time.Time `json:"last_change"`
	Status      bool      `json:"status"`
}

type AlternativeModel struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (a *AlternativeModel) UnmarshalJSON(data []byte) error {
	result := struct {
		Title       *string `json:"title"`
		Description string  `json:"description"`
	}{}

	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.Title == nil {
		return errors.New("invalid value for required title")
	} else {
		a.Title = *result.Title
		a.Description = result.Description
	}
	return nil
}

type Alts []AlternativeModel

func (a Alts) Value() (driver.Value, error) {
	data, err := json.Marshal(a)
	return string(data), err
}

func (a *Alts) Scan(src interface{}) error {
	var tmp Alts
	var err error
	switch src.(type) {
	case string:
		err = json.Unmarshal([]byte(src.(string)), &tmp)
	case []byte:
		err = json.Unmarshal(src.([]byte), &tmp)
	case nil:
		return nil
	default:
		return errors.New("incompatible type for Alts")
	}
	if err != nil {
		return err
	}
	*a = tmp
	return nil
}

type CriterionModel struct {
	Title           string      `json:"title"`
	Description     string      `json:"description"`
	Weight          eval.Rating `json:"weight"`
	TypeOfCriterion bool        `json:"type_of_criterion"`
}

func (c *CriterionModel) UnmarshalJSON(data []byte) error {
	result := struct {
		Title           *string      `json:"title"`
		Description     string       `json:"description"`
		Weight          *eval.Rating `json:"weight"`
		TypeOfCriterion *bool        `json:"type_of_criterion"`
	}{}

	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.Title == nil || result.Weight == nil || result.TypeOfCriterion == nil {
		return errors.New("invalid value for required criterion field")
	} else {
		c.Title = *result.Title
		c.Description = result.Description
		c.Weight = *result.Weight
		c.TypeOfCriterion = *result.TypeOfCriterion
	}
	return nil
}

type Criteria []CriterionModel

func (c Criteria) Value() (driver.Value, error) {
	data, err := json.Marshal(c)
	return string(data), err
}

func (c *Criteria) Scan(src interface{}) error {
	var tmp Criteria
	var err error
	switch src.(type) {
	case string:
		err = json.Unmarshal([]byte(src.(string)), &tmp)
	case []byte:
		err = json.Unmarshal(src.([]byte), &tmp)
	case nil:
		return nil
	default:
		return errors.New("incompatible type for Criteria")
	}
	if err != nil {
		return err
	}
	*c = tmp
	return nil
}

type Weights []eval.Rating

func (w Weights) Value() (driver.Value, error) {
	data, err := json.Marshal(w)
	return string(data), err
}

func (w *Weights) Scan(src interface{}) error {
	var tmp Weights
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
	*w = tmp
	return nil
}
