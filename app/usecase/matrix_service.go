package usecase

import (
	"context"
	"errors"
	"webApp/lib/eval"
	"webApp/repository"
)

type MatrixService struct {
	repo     repository.Matrix
	taskRepo repository.Task
	factory  repository.IConnectionFactory
}

func NewMatrixService(matrix repository.Matrix, task repository.Task, factory repository.IConnectionFactory) *MatrixService {
	return &MatrixService{
		repo:     matrix,
		taskRepo: task,
		factory:  factory,
	}
}

func (m *MatrixService) CreateMatrix(ctx context.Context, uid, sid int64) (int64, error) {
	if err := m.factory.StartTransaction(); err != nil {
		return 0, err
	}

	mid, err := m.repo.CreateMatrix(ctx, uid, sid)
	if err != nil {
		return 0, errors.Join(err, m.factory.Rollback())
	}

	if err := m.taskRepo.SetLastChange(ctx, sid); err != nil {
		return 0, errors.Join(err, m.factory.Rollback())
	}
	return mid, m.factory.Commit()
}

func (m *MatrixService) GetMID(ctx context.Context, uid, sid int64) (int64, error) {
	return m.repo.GetMID(ctx, uid, sid)
}

func (m *MatrixService) UpdateMatrix(ctx context.Context, sid, mid, ord int64, rating []eval.Rating) error {
	if err := m.factory.StartTransaction(); err != nil {
		return err
	}

	if err := m.repo.UpdateMatrix(ctx, mid, ord, rating); err != nil {
		return errors.Join(err, m.factory.Rollback())
	}

	if err := m.taskRepo.SetLastChange(ctx, sid); err != nil {
		return errors.Join(err, m.factory.Rollback())
	}
	return m.factory.Commit()
}

func (m *MatrixService) GetRatings(ctx context.Context, uid, sid, ord int64) ([]eval.Rating, error) {
	mid, err := m.GetMID(ctx, uid, sid)
	if err != nil {
		return nil, err
	}

	matrix, err := m.repo.GetMatrix(ctx, mid)
	if err != nil {
		return nil, err
	}

	return matrix.GetAlternativeRatings(int(ord))
}

func (m *MatrixService) GetUIDsRelateToTask(ctx context.Context, sid int64) ([]int64, error) {
	return m.repo.GetUIDsRelateToTask(ctx, sid)
}
