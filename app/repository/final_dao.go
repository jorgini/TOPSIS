package repository

import (
	"context"
	"errors"
	"fmt"
	"time"
	"webApp/app/configs"
	"webApp/app/entity"
)

type FinalDao struct {
	c   IConnectionFactory
	cfg *configs.DbConfig
}

func NewFinalDao(factory IConnectionFactory, config *configs.DbConfig) *FinalDao {
	return &FinalDao{
		c:   factory,
		cfg: config,
	}
}

func (f *FinalDao) GetFinal(ctx context.Context, sid int64) (*entity.FinalModel, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE fid=$1", f.cfg.FinalTable)

	conn := f.c.getConnection()
	if conn == nil {
		return nil, errors.New("cant connect to db")
	}

	var final entity.FinalModel
	if err := conn.GetContext(ctx, &final, query, sid); err != nil {
		return nil, errors.Join(err, f.c.closeConnection())
	}
	return &final, f.c.closeConnection()
}

func (f *FinalDao) SetFinal(ctx context.Context, final *entity.FinalModel) error {
	query := fmt.Sprintf("INSERT INTO %s (fid, result, sens_analysis, last_change) values ($1, $2, $3, $4)", f.cfg.FinalTable)

	conn := f.c.getConnection()
	if conn == nil {
		return errors.New("cant connect to db")
	}

	if _, err := conn.ExecContext(ctx, query, final.FID, final.Result, final.SensAnalysis, time.Now()); err != nil {
		return errors.Join(err, f.c.closeConnection())
	}
	return f.c.closeConnection()
}

func (f *FinalDao) UpdateFinal(ctx context.Context, final *entity.FinalModel) error {
	query := fmt.Sprintf("UPDATE %s result=$1, sens_analysis=$2, last_change=$3 WHERE fid=$4", f.cfg.FinalTable)

	conn := f.c.getConnection()
	if conn == nil {
		return errors.New("cant connect to db")
	}

	if _, err := conn.ExecContext(ctx, query, final.Result, final.SensAnalysis, time.Now(), final.FID); err != nil {
		return errors.Join(err, f.c.closeConnection())
	}
	return f.c.closeConnection()
}
