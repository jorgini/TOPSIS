package repository

import (
	"context"
	"errors"
	"fmt"
	"webApp/app/configs"
	"webApp/app/entity"
	"webApp/lib/eval"
	"webApp/lib/matrix"
)

type MatrixDao struct {
	c   IConnectionFactory
	cfg *configs.DbConfig
}

func NewMatrixDao(factory IConnectionFactory, config *configs.DbConfig) *MatrixDao {
	return &MatrixDao{
		c:   factory,
		cfg: config,
	}
}

func (m *MatrixDao) getSizesOfMatrix(ctx context.Context, sid int64) (int, int, error) {
	query := fmt.Sprintf("SELECT alternatives, criteria FROM %s WHERE sid=$1", m.cfg.TaskTable)

	conn := m.c.getConnection()
	if conn == nil {
		return 0, 0, errors.New("cant connect to db")
	}

	var alts entity.Alts
	var criteria entity.Criteria
	row := conn.QueryRowxContext(ctx, query, sid)
	if err := row.Scan(&alts, &criteria); err != nil {
		return 0, 0, errors.Join(err, m.c.closeConnection())
	}
	return len(alts), len(criteria), m.c.closeConnection()
}

func (m *MatrixDao) DeleteDependencies(ctx context.Context, sid, mainUid int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE sid=$1 AND uid!=$2", m.cfg.MatrixTable)

	conn := m.c.getConnection()
	if conn == nil {
		return errors.New("cant connect to db")
	}

	_, err := conn.ExecContext(ctx, query, sid, mainUid)
	if err != nil {
		return errors.Join(err, m.c.closeConnection())
	}
	return m.c.closeConnection()
}

func (m *MatrixDao) NullifyMatrices(ctx context.Context, sid int64, alts, criteria int) error {
	query := fmt.Sprintf("UPDATE %s SET matrix=$1 WHERE sid=$2", m.cfg.MatrixTable)

	conn := m.c.getConnection()
	if conn == nil {
		return errors.New("cant connect to db")
	}

	nullMatrix := matrix.NewMatrix(alts, criteria)
	if _, err := conn.ExecContext(ctx, query, nullMatrix, sid); err != nil {
		return errors.Join(err, m.c.closeConnection())
	}
	return m.c.closeConnection()
}

func (m *MatrixDao) CreateMatrix(ctx context.Context, uid, sid int64) (int64, error) {
	query := fmt.Sprintf("INSERT INTO %s (sid, uid, matrix, status) values ($1,$2,$3,$4) RETURNING mid",
		m.cfg.MatrixTable)

	x, y, err := m.getSizesOfMatrix(ctx, sid)
	if err != nil {
		return 0, err
	}
	newMatrix := matrix.NewMatrix(x, y)

	conn := m.c.getConnection()
	if conn == nil {
		return 0, errors.New("cant connect to db")
	}

	var mid int64
	row := conn.QueryRowxContext(ctx, query, sid, uid, newMatrix, entity.Draft)
	if err := row.Scan(&mid); err != nil {
		return 0, errors.Join(err, m.c.closeConnection())
	}
	return mid, m.c.closeConnection()
}

func (m *MatrixDao) CheckAccess(ctx context.Context, uid, sid int64) error {
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE uid=$1 AND sid=$2", m.cfg.MatrixTable)

	conn := m.c.getConnection()
	if conn == nil {
		return errors.New("cant connect to db")
	}

	var result int
	row := conn.QueryRowxContext(ctx, query, uid, sid)
	if err := row.Scan(&result); err != nil {
		return errors.Join(err, m.c.closeConnection())
	}
	return m.c.closeConnection()
}

func (m *MatrixDao) DeleteMatrix(ctx context.Context, uid, sid int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE uid=$1 AND sid=$2", m.cfg.MatrixTable)

	conn := m.c.getConnection()
	if conn == nil {
		return errors.New("cant connect to db")
	}

	if _, err := conn.ExecContext(ctx, query, uid, sid); err != nil {
		return errors.Join(err, m.c.closeConnection())
	}
	return m.c.closeConnection()
}

func (m *MatrixDao) GetMID(ctx context.Context, uid, sid int64) (int64, error) {
	query := fmt.Sprintf("SELECT mid FROM %s WHERE uid=$1 AND sid=$2", m.cfg.MatrixTable)

	conn := m.c.getConnection()
	if conn == nil {
		return 0, errors.New("cant connect to db")
	}

	var mid int64
	row := conn.QueryRowxContext(ctx, query, uid, sid)
	if err := row.Scan(&mid); err != nil {
		return 0, errors.Join(err, m.c.closeConnection())
	}
	return mid, m.c.closeConnection()
}

func (m *MatrixDao) GetMatrix(ctx context.Context, mid int64) (*matrix.Matrix, error) {
	query := fmt.Sprintf("SELECT matrix FROM %s WHERE mid=$1", m.cfg.MatrixTable)

	conn := m.c.getConnection()
	if conn == nil {
		return nil, errors.New("cant connect to db")
	}

	var resultMatrix matrix.Matrix
	if err := conn.GetContext(ctx, &resultMatrix, query, mid); err != nil {
		return nil, errors.Join(err, m.c.closeConnection())
	}
	return &resultMatrix, nil
}

func (m *MatrixDao) UpdateMatrix(ctx context.Context, mid, ord int64, rating []eval.Rating) error {
	prevMatrix, err := m.GetMatrix(ctx, mid)
	if err != nil {
		return err
	}

	err = prevMatrix.UpdateAlternativeRatings(int(ord), rating)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET matrix=$1 WHERE mid=$2", m.cfg.MatrixTable)

	conn := m.c.getConnection()
	if conn == nil {
		return errors.New("cant connect to db")
	}

	if _, err := conn.ExecContext(ctx, query, prevMatrix, mid); err != nil {
		return errors.Join(err, m.c.closeConnection())
	}
	return m.c.closeConnection()
}

func (m *MatrixDao) GetUIDsRelateToTask(ctx context.Context, sid int64) ([]int64, error) {
	query := fmt.Sprintf("SELECT uid FROM %s WHERE sid=$1", m.cfg.MatrixTable)

	conn := m.c.getConnection()
	if conn == nil {
		return nil, errors.New("cant connect to db")
	}

	var uids []int64
	if err := conn.SelectContext(ctx, &uids, query, sid); err != nil {
		return nil, errors.Join(err, m.c.closeConnection())
	}
	return uids, m.c.closeConnection()
}

func (m *MatrixDao) GetMatricesRelateToTask(ctx context.Context, sid int64) ([]entity.MatrixModel, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE sid=$1", m.cfg.MatrixTable)

	conn := m.c.getConnection()
	if conn == nil {
		return nil, errors.New("cant connect to db")
	}

	var matrices []entity.MatrixModel
	if err := conn.SelectContext(ctx, &matrices, query, sid); err != nil {
		return nil, err
	}
	return matrices, nil
}
