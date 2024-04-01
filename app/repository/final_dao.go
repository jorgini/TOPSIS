package repository

import (
	"context"
	"errors"
	"fmt"
	"time"
	"webApp/configs"
	"webApp/entity"
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

	conn := f.c.GetConnection()
	if conn == nil {
		return nil, errors.New("cant connect to db")
	}

	var final entity.FinalModel
	if err := conn.GetContext(ctx, &final, query, sid); err != nil {
		return nil, errors.Join(err, f.c.CloseConnection())
	}
	return &final, f.c.CloseConnection()
}

func (f *FinalDao) SetFinal(ctx context.Context, final *entity.FinalModel) error {
	query := fmt.Sprintf("INSERT INTO %s (fid, result, sens_analysis, threshold, last_change) values ($1, $2, $3, $4, $5)",
		f.cfg.FinalTable)

	conn := f.c.GetConnection()
	if conn == nil {
		return errors.New("cant connect to db")
	}

	if _, err := conn.ExecContext(ctx, query, final.FID, final.Result, final.SensAnalysis, final.Threshold, time.Now()); err != nil {
		return errors.Join(err, f.c.CloseConnection())
	}
	return f.c.CloseConnection()
}

func (f *FinalDao) UpdateFinal(ctx context.Context, final *entity.FinalModel) error {
	query := fmt.Sprintf("UPDATE %s SET result=$1, sens_analysis=$2, threshold=$3, last_change=$4 WHERE fid=$5", f.cfg.FinalTable)

	conn := f.c.GetConnection()
	if conn == nil {
		return errors.New("cant connect to db")
	}

	if result, err := conn.ExecContext(ctx, query, final.Result, final.SensAnalysis, final.Threshold, time.Now(), final.FID); err != nil {
		return errors.Join(err, f.c.CloseConnection())
	} else if n, err := result.RowsAffected(); err != nil || n == 0 {
		return errors.Join(errors.New("nothing to update"), f.c.CloseConnection())
	}
	return f.c.CloseConnection()
}
