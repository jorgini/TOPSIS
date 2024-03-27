package repository

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type PgConnectionFactory struct {
	conn   Connection
	source *sqlx.DB
}

func NewPgConnectionFactory(db *sqlx.DB) *PgConnectionFactory {
	return &PgConnectionFactory{
		source: db,
	}
}

func (f *PgConnectionFactory) getConnection() Connection {
	if f.conn != nil {
		return f.conn
	}

	var err error
	f.conn, err = f.source.Connx(context.Background())
	if err != nil {
		logrus.Info(err)
		return nil
	}
	return f.conn
}

func (f *PgConnectionFactory) StartTransaction() error {
	if f.conn != nil {
		if err := f.closeConnection(); err != nil {
			return err
		}
	}

	var err error
	f.conn, err = f.source.Beginx()
	if err != nil {
		return err
	}
	return nil
}

func (f *PgConnectionFactory) Rollback() error {
	defer func() {
		f.conn = nil
	}()

	switch f.conn.(type) {
	case *sqlx.Tx:
		if err := f.conn.(*sqlx.Tx).Rollback(); err != nil {
			return err
		}
	case *sqlx.Conn:
		return errors.Join(errors.New("transaction lost"), f.closeConnection())
	default:
		return errors.New("invalid connection")
	}
	return nil
}

func (f *PgConnectionFactory) Commit() error {
	defer func() {
		f.conn = nil
	}()

	switch f.conn.(type) {
	case *sqlx.Tx:
		if err := f.conn.(*sqlx.Tx).Commit(); err != nil {
			return err
		}
	case *sqlx.Conn:
		return errors.Join(errors.New("transaction lost"), f.closeConnection())
	default:
		return errors.New("invalid connection")
	}
	return nil
}

func (f *PgConnectionFactory) closeConnection() error {
	switch f.conn.(type) {
	case *sqlx.Tx:
		return nil
	case *sqlx.Conn:
		defer func() {
			f.conn = nil
		}()

		if err := f.conn.(*sqlx.Conn).Close(); err != nil {
			return err
		}
	default:
		return errors.New("invalid connection")
	}
	return nil
}
