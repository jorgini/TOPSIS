package eval

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

var (
	Bad       = 0
	Normal    = 1
	Good      = 2
	Excellent = 3
)

type LinguisticScale struct {
	Ratings []Rating `json:"ratings"`
	Marks   []string `json:"marks"`
}

var DefaultNumberScale = &LinguisticScale{
	Ratings: []Rating{{Number(1)}, {Number(3)}, {Number(6)}, {Number(9)}},
	Marks:   []string{"Bad", "Normal", "Good", "Excellent"},
}

var DefaultIntervalScale = &LinguisticScale{
	Ratings: []Rating{{Evaluated: Interval{Start: 1, End: 3}}, {Evaluated: Interval{Start: 3, End: 5}},
		{Evaluated: Interval{Start: 5, End: 7}}, {Evaluated: Interval{Start: 7, End: 9}}},
	Marks: []string{"Bad", "Normal", "Good", "Excellent"},
}

var DefaultT1FSScale = &LinguisticScale{
	Ratings: []Rating{{NewT1FS(0, 2, 3)}, {NewT1FS(3, 4, 5)},
		{NewT1FS(5, 6, 7)}, {NewT1FS(8, 9, 10)}},
	Marks: []string{"Bad", "Normal", "Good", "Excellent"},
}

func (l LinguisticScale) Value() (driver.Value, error) {
	data, err := json.Marshal(l)
	return string(data), err
}

func (l *LinguisticScale) Scan(src interface{}) error {
	var ling LinguisticScale
	var err error
	switch src.(type) {
	case string:
		err = json.Unmarshal([]byte(src.(string)), &ling)
	case []byte:
		err = json.Unmarshal(src.([]byte), &ling)
	default:
		return errors.New("incompatible type for LinguisticScale")
	}
	if err != nil {
		return err
	}
	*l = ling
	return nil
}

type Linguistic struct {
	Mark string `json:"mark"`
	Rating
}

func (l Linguistic) CopyEval() Rating {
	return Rating{Linguistic{Mark: l.Mark, Rating: l.Rating.CopyEval()}}
}
