package usecase

import (
	"context"
	"errors"
	"webApp/entity"
	"webApp/lib/eval"
	"webApp/repository"
)

type TaskService struct {
	repo       repository.Task
	matrixRepo repository.Matrix
	factory    repository.IConnectionFactory
}

func NewTaskService(task repository.Task, matrix repository.Matrix, factory repository.IConnectionFactory) *TaskService {
	return &TaskService{
		repo:       task,
		matrixRepo: matrix,
		factory:    factory,
	}
}

func (t *TaskService) CreateNewTask(ctx context.Context, input *entity.TaskModel) (int64, error) {
	return t.repo.CreateNewTask(ctx, input)
}

func (t *TaskService) ValidateUser(ctx context.Context, uid, sid int64) error {
	return t.repo.ValidateUser(ctx, uid, sid)
}

func (t *TaskService) CheckAccess(ctx context.Context, uid, sid int64) error {
	if err := t.repo.ValidateUser(ctx, uid, sid); err != nil {
		return t.matrixRepo.CheckAccess(ctx, uid, sid)
	}
	return nil
}

func (t *TaskService) GetTask(ctx context.Context, sid int64) (*entity.TaskModel, error) {
	task, err := t.repo.GetTask(ctx, sid)
	if err != nil {
		return nil, err
	}

	return &entity.TaskModel{
		Title:        task.Title,
		Description:  task.Description,
		TaskType:     task.TaskType,
		Method:       task.Method,
		CalcSettings: task.CalcSettings,
		LingScale:    task.LingScale,
	}, nil
}

func (t *TaskService) UpdateTask(ctx context.Context, sid int64, task *entity.TaskModel) error {
	prev, err := t.repo.GetTask(ctx, sid)
	if err != nil {
		return err
	}

	if err := t.factory.StartTransaction(); err != nil {
		return err
	}

	if prev.TaskType == entity.Group && task.TaskType == entity.Individual {
		if err := t.matrixRepo.DeleteDependencies(ctx, sid, prev.MaintainerID); err != nil {
			return errors.Join(err, t.factory.Rollback())
		}
	}

	if err := t.repo.UpdateTask(ctx, sid, task); err != nil {
		return errors.Join(err, t.factory.Rollback())
	}

	return t.factory.Commit()
}

func (t *TaskService) DeleteTask(ctx context.Context, uid, sid int64) error {
	if t.repo.ValidateUser(ctx, uid, sid) == nil {
		return t.repo.DeleteTask(ctx, sid)
	} else {
		return t.matrixRepo.DeleteMatrix(ctx, uid, sid)
	}
}

func (t *TaskService) SetPassword(ctx context.Context, sid int64, password string) error {
	return t.repo.SetPassword(ctx, sid, getPasswordHash(password))
}

func (t *TaskService) SetCriteria(ctx context.Context, sid int64, criteria entity.Criteria) error {
	return t.repo.SetCriteria(ctx, sid, criteria)
}

func (t *TaskService) SetAlts(ctx context.Context, sid int64, alts entity.Alts) error {
	return t.repo.SetAlts(ctx, sid, alts)
}

func (t *TaskService) UpdateCriteria(ctx context.Context, sid int64, newCr entity.Criteria) error {
	prevCr, err := t.GetCriteria(ctx, sid)
	if err != nil {
		return err
	}

	alts, err := t.GetAlts(ctx, sid)
	if err != nil {
		return err
	}

	if err := t.factory.StartTransaction(); err != nil {
		return err
	}

	if err := t.repo.UpdateCriteria(ctx, sid, newCr); err != nil {
		return errors.Join(err, t.factory.Rollback())
	}

	if len(prevCr) != len(newCr) {
		if err := t.matrixRepo.NullifyMatrices(ctx, sid, len(alts), len(newCr)); err != nil {
			return errors.Join(err, t.factory.Rollback())
		}
	}
	return t.factory.Commit()
}

func (t *TaskService) UpdateAlts(ctx context.Context, sid int64, newAlts entity.Alts) error {
	criteria, err := t.GetCriteria(ctx, sid)
	if err != nil {
		return err
	}

	prevAlts, err := t.GetAlts(ctx, sid)
	if err != nil {
		return err
	}

	if err := t.factory.StartTransaction(); err != nil {
		return err
	}

	if err := t.repo.UpdateAlts(ctx, sid, newAlts); err != nil {
		return errors.Join(err, t.factory.Rollback())
	}

	if len(newAlts) != len(prevAlts) {
		if err := t.matrixRepo.NullifyMatrices(ctx, sid, len(newAlts), len(criteria)); err != nil {
			return errors.Join(err, t.factory.Rollback())
		}
	}
	return t.factory.Commit()
}

func (t *TaskService) GetCriteria(ctx context.Context, sid int64) (entity.Criteria, error) {
	return t.repo.GetCriteria(ctx, sid)
}

func (t *TaskService) GetAlts(ctx context.Context, sid int64) (entity.Alts, error) {
	return t.repo.GetAlts(ctx, sid)
}

func (t *TaskService) GetAllSolutions(ctx context.Context, uid int64) ([]entity.TaskShortCard, error) {
	tasks, err := t.repo.GetAllSolutions(ctx, uid)
	if err != nil {
		return nil, err
	}

	shortcards := make([]entity.TaskShortCard, len(tasks))
	for i := range shortcards {
		shortcards[i] = entity.TaskShortCard{
			Title:       tasks[i].Title,
			Description: tasks[i].Description,
			Method:      tasks[i].Method,
			TaskType:    tasks[i].TaskType,
			LastChange:  tasks[i].LastChangesAt,
			Status:      tasks[i].Status,
		}
	}
	return shortcards, nil
}

func (t *TaskService) ConnectToTask(ctx context.Context, sid int64, password string) error {
	return t.repo.ConnectToTask(ctx, sid, getPasswordHash(password))
}

func (t *TaskService) SetExpertsWeights(ctx context.Context, sid int64, weights entity.Weights) error {
	experts, err := t.matrixRepo.GetExpertsRelateToTask(ctx, sid)
	if err != nil {
		return err
	}

	if len(weights) == 0 && len(experts) != 0 {
		weights = make([]eval.Rating, len(experts))
		for i := range weights {
			weights[i] = eval.Rating{Evaluated: eval.Number(1.0 / float64(len(experts)))}
		}
	} else if len(weights) != len(experts) || len(experts) == 0 {
		return errors.New("invalid size of weights for experts")
	}

	return t.repo.SetExpertsWeights(ctx, sid, weights)
}
