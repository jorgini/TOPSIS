package entity

import "time"

type Session struct {
	Uid       int64     `db:"uid"`
	Token     string    `db:"token"`
	ExpiredAt time.Time `db:"exp_at"`
}

type Tokens struct {
	Access  string
	Refresh string
}
