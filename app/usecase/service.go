package usecase

import (
	"context"
	"webApp/configs"
	"webApp/entity"
	"webApp/lib/eval"
	"webApp/repository"
)

//go:generate mockgen -source=service.go -destination=mocks-service/mock.go

type User interface {
	CreateNewUser(ctx context.Context, user *entity.UserModel) (int64, error)
	GetUID(ctx context.Context, user *entity.UserModel) (int64, error)
	UpdateUser(ctx context.Context, uid int64, update *entity.UserModel) error
	DeleteUser(ctx context.Context, uid int64) error
	GetUsersRelateToTask(ctx context.Context, sid int64) ([]entity.Expert, error)
}

type Session interface {
	GenerateToken(ctx context.Context, uid int64, cfg *configs.AppConfig) (entity.Tokens, error)
	RefreshToken(ctx context.Context, refresh string, cfg *configs.AppConfig) (entity.Tokens, error)
}

type Task interface {
	CreateNewTask(ctx context.Context, task *entity.TaskModel) (int64, error)
	ValidateUser(ctx context.Context, uid, sid int64) error
	CheckAccess(ctx context.Context, uid, sid int64) error
	UpdateTask(ctx context.Context, sid int64, input *entity.TaskModel) error
	GetTask(ctx context.Context, sid int64) (*entity.TaskModel, error)
	DeleteTask(ctx context.Context, uid, sid int64) error
	SetPassword(ctx context.Context, sid int64, password string) error
	SetCriteria(ctx context.Context, sid int64, criteria entity.Criteria) error
	SetAlts(ctx context.Context, sid int64, alts entity.Alts) error
	UpdateCriteria(ctx context.Context, sid int64, criteria entity.Criteria) error
	UpdateAlts(ctx context.Context, sid int64, alts entity.Alts) error
	GetCriteria(ctx context.Context, sid int64) (entity.Criteria, error)
	GetAlts(ctx context.Context, sid int64) (entity.Alts, error)
	GetAllSolutions(ctx context.Context, uid int64) ([]entity.TaskShortCard, error)
	ConnectToTask(ctx context.Context, sid int64, password string) error
	SetExpertsWeights(ctx context.Context, sid int64, weights entity.Weights) error
}

type Matrix interface {
	CreateMatrix(ctx context.Context, uid, sid int64) (int64, error)
	GetMID(ctx context.Context, uid, sid int64) (int64, error)
	UpdateMatrix(ctx context.Context, sid, mid, ord int64, rating []eval.Rating) error
	GetRatings(ctx context.Context, uid, sid, ord int64) ([]eval.Rating, error)
	GetExpertsRelateToTask(ctx context.Context, sid int64) ([]entity.ExpertStatus, error)
	SetStatusComplete(ctx context.Context, mid int64) error
	DeactivateStatuses(ctx context.Context, sid int64) error
	IsAllStatusesComplete(ctx context.Context, sid int64) error
}

type Final interface {
	PresentFinal(ctx context.Context, sid int64, threshold float64) (*entity.FinalModel, error)
	GetFinal(ctx context.Context, sid int64) (*entity.FinalModel, error)
}

type DiService interface {
	GetInstanceService() *Service
}

type Service struct {
	User
	Session
	Task
	Matrix
	Final
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User:    NewUserService(repo.User, repo.Matrix),
		Session: NewSessionService(repo.Session),
		Task:    NewTaskService(repo.Task, repo.Matrix, repo.IConnectionFactory),
		Matrix:  NewMatrixService(repo.Matrix, repo.Task, repo.IConnectionFactory),
		Final:   NewFinalService(repo.Final, repo.Task, repo.Matrix),
	}
}
