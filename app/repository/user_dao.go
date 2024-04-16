package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"webApp/configs"
	"webApp/entity"
)

type UserDao struct {
	c   IConnectionFactory
	cfg *configs.DbConfig
}

func NewUserDao(factory IConnectionFactory, config *configs.DbConfig) *UserDao {
	return &UserDao{
		c:   factory,
		cfg: config,
	}
}

func getUserUpdateQuery(user *entity.UserModel) (string, int) {
	query := make([]string, 0, 2)
	if user.Login != "" {
		query = append(query, fmt.Sprintf("login=$%d ", len(query)+1))
	}
	if user.Email != "" {
		query = append(query, fmt.Sprintf("email=$%d ", len(query)+1))
	}
	if user.Password != "" {
		query = append(query, fmt.Sprintf("password=$%d ", len(query)+1))
	}
	return strings.Join(query, ","), len(query) + 1
}

func getUserUpdateArgs(user *entity.UserModel) []interface{} {
	args := make([]interface{}, 0, 2)
	if user.Login != "" {
		args = append(args, user.Login)
	}
	if user.Email != "" {
		args = append(args, user.Email)
	}
	if user.Password != "" {
		args = append(args, user.Password)
	}
	return args
}

func logInBy(user *entity.UserModel) string {
	if user.Email != "" {
		return "email=$1"
	} else if user.Login != "" {
		return "login=$1"
	}
	return "error"
}

func logInByValue(user *entity.UserModel) interface{} {
	if user.Email != "" {
		return user.Email
	} else if user.Login != "" {
		return user.Login
	}
	return nil
}

func (u *UserDao) CreateNewUser(ctx context.Context, user *entity.UserModel) (int64, error) {
	query := fmt.Sprintf("INSERT INTO %s (login, email, password) values ($1, $2, $3) RETURNING uid", u.cfg.UserTable)

	var uid int64

	conn := u.c.GetConnection()
	if conn == nil {
		return 0, errors.New("cant connect to db")
	}

	row := conn.QueryRowxContext(ctx, query, user.Login, user.Email, user.Password)
	if err := row.Scan(&uid); err != nil {
		return 0, errors.Join(err, u.c.CloseConnection())
	}
	return uid, u.c.CloseConnection()
}

func (u *UserDao) GetUID(ctx context.Context, user *entity.UserModel) (int64, error) {
	query := fmt.Sprintf("SELECT uid FROM %s WHERE %s and password=$2", u.cfg.UserTable, logInBy(user))

	conn := u.c.GetConnection()
	if conn == nil {
		return 0, errors.New("cant connect to db")
	}

	var uid int64
	row := conn.QueryRowxContext(ctx, query, logInByValue(user), user.Password)
	if err := row.Scan(&uid); err != nil {
		return 0, errors.Join(err, u.c.CloseConnection())
	}
	return uid, u.c.CloseConnection()
}

func (u *UserDao) UpdateUser(ctx context.Context, uid int64, user *entity.UserModel) error {
	update, ord := getUserUpdateQuery(user)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE uid=$%d", u.cfg.UserTable, update, ord)
	logrus.Info(query)
	conn := u.c.GetConnection()
	if conn == nil {
		return errors.New("cant connect to db")
	}

	if result, err := conn.ExecContext(ctx, query, append(getUserUpdateArgs(user), uid)...); err != nil {
		return errors.Join(err, u.c.CloseConnection())
	} else if n, err := result.RowsAffected(); err != nil || n == 0 {
		return errors.Join(errors.New("nothing to update"), u.c.CloseConnection())
	}
	return u.c.CloseConnection()
}

func (u *UserDao) DeleteUser(ctx context.Context, uid int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE uid=$1", u.cfg.UserTable)

	conn := u.c.GetConnection()
	if conn == nil {
		return errors.New("cant connect to db")
	}

	if result, err := conn.ExecContext(ctx, query, uid); err != nil {
		return errors.Join(err, u.c.CloseConnection())
	} else if n, err := result.RowsAffected(); err != nil || n == 0 {
		return errors.Join(errors.New("nothing to update"), u.c.CloseConnection())
	}
	return u.c.CloseConnection()
}

func (u *UserDao) GetUserByUID(ctx context.Context, uid int64) (*entity.UserModel, error) {
	query := fmt.Sprintf("SELECT login, email FROM %s WHERE uid=$1", u.cfg.UserTable)

	conn := u.c.GetConnection()
	if conn == nil {
		return nil, errors.New("cant connect to db")
	}

	var user entity.UserModel
	if err := conn.GetContext(ctx, &user, query, uid); err != nil {
		return nil, err
	}
	return &user, nil
}
