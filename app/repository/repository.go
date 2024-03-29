package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"webApp/configs"
	"webApp/entity"
	"webApp/lib/eval"
	"webApp/lib/matrix"
)

type User interface {
	CreateNewUser(ctx context.Context, user *entity.UserModel) (int64, error)
	GetUID(ctx context.Context, user *entity.UserModel) (int64, error)
	UpdateUser(ctx context.Context, uid int64, update *entity.UserModel) error
	DeleteUser(ctx context.Context, uid int64) error
	GetUserByUID(ctx context.Context, uid int64) (string, error)
}

type Session interface {
	GetUIDByToken(ctx context.Context, refresh string) (int64, error)
	InsertRefreshToken(ctx context.Context, refresh entity.Session) error
}

type Task interface {
	CreateNewTask(ctx context.Context, task *entity.TaskModel) (int64, error)
	ValidateUser(ctx context.Context, uid, sid int64) error
	UpdateTask(ctx context.Context, sid int64, input *entity.TaskModel) error
	GetTask(ctx context.Context, sid int64) (*entity.TaskModel, error)
	DeleteTask(ctx context.Context, sid int64) error
	SetPassword(ctx context.Context, sid int64, password string) error
	SetCriteria(ctx context.Context, sid int64, alts entity.Criteria) error
	SetAlts(ctx context.Context, sid int64, alts entity.Alts) error
	UpdateCriteria(ctx context.Context, sid int64, criteria entity.Criteria) error
	UpdateAlts(ctx context.Context, sid int64, alts entity.Alts) error
	GetCriteria(ctx context.Context, sid int64) (entity.Criteria, error)
	GetAlts(ctx context.Context, sid int64) (entity.Alts, error)
	GetAllSolutions(ctx context.Context, uid int64) ([]entity.TaskModel, error)
	ConnectToTask(ctx context.Context, sid int64, password string) error
	SetLastChange(ctx context.Context, sid int64) error
	SetExpertsWeights(ctx context.Context, sid int64, weights entity.Weights) error
}

type Matrix interface {
	CreateMatrix(ctx context.Context, uid, sid int64) (int64, error)
	DeleteDependencies(ctx context.Context, sid, mainUid int64) error
	NullifyMatrices(ctx context.Context, sid int64, alts, criteria int) error
	CheckAccess(ctx context.Context, uid, sid int64) error
	DeleteMatrix(ctx context.Context, uid, sid int64) error
	GetMID(ctx context.Context, uid, sid int64) (int64, error)
	UpdateMatrix(ctx context.Context, mid, ord int64, rating []eval.Rating) error
	GetMatrix(ctx context.Context, mid int64) (*matrix.Matrix, error)
	GetExpertsRelateToTask(ctx context.Context, sid int64) ([]entity.ExpertStatus, error)
	GetMatricesRelateToTask(ctx context.Context, sid int64) ([]entity.MatrixModel, error)
	SetStatusComplete(ctx context.Context, mid int64) error
	DeactivateStatuses(ctx context.Context, sid int64) error
}

type Final interface {
	GetFinal(ctx context.Context, sid int64) (*entity.FinalModel, error)
	SetFinal(ctx context.Context, final *entity.FinalModel) error
	UpdateFinal(ctx context.Context, final *entity.FinalModel) error
}

type Connection interface {
	sqlx.ExecerContext
	sqlx.QueryerContext
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type IConnectionFactory interface {
	StartTransaction() error
	Rollback() error
	Commit() error
	getConnection() Connection
	closeConnection() error
}

type Repository struct {
	User
	Session
	Task
	Matrix
	Final
	IConnectionFactory
}

func NewRepository(db *sqlx.DB, config *configs.DbConfig) *Repository {
	factory := NewPgConnectionFactory(db)
	return &Repository{
		IConnectionFactory: factory,
		User:               NewUserDao(factory, config),
		Session:            NewSessionDao(factory, config),
		Task:               NewTaskDao(factory, config),
		Matrix:             NewMatrixDao(factory, config),
		Final:              NewFinalDao(factory, config),
	}
}
