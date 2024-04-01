package matrix

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"
	"webApp/lib/eval"
)

type RankedList struct {
	Coeffs []eval.Rating `json:"coeffs"`
	Order  []int         `json:"order"`
}

func (r RankedList) String() string {
	s := ""
	for i, coeff := range r.Coeffs {
		s += strconv.Itoa(i+1) + ": alt №" + strconv.Itoa(r.Order[i]) + " - " + coeff.String() + "\n"
	}
	return s
}

func (r RankedList) Value() (driver.Value, error) {
	data, err := json.Marshal(r)
	return string(data), err
}

func (r *RankedList) Scan(src interface{}) error {
	var tmp RankedList
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
	*r = tmp
	return nil
}
