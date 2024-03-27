package repository

import (
	"context"
	"errors"
	"fmt"
	"webApp/app/configs"
	"webApp/app/entity"
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

	conn := s.c.getConnection()
	if conn == nil {
		return 0, errors.New("cant connect to db")
	}

	var uid int64
	row := conn.QueryRowxContext(ctx, query, refresh)
	if err := row.Scan(&uid); err != nil {
		return 0, errors.Join(err, s.c.closeConnection())
	}
	return uid, s.c.closeConnection()
}

func (s *SessionDao) InsertRefreshToken(ctx context.Context, refresh entity.Session) error {
	query := fmt.Sprintf("INSERT INTO %s (uid, token, exp_at) values ($1, $2, $3)", s.cfg.SessionTable)

	conn := s.c.getConnection()
	if conn == nil {
		return errors.New("cant connect to db")
	}

	if _, err := conn.ExecContext(ctx, query, refresh.Uid, refresh.Token, refresh.ExpiredAt); err != nil {
		return errors.Join(err, s.c.closeConnection())
	}
	return s.c.closeConnection()
}
