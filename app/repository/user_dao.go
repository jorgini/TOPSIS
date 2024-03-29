package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"webApp/app/configs"
	"webApp/app/entity"
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
	if user.Email != "" {
		args = append(args, user.Email)
	}
	if user.Password != "" {
		args = append(args, user.Password)
	}
	return args
}

func (u *UserDao) CreateNewUser(ctx context.Context, user *entity.UserModel) (int64, error) {
	query := fmt.Sprintf("INSERT INTO %s (email, password) values ($1, $2) RETURNING uid", u.cfg.UserTable)

	var uid int64

	conn := u.c.getConnection()
	if conn == nil {
		return 0, errors.New("cant connect to db")
	}

	row := conn.QueryRowxContext(ctx, query, user.Email, user.Password)
	if err := row.Scan(&uid); err != nil {
		return 0, errors.Join(err, u.c.closeConnection())
	}
	return uid, u.c.closeConnection()
}

func (u *UserDao) GetUID(ctx context.Context, email, password string) (int64, error) {
	query := fmt.Sprintf("SELECT uid FROM %s WHERE email=$1 and password=$2", u.cfg.UserTable)

	conn := u.c.getConnection()
	if conn == nil {
		return 0, errors.New("cant connect to db")
	}

	var uid int64
	row := conn.QueryRowxContext(ctx, query, email, password)
	if err := row.Scan(&uid); err != nil {
		return 0, errors.Join(err, u.c.closeConnection())
	}
	return uid, u.c.closeConnection()
}

func (u *UserDao) UpdateUser(ctx context.Context, uid int64, user *entity.UserModel) error {
	update, ord := getUserUpdateQuery(user)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE uid=$%d", u.cfg.UserTable, update, ord)
	logrus.Info(query)
	conn := u.c.getConnection()
	if conn == nil {
		return errors.New("cant connect to db")
	}

	_, err := conn.ExecContext(ctx, query, append(getUserUpdateArgs(user), uid)...)
	if err != nil {
		return errors.Join(err, u.c.closeConnection())
	}
	return u.c.closeConnection()
}

func (u *UserDao) DeleteUser(ctx context.Context, uid int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE uid=$1", u.cfg.UserTable)

	conn := u.c.getConnection()
	if conn == nil {
		return errors.New("cant connect to db")
	}

	_, err := conn.ExecContext(ctx, query, uid)
	if err != nil {
		return errors.Join(err, u.c.closeConnection())
	}
	return u.c.closeConnection()
}

func (u *UserDao) GetUserByUID(ctx context.Context, uid int64) (string, error) {
	query := fmt.Sprintf("SELECT email FROM %s WHERE uid=$1", u.cfg.UserTable)

	conn := u.c.getConnection()
	if conn == nil {
		return "", errors.New("cant connect to db")
	}

	var email string
	if err := conn.GetContext(ctx, &email, query, uid); err != nil {
		return "", err
	}
	return email, nil
}
