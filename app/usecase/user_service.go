package usecase

import (
	"context"
	"crypto/sha1"
	"fmt"
	"webApp/entity"
	"webApp/lib/eval"
	"webApp/repository"
)

var (
	salt = "fmm39ijhn)(J@mfj0293NI)"
)

type UserService struct {
	repo       repository.User
	taskRepo   repository.Task
	matrixRepo repository.Matrix
}

func NewUserService(repo repository.User, task repository.Task, matrix repository.Matrix) *UserService {
	return &UserService{
		repo:       repo,
		taskRepo:   task,
		matrixRepo: matrix,
	}
}

func (u *UserService) CreateNewUser(ctx context.Context, user *entity.UserModel) (int64, error) {
	user.Password = getPasswordHash(user.Password)
	return u.repo.CreateNewUser(ctx, user)
}

func (u *UserService) GetUID(ctx context.Context, user *entity.UserModel) (int64, error) {
	user.Password = getPasswordHash(user.Password)
	return u.repo.GetUID(ctx, user)
}

func (u *UserService) UpdateUser(ctx context.Context, uid int64, update *entity.UserModel) error {
	return u.repo.UpdateUser(ctx, uid, update)
}

func (u *UserService) DeleteUser(ctx context.Context, uid int64) error {
	return u.repo.DeleteUser(ctx, uid)
}

func (u *UserService) GetUsersRelateToTask(ctx context.Context, sid int64) ([]entity.Expert, error) {
	uids, err := u.matrixRepo.GetExpertsRelateToTask(ctx, sid)
	if err != nil {
		return nil, err
	}

	task, err := u.taskRepo.GetTask(ctx, sid)
	if err != nil {
		return nil, err
	}

	experts := make([]entity.Expert, len(uids))
	for i := range experts {
		tmp, err := u.repo.GetUserByUID(ctx, uids[i].UID)
		if err != nil {
			return nil, err
		}
		experts[i].Login = tmp.Login
		experts[i].Status = uids[i].Status
		if task.ExpertsWeights != nil {
			if i >= len(task.ExpertsWeights) {
				experts[i].Weight = eval.Rating{eval.Number(0.0)}
			} else {
				experts[i].Weight = task.ExpertsWeights[i]
			}
		}
	}
	return experts, nil
}

func (u *UserService) GetUserInfo(ctx context.Context, uid int64) (*entity.UserModel, error) {
	return u.repo.GetUserByUID(ctx, uid)
}

func getPasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
