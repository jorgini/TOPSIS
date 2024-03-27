package usecase

import (
	"context"
	"crypto/sha1"
	"fmt"
	"time"
	"webApp/app/entity"
	"webApp/app/repository"
)

var (
	salt = "fmm39ijhn)(J@mfj0293NI)"
	ttl  = 48 * time.Hour
)

type UserService struct {
	repo       repository.User
	matrixRepo repository.Matrix
}

func NewUserService(repo repository.User, matrix repository.Matrix) *UserService {
	return &UserService{
		repo:       repo,
		matrixRepo: matrix,
	}
}

func (u *UserService) CreateNewUser(ctx context.Context, user *entity.UserModel) (int64, error) {
	user.Password = getPasswordHash(user.Password)
	return u.repo.CreateNewUser(ctx, user)
}

func (u *UserService) GetUID(ctx context.Context, email, password string) (int64, error) {
	password = getPasswordHash(password)
	return u.repo.GetUID(ctx, email, password)
}

func (u *UserService) UpdateUser(ctx context.Context, uid int64, update *entity.UserModel) error {
	return u.repo.UpdateUser(ctx, uid, update)
}

func (u *UserService) DeleteUser(ctx context.Context, uid int64) error {
	return u.repo.DeleteUser(ctx, uid)
}

func getPasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (u *UserService) GetUsersRelateToTask(ctx context.Context, sid int64) ([]string, error) {
	uids, err := u.matrixRepo.GetUIDsRelateToTask(ctx, sid)
	if err != nil {
		return nil, err
	}

	users := make([]string, len(uids))
	for i := range users {
		users[i], err = u.repo.GetUserByUID(ctx, uids[i])
		if err != nil {
			return nil, err
		}
	}
	return users, nil
}
