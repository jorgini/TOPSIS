package entity

import (
	"encoding/json"
	"errors"
)

type UserModel struct {
	UID      int64  `json:"uid" db:"uid"`
	Login    string `json:"login" db:"login"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

func (u *UserModel) UnmarshalJSON(data []byte) error {
	result := struct {
		UID      int64   `json:"uid"`
		Login    *string `json:"login"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
	}{}

	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.Email == nil || result.Password == nil || result.Login == nil {
		return errors.New("invalid input data for sign up")
	} else {
		u.UID = result.UID
		u.Login = *result.Login
		u.Email = *result.Email
		u.Password = *result.Password
	}
	return nil
}

type Expert struct {
	Login  string `json:"login"`
	Status bool   `json:"status"`
}

type ExpertStatus struct {
	UID    int64 `db:"uid"`
	Status bool  `db:"status"`
}
