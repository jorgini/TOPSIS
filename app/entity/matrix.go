package entity

import (
	"webApp/lib/matrix"
)

type MatrixModel struct {
	MID    uint64         `bson:"_id"`
	SID    uint64         `bson:"sid"`
	UID    uint64         `bson:"uid"`
	Matrix *matrix.Matrix `bson:"matrix"`
	Status bool           `bson:"status"`
}
