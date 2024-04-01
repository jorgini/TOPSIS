package repository

import (
	"context"
	"errors"
	"fmt"
	"webApp/configs"
	"webApp/entity"
)

type SessionDao struct {
	c   IConnectionFactory
	cfg *configs.DbConfig
}

func NewSessionDao(factory IConnectionFactory, config *configs.DbConfig) *SessionDao {
	return &SessionDao{
		c:   factory,
		cfg: config,
	}
}

func (s *SessionDao) GetUIDByToken(ctx context.Context, refresh string) (int64, error) {
	query := fmt.Sprintf("SELECT uid FROM %s WHERE token=$1", s.cfg.SessionTable)

	conn := s.c.GetConnection()
	if conn == nil {
		return 0, errors.New("cant connect to db")
	}

	var uid int64
	row := conn.QueryRowxContext(ctx, query, refresh)
	if err := row.Scan(&uid); err != nil {
		return 0, errors.Join(err, s.c.CloseConnection())
	}
	return uid, s.c.CloseConnection()
}

func (s *SessionDao) InsertRefreshToken(ctx context.Context, refresh entity.Session) error {
	query := fmt.Sprintf("INSERT INTO %s (uid, token, exp_at) values ($1, $2, $3)", s.cfg.SessionTable)

	conn := s.c.GetConnection()
	if conn == nil {
		return errors.New("cant connect to db")
	}

	if _, err := conn.ExecContext(ctx, query, refresh.Uid, refresh.Token, refresh.ExpiredAt); err != nil {
		return errors.Join(err, s.c.CloseConnection())
	}
	return s.c.CloseConnection()
}

func (s *SessionDao) UpdateRefreshToken(ctx context.Context, refresh entity.Session) error {
	query := fmt.Sprintf("UPDATE %s SET token=$1, exp_at=$2 WHERE uid=$3", s.cfg.SessionTable)

	conn := s.c.GetConnection()
	if conn == nil {
		return errors.New("cant connect to db")
	}

	if result, err := conn.ExecContext(ctx, query, refresh.Token, refresh.ExpiredAt, refresh.Uid); err != nil {
		return errors.Join(err, s.c.CloseConnection())
	} else if n, err := result.RowsAffected(); err != nil || n == 0 {
		return errors.Join(errors.New("nothing to update"), s.c.CloseConnection())
	}
	return s.c.CloseConnection()
}
