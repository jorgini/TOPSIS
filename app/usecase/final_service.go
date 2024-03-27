package usecase

import (
	"context"
	"webApp/app/entity"
	"webApp/app/repository"
)

const (
	standardThreshold = 0.03
)

type FinalService struct {
	repo       repository.Final
	taskRepo   repository.Task
	matrixRepo repository.Matrix
}

func NewFinalService(repo repository.Final, taskRepo repository.Task, matrixRepo repository.Matrix) *FinalService {
	return &FinalService{
		repo:       repo,
		taskRepo:   taskRepo,
		matrixRepo: matrixRepo,
	}
}

func (f *FinalService) GetFinal(ctx context.Context, sid int64) (*entity.FinalModel, error) {
	return f.repo.GetFinal(ctx, sid)
}

func (f *FinalService) SetFinal(ctx context.Context, final *entity.FinalModel) error {
	return f.repo.SetFinal(ctx, final)
}

func (f *FinalService) PresentFinal(ctx context.Context, sid int64, threshold float64) (*entity.FinalModel, error) {
	task, err := f.taskRepo.GetTask(ctx, sid)
	if err != nil {
		return nil, err
	}

	prev, err := f.GetFinal(ctx, sid)
	if err == nil {
		if task.LastChangesAt.Before(prev.LastChange) && (threshold == prev.Threshold || threshold == -1) {
			return prev, nil
		}
	}

	if threshold == -1 {
		if err == nil {
			threshold = prev.Threshold
		} else {
			threshold = standardThreshold
		}
	}

	matrices, err := f.matrixRepo.GetMatricesRelateToTask(ctx, sid)
	if err != nil {
		return nil, err
	}

	result, err := entity.CalcFinal(matrices, task, threshold)
	if err != nil {
		return nil, err
	}
	return result, f.SetFinal(ctx, result)
}
