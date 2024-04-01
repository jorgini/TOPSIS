package usecase

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"webApp/entity"
	"webApp/lib/eval"
	v "webApp/lib/variables"
	"webApp/repository"
	mock_repository "webApp/repository/mocks-repository"
)

func TestTaskService_UpdateTask(t *testing.T) {
	// Init test cases
	type mockBehavior func(t *mock_repository.MockTask, m *mock_repository.MockMatrix, f *mock_repository.MockIConnectionFactory,
		input *entity.TaskModel, output *entity.TaskModel)

	tests := []struct {
		name          string
		inputTask     entity.TaskModel
		outputTask    entity.TaskModel
		mockBehavior  mockBehavior
		expectedError error
	}{
		{
			name: "Ok",
			inputTask: entity.TaskModel{
				Title:        "title",
				Description:  "description",
				MaintainerID: 1,
				TaskType:     "group",
				Method:       "topsis",
				CalcSettings: 42,
				LingScale:    *eval.DefaultT1FSScale,
				Status:       entity.Draft,
			},
			outputTask: entity.TaskModel{
				Title:        "title",
				Description:  "description",
				MaintainerID: 1,
				TaskType:     "individual",
				Method:       "smart",
				CalcSettings: 42,
				LingScale:    *eval.DefaultT1FSScale,
				Status:       entity.Draft,
			},
			mockBehavior: func(t *mock_repository.MockTask, m *mock_repository.MockMatrix, f *mock_repository.MockIConnectionFactory,
				input *entity.TaskModel, output *entity.TaskModel) {
				t.EXPECT().GetTask(context.Background(), int64(1)).Return(output, nil)
				f.EXPECT().StartTransaction().Return(nil)
				t.EXPECT().UpdateTask(context.Background(), int64(1), input).Return(nil)
				f.EXPECT().Commit().Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Ok with delete dependency",
			inputTask: entity.TaskModel{
				Title:        "title",
				Description:  "description",
				MaintainerID: 1,
				TaskType:     "individual",
				Method:       "topsis",
				CalcSettings: 42,
				LingScale:    *eval.DefaultT1FSScale,
				Status:       entity.Draft,
			},
			outputTask: entity.TaskModel{
				Title:        "title",
				Description:  "description",
				MaintainerID: 1,
				TaskType:     "group",
				Method:       "smart",
				CalcSettings: 42,
				LingScale:    *eval.DefaultT1FSScale,
				Status:       entity.Draft,
			},
			mockBehavior: func(t *mock_repository.MockTask, m *mock_repository.MockMatrix, f *mock_repository.MockIConnectionFactory,
				input *entity.TaskModel, output *entity.TaskModel) {
				t.EXPECT().GetTask(context.Background(), int64(1)).Return(output, nil)
				f.EXPECT().StartTransaction().Return(nil)
				m.EXPECT().DeleteDependencies(context.Background(), int64(1), output.MaintainerID).Return(nil)
				t.EXPECT().UpdateTask(context.Background(), int64(1), input).Return(nil)
				f.EXPECT().Commit().Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Fail find task",
			mockBehavior: func(t *mock_repository.MockTask, m *mock_repository.MockMatrix, f *mock_repository.MockIConnectionFactory,
				input *entity.TaskModel, output *entity.TaskModel) {
				t.EXPECT().GetTask(context.Background(), int64(1)).Return(nil, errors.New("something went wrong"))
			},
			expectedError: errors.New("something went wrong"),
		},
		{
			name: "Fail start transaction",
			outputTask: entity.TaskModel{
				Title:        "title",
				Description:  "description",
				MaintainerID: 1,
				TaskType:     "group",
				Method:       "smart",
				CalcSettings: 42,
				LingScale:    *eval.DefaultT1FSScale,
				Status:       entity.Draft,
			},
			mockBehavior: func(t *mock_repository.MockTask, m *mock_repository.MockMatrix, f *mock_repository.MockIConnectionFactory,
				input *entity.TaskModel, output *entity.TaskModel) {
				t.EXPECT().GetTask(context.Background(), int64(1)).Return(output, nil)
				f.EXPECT().StartTransaction().Return(errors.New("something went wrong"))
			},
			expectedError: errors.New("something went wrong"),
		},
		{
			name: "Fail delete dependency",
			inputTask: entity.TaskModel{
				Title:        "title",
				Description:  "description",
				MaintainerID: 1,
				TaskType:     "individual",
				Method:       "topsis",
				CalcSettings: 42,
				LingScale:    *eval.DefaultT1FSScale,
				Status:       entity.Draft,
			},
			outputTask: entity.TaskModel{
				Title:        "title",
				Description:  "description",
				MaintainerID: 1,
				TaskType:     "group",
				Method:       "smart",
				CalcSettings: 42,
				LingScale:    *eval.DefaultT1FSScale,
				Status:       entity.Draft,
			},
			mockBehavior: func(t *mock_repository.MockTask, m *mock_repository.MockMatrix, f *mock_repository.MockIConnectionFactory,
				input *entity.TaskModel, output *entity.TaskModel) {
				t.EXPECT().GetTask(context.Background(), int64(1)).Return(output, nil)
				f.EXPECT().StartTransaction().Return(nil)
				m.EXPECT().DeleteDependencies(context.Background(), int64(1), output.MaintainerID).Return(errors.New("something went wrong"))
				f.EXPECT().Rollback().Return(nil)
			},
			expectedError: errors.New("something went wrong"),
		},
		{
			name: "Fail update task",
			inputTask: entity.TaskModel{
				Title:        "title",
				Description:  "description",
				MaintainerID: 1,
				TaskType:     "individual",
				Method:       "topsis",
				CalcSettings: 42,
				LingScale:    *eval.DefaultT1FSScale,
				Status:       entity.Draft,
			},
			outputTask: entity.TaskModel{
				Title:        "title",
				Description:  "description",
				MaintainerID: 1,
				TaskType:     "group",
				Method:       "smart",
				CalcSettings: 42,
				LingScale:    *eval.DefaultT1FSScale,
				Status:       entity.Draft,
			},
			mockBehavior: func(t *mock_repository.MockTask, m *mock_repository.MockMatrix, f *mock_repository.MockIConnectionFactory,
				input *entity.TaskModel, output *entity.TaskModel) {
				t.EXPECT().GetTask(context.Background(), int64(1)).Return(output, nil)
				f.EXPECT().StartTransaction().Return(nil)
				m.EXPECT().DeleteDependencies(context.Background(), int64(1), output.MaintainerID).Return(nil)
				t.EXPECT().UpdateTask(context.Background(), int64(1), input).Return(errors.New("something went wrong"))
				f.EXPECT().Rollback().Return(nil)
			},
			expectedError: errors.New("something went wrong"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			taskRepo := mock_repository.NewMockTask(c)
			matrixRepo := mock_repository.NewMockMatrix(c)
			factory := mock_repository.NewMockIConnectionFactory(c)

			repo := repository.Repository{Task: taskRepo, Matrix: matrixRepo, IConnectionFactory: factory}
			test.mockBehavior(taskRepo, matrixRepo, factory, &test.inputTask, &test.outputTask)

			svc := NewService(&repo)

			err := svc.UpdateTask(context.Background(), 1, &test.inputTask)

			// Assert
			if test.expectedError != nil {
				assert.Equal(t, test.expectedError.Error(), err.Error())
			} else {
				assert.Equal(t, test.expectedError, err)
			}
		})
	}
}

func TestTaskService_UpdateAlts(t *testing.T) {
	// Init test cases
	type mockBehavior func(t *mock_repository.MockTask, m *mock_repository.MockMatrix, f *mock_repository.MockIConnectionFactory,
		input entity.Alts, output entity.Alts, criteria entity.Criteria)

	tests := []struct {
		name           string
		inputAlts      entity.Alts
		outputAlts     entity.Alts
		outputCriteria entity.Criteria
		mockBehavior   mockBehavior
		expectedError  error
	}{
		{
			name: "Ok",
			inputAlts: entity.Alts{
				{"a1", "smt"},
				{"a2", "smt"},
				{"a3", "smt"},
			},
			outputAlts: entity.Alts{
				{"title1", "test"},
				{"title2", "test"},
				{"title3", "test"},
			},
			outputCriteria: entity.Criteria{
				{"c1", "smt", eval.Rating{eval.Number(0.5)}, v.Benefit},
				{"c2", "smt", eval.Rating{eval.Number(0.5)}, v.Cost},
			},
			mockBehavior: func(t *mock_repository.MockTask, m *mock_repository.MockMatrix, f *mock_repository.MockIConnectionFactory,
				input entity.Alts, output entity.Alts, criteria entity.Criteria) {
				t.EXPECT().GetCriteria(context.Background(), int64(1)).Return(criteria, nil)
				t.EXPECT().GetAlts(context.Background(), int64(1)).Return(output, nil)
				f.EXPECT().StartTransaction().Return(nil)
				t.EXPECT().UpdateAlts(context.Background(), int64(1), input).Return(nil)
				f.EXPECT().Commit().Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Ok with nullify matrix",
			inputAlts: entity.Alts{
				{"a1", "smt"},
				{"a2", "smt"},
			},
			outputAlts: entity.Alts{
				{"title1", "test"},
				{"title2", "test"},
				{"title3", "test"},
			},
			outputCriteria: entity.Criteria{
				{"c1", "smt", eval.Rating{eval.Number(0.5)}, v.Benefit},
				{"c2", "smt", eval.Rating{eval.Number(0.5)}, v.Cost},
			},
			mockBehavior: func(t *mock_repository.MockTask, m *mock_repository.MockMatrix, f *mock_repository.MockIConnectionFactory,
				input entity.Alts, output entity.Alts, criteria entity.Criteria) {
				t.EXPECT().GetCriteria(context.Background(), int64(1)).Return(criteria, nil)
				t.EXPECT().GetAlts(context.Background(), int64(1)).Return(output, nil)
				f.EXPECT().StartTransaction().Return(nil)
				t.EXPECT().UpdateAlts(context.Background(), int64(1), input).Return(nil)
				m.EXPECT().NullifyMatrices(context.Background(), int64(1), len(input), len(criteria)).Return(nil)
				f.EXPECT().Commit().Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Fail find criteria",
			mockBehavior: func(t *mock_repository.MockTask, m *mock_repository.MockMatrix, f *mock_repository.MockIConnectionFactory,
				input entity.Alts, output entity.Alts, criteria entity.Criteria) {
				t.EXPECT().GetCriteria(context.Background(), int64(1)).Return(nil, errors.New("something went wrong"))
			},
			expectedError: errors.New("something went wrong"),
		},
		{
			name: "Fail find alts",
			outputCriteria: entity.Criteria{
				{"c1", "smt", eval.Rating{eval.Number(0.5)}, v.Benefit},
				{"c2", "smt", eval.Rating{eval.Number(0.5)}, v.Cost},
			},
			mockBehavior: func(t *mock_repository.MockTask, m *mock_repository.MockMatrix, f *mock_repository.MockIConnectionFactory,
				input entity.Alts, output entity.Alts, criteria entity.Criteria) {
				t.EXPECT().GetCriteria(context.Background(), int64(1)).Return(criteria, nil)
				t.EXPECT().GetAlts(context.Background(), int64(1)).Return(nil, errors.New("something went wrong"))
			},
			expectedError: errors.New("something went wrong"),
		},
		{
			name: "Fail start transaction",
			outputAlts: entity.Alts{
				{"title1", "test"},
				{"title2", "test"},
				{"title3", "test"},
			},
			outputCriteria: entity.Criteria{
				{"c1", "smt", eval.Rating{eval.Number(0.5)}, v.Benefit},
				{"c2", "smt", eval.Rating{eval.Number(0.5)}, v.Cost},
			},
			mockBehavior: func(t *mock_repository.MockTask, m *mock_repository.MockMatrix, f *mock_repository.MockIConnectionFactory,
				input entity.Alts, output entity.Alts, criteria entity.Criteria) {
				t.EXPECT().GetCriteria(context.Background(), int64(1)).Return(criteria, nil)
				t.EXPECT().GetAlts(context.Background(), int64(1)).Return(output, nil)
				f.EXPECT().StartTransaction().Return(errors.New("something went wrong"))
			},
			expectedError: errors.New("something went wrong"),
		},
		{
			name: "Fail update alts",
			inputAlts: entity.Alts{
				{"a1", "smt"},
				{"a2", "smt"},
			},
			outputAlts: entity.Alts{
				{"title1", "test"},
				{"title2", "test"},
				{"title3", "test"},
			},
			outputCriteria: entity.Criteria{
				{"c1", "smt", eval.Rating{eval.Number(0.5)}, v.Benefit},
				{"c2", "smt", eval.Rating{eval.Number(0.5)}, v.Cost},
			},
			mockBehavior: func(t *mock_repository.MockTask, m *mock_repository.MockMatrix, f *mock_repository.MockIConnectionFactory,
				input entity.Alts, output entity.Alts, criteria entity.Criteria) {
				t.EXPECT().GetCriteria(context.Background(), int64(1)).Return(criteria, nil)
				t.EXPECT().GetAlts(context.Background(), int64(1)).Return(output, nil)
				f.EXPECT().StartTransaction().Return(nil)
				t.EXPECT().UpdateAlts(context.Background(), int64(1), input).Return(errors.New("something went wrong"))
				f.EXPECT().Rollback().Return(nil)
			},
			expectedError: errors.New("something went wrong"),
		},
		{
			name: "Fail nullify matrix",
			inputAlts: entity.Alts{
				{"a1", "smt"},
				{"a2", "smt"},
			},
			outputAlts: entity.Alts{
				{"title1", "test"},
				{"title2", "test"},
				{"title3", "test"},
			},
			outputCriteria: entity.Criteria{
				{"c1", "smt", eval.Rating{eval.Number(0.5)}, v.Benefit},
				{"c2", "smt", eval.Rating{eval.Number(0.5)}, v.Cost},
			},
			mockBehavior: func(t *mock_repository.MockTask, m *mock_repository.MockMatrix, f *mock_repository.MockIConnectionFactory,
				input entity.Alts, output entity.Alts, criteria entity.Criteria) {
				t.EXPECT().GetCriteria(context.Background(), int64(1)).Return(criteria, nil)
				t.EXPECT().GetAlts(context.Background(), int64(1)).Return(output, nil)
				f.EXPECT().StartTransaction().Return(nil)
				t.EXPECT().UpdateAlts(context.Background(), int64(1), input).Return(nil)
				m.EXPECT().NullifyMatrices(context.Background(), int64(1), len(input), len(criteria)).Return(errors.New("something went wrong"))
				f.EXPECT().Rollback().Return(nil)
			},
			expectedError: errors.New("something went wrong"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			taskRepo := mock_repository.NewMockTask(c)
			matrixRepo := mock_repository.NewMockMatrix(c)
			factory := mock_repository.NewMockIConnectionFactory(c)

			repo := repository.Repository{Task: taskRepo, Matrix: matrixRepo, IConnectionFactory: factory}
			test.mockBehavior(taskRepo, matrixRepo, factory, test.inputAlts, test.outputAlts, test.outputCriteria)

			svc := NewService(&repo)

			err := svc.UpdateAlts(context.Background(), 1, test.inputAlts)

			// Assert
			if test.expectedError != nil {
				assert.Equal(t, test.expectedError.Error(), err.Error())
			} else {
				assert.Equal(t, test.expectedError, err)
			}
		})
	}
}

func TestTaskService_UpdateCriteria(t *testing.T) {
	// Init test cases
	type mockBehavior func(t *mock_repository.MockTask, m *mock_repository.MockMatrix, f *mock_repository.MockIConnectionFactory,
		input entity.Criteria, output entity.Criteria, alts entity.Alts)

	tests := []struct {
		name           string
		inputCriteria  entity.Criteria
		outputCriteria entity.Criteria
		outputAlts     entity.Alts
		mockBehavior   mockBehavior
		expectedError  error
	}{
		{
			name: "Ok",
			inputCriteria: entity.Criteria{
				{"test", "smt", eval.Rating{eval.Number(0.4)}, v.Cost},
				{"test2", "smt", eval.Rating{eval.Number(0.6)}, v.Cost},
			},
			outputCriteria: entity.Criteria{
				{"c1", "smt", eval.Rating{eval.Number(0.5)}, v.Benefit},
				{"c2", "smt", eval.Rating{eval.Number(0.5)}, v.Cost},
			},
			outputAlts: entity.Alts{
				{"title1", "test"},
				{"title2", "test"},
				{"title3", "test"},
			},
			mockBehavior: func(t *mock_repository.MockTask, m *mock_repository.MockMatrix, f *mock_repository.MockIConnectionFactory,
				input entity.Criteria, output entity.Criteria, alts entity.Alts) {
				t.EXPECT().GetAlts(context.Background(), int64(1)).Return(alts, nil)
				t.EXPECT().GetCriteria(context.Background(), int64(1)).Return(output, nil)
				f.EXPECT().StartTransaction().Return(nil)
				t.EXPECT().UpdateCriteria(context.Background(), int64(1), input).Return(nil)
				f.EXPECT().Commit().Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Ok with nullify matrix",
			inputCriteria: entity.Criteria{
				{"c1", "smt", eval.Rating{eval.Number(0.5)}, v.Benefit},
				{"c2", "smt", eval.Rating{eval.Number(0.5)}, v.Cost},
				{"c3", "smt", eval.Rating{eval.Number(0.5)}, v.Cost},
			},
			outputCriteria: entity.Criteria{
				{"c1", "smt", eval.Rating{eval.Number(0.5)}, v.Benefit},
				{"c2", "smt", eval.Rating{eval.Number(0.5)}, v.Cost},
			},
			outputAlts: entity.Alts{
				{"title1", "test"},
				{"title2", "test"},
				{"title3", "test"},
			},
			mockBehavior: func(t *mock_repository.MockTask, m *mock_repository.MockMatrix, f *mock_repository.MockIConnectionFactory,
				input entity.Criteria, output entity.Criteria, alts entity.Alts) {
				t.EXPECT().GetAlts(context.Background(), int64(1)).Return(alts, nil)
				t.EXPECT().GetCriteria(context.Background(), int64(1)).Return(output, nil)
				f.EXPECT().StartTransaction().Return(nil)
				t.EXPECT().UpdateCriteria(context.Background(), int64(1), input).Return(nil)
				m.EXPECT().NullifyMatrices(context.Background(), int64(1), len(alts), len(input)).Return(nil)
				f.EXPECT().Commit().Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Fail find criteria",
			mockBehavior: func(t *mock_repository.MockTask, m *mock_repository.MockMatrix, f *mock_repository.MockIConnectionFactory,
				input entity.Criteria, output entity.Criteria, alts entity.Alts) {
				t.EXPECT().GetCriteria(context.Background(), int64(1)).Return(nil, errors.New("something went wrong"))
			},
			expectedError: errors.New("something went wrong"),
		},
		{
			name: "Fail find alts",
			outputCriteria: entity.Criteria{
				{"c1", "smt", eval.Rating{eval.Number(0.5)}, v.Benefit},
				{"c2", "smt", eval.Rating{eval.Number(0.5)}, v.Cost},
			},
			mockBehavior: func(t *mock_repository.MockTask, m *mock_repository.MockMatrix, f *mock_repository.MockIConnectionFactory,
				input entity.Criteria, output entity.Criteria, alts entity.Alts) {
				t.EXPECT().GetCriteria(context.Background(), int64(1)).Return(output, nil)
				t.EXPECT().GetAlts(context.Background(), int64(1)).Return(nil, errors.New("something went wrong"))
			},
			expectedError: errors.New("something went wrong"),
		},
		{
			name: "Fail start transaction",
			outputAlts: entity.Alts{
				{"title1", "test"},
				{"title2", "test"},
				{"title3", "test"},
			},
			outputCriteria: entity.Criteria{
				{"c1", "smt", eval.Rating{eval.Number(0.5)}, v.Benefit},
				{"c2", "smt", eval.Rating{eval.Number(0.5)}, v.Cost},
			},
			mockBehavior: func(t *mock_repository.MockTask, m *mock_repository.MockMatrix, f *mock_repository.MockIConnectionFactory,
				input entity.Criteria, output entity.Criteria, alts entity.Alts) {
				t.EXPECT().GetAlts(context.Background(), int64(1)).Return(alts, nil)
				t.EXPECT().GetCriteria(context.Background(), int64(1)).Return(output, nil)
				f.EXPECT().StartTransaction().Return(errors.New("something went wrong"))
			},
			expectedError: errors.New("something went wrong"),
		},
		{
			name: "Fail update alts",
			inputCriteria: entity.Criteria{
				{"c1", "smt", eval.Rating{eval.Number(0.5)}, v.Benefit},
				{"c2", "smt", eval.Rating{eval.Number(0.5)}, v.Cost},
				{"c3", "smt", eval.Rating{eval.Number(0.5)}, v.Cost},
			},
			outputCriteria: entity.Criteria{
				{"c1", "smt", eval.Rating{eval.Number(0.5)}, v.Benefit},
				{"c2", "smt", eval.Rating{eval.Number(0.5)}, v.Cost},
			},
			outputAlts: entity.Alts{
				{"title1", "test"},
				{"title2", "test"},
				{"title3", "test"},
			},
			mockBehavior: func(t *mock_repository.MockTask, m *mock_repository.MockMatrix, f *mock_repository.MockIConnectionFactory,
				input entity.Criteria, output entity.Criteria, alts entity.Alts) {
				t.EXPECT().GetAlts(context.Background(), int64(1)).Return(alts, nil)
				t.EXPECT().GetCriteria(context.Background(), int64(1)).Return(output, nil)
				f.EXPECT().StartTransaction().Return(nil)
				t.EXPECT().UpdateCriteria(context.Background(), int64(1), input).Return(errors.New("something went wrong"))
				f.EXPECT().Rollback().Return(nil)
			},
			expectedError: errors.New("something went wrong"),
		},
		{
			name: "Fail nullify matrix",
			inputCriteria: entity.Criteria{
				{"c1", "smt", eval.Rating{eval.Number(0.5)}, v.Benefit},
				{"c2", "smt", eval.Rating{eval.Number(0.5)}, v.Cost},
				{"c3", "smt", eval.Rating{eval.Number(0.5)}, v.Cost},
			},
			outputCriteria: entity.Criteria{
				{"c1", "smt", eval.Rating{eval.Number(0.5)}, v.Benefit},
				{"c2", "smt", eval.Rating{eval.Number(0.5)}, v.Cost},
			},
			outputAlts: entity.Alts{
				{"title1", "test"},
				{"title2", "test"},
				{"title3", "test"},
			},
			mockBehavior: func(t *mock_repository.MockTask, m *mock_repository.MockMatrix, f *mock_repository.MockIConnectionFactory,
				input entity.Criteria, output entity.Criteria, alts entity.Alts) {
				t.EXPECT().GetAlts(context.Background(), int64(1)).Return(alts, nil)
				t.EXPECT().GetCriteria(context.Background(), int64(1)).Return(output, nil)
				f.EXPECT().StartTransaction().Return(nil)
				t.EXPECT().UpdateCriteria(context.Background(), int64(1), input).Return(nil)
				m.EXPECT().NullifyMatrices(context.Background(), int64(1), len(alts), len(input)).Return(errors.New("something went wrong"))
				f.EXPECT().Rollback().Return(nil)
			},
			expectedError: errors.New("something went wrong"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			taskRepo := mock_repository.NewMockTask(c)
			matrixRepo := mock_repository.NewMockMatrix(c)
			factory := mock_repository.NewMockIConnectionFactory(c)

			repo := repository.Repository{Task: taskRepo, Matrix: matrixRepo, IConnectionFactory: factory}
			test.mockBehavior(taskRepo, matrixRepo, factory, test.inputCriteria, test.outputCriteria, test.outputAlts)

			svc := NewService(&repo)

			err := svc.UpdateCriteria(context.Background(), 1, test.inputCriteria)

			// Assert
			if test.expectedError != nil {
				assert.Equal(t, test.expectedError.Error(), err.Error())
			} else {
				assert.Equal(t, test.expectedError, err)
			}
		})
	}
}

func TestTaskService_SetExpertsWeights(t *testing.T) {
	// Init test cases
	type mockBehavior func(t *mock_repository.MockTask, m *mock_repository.MockMatrix,
		input entity.Weights, output []entity.ExpertStatus)

	tests := []struct {
		name          string
		inputWeights  entity.Weights
		outputExperts []entity.ExpertStatus
		mockBehavior  mockBehavior
		expectedError error
	}{
		{
			name: "Ok",
			inputWeights: entity.Weights{
				eval.Rating{eval.Number(0.5)},
				eval.Rating{eval.Number(0.5)},
				eval.Rating{eval.Number(0.5)},
			},
			outputExperts: []entity.ExpertStatus{
				{1, entity.Draft},
				{2, entity.Complete},
				{3, entity.Draft},
			},
			mockBehavior: func(t *mock_repository.MockTask, m *mock_repository.MockMatrix,
				input entity.Weights, output []entity.ExpertStatus) {
				m.EXPECT().GetExpertsRelateToTask(context.Background(), int64(1)).Return(output, nil)
				t.EXPECT().SetExpertsWeights(context.Background(), int64(1), input).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Ok with default",
			outputExperts: []entity.ExpertStatus{
				{1, entity.Draft},
				{2, entity.Complete},
			},
			mockBehavior: func(t *mock_repository.MockTask, m *mock_repository.MockMatrix,
				input entity.Weights, output []entity.ExpertStatus) {
				m.EXPECT().GetExpertsRelateToTask(context.Background(), int64(1)).Return(output, nil)
				t.EXPECT().SetExpertsWeights(context.Background(), int64(1), entity.Weights{
					eval.Rating{eval.Number(0.5)},
					eval.Rating{eval.Number(0.5)},
				}).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Fail find experts",
			mockBehavior: func(t *mock_repository.MockTask, m *mock_repository.MockMatrix,
				input entity.Weights, output []entity.ExpertStatus) {
				m.EXPECT().GetExpertsRelateToTask(context.Background(), int64(1)).Return(nil, errors.New("something went wrong"))
			},
			expectedError: errors.New("something went wrong"),
		},
		{
			name: "Wrong sizes",
			inputWeights: entity.Weights{
				eval.Rating{eval.Number(0.5)},
				eval.Rating{eval.Number(0.5)},
				eval.Rating{eval.Number(0.5)},
				eval.Rating{eval.Number(0.5)},
			},
			outputExperts: []entity.ExpertStatus{
				{1, entity.Draft},
				{2, entity.Complete},
				{3, entity.Draft},
			},
			mockBehavior: func(t *mock_repository.MockTask, m *mock_repository.MockMatrix,
				input entity.Weights, output []entity.ExpertStatus) {
				m.EXPECT().GetExpertsRelateToTask(context.Background(), int64(1)).Return(output, nil)
			},
			expectedError: errors.New("invalid size of weights for experts"),
		},
		{
			name: "Zero experts",
			inputWeights: entity.Weights{
				eval.Rating{eval.Number(0.5)},
				eval.Rating{eval.Number(0.5)},
				eval.Rating{eval.Number(0.5)},
				eval.Rating{eval.Number(0.5)},
			},
			mockBehavior: func(t *mock_repository.MockTask, m *mock_repository.MockMatrix,
				input entity.Weights, output []entity.ExpertStatus) {
				m.EXPECT().GetExpertsRelateToTask(context.Background(), int64(1)).Return(nil, nil)
			},
			expectedError: errors.New("invalid size of weights for experts"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			taskRepo := mock_repository.NewMockTask(c)
			matrixRepo := mock_repository.NewMockMatrix(c)

			repo := repository.Repository{Task: taskRepo, Matrix: matrixRepo}
			test.mockBehavior(taskRepo, matrixRepo, test.inputWeights, test.outputExperts)

			svc := NewService(&repo)

			err := svc.SetExpertsWeights(context.Background(), 1, test.inputWeights)

			// Assert
			if test.expectedError != nil {
				assert.Equal(t, test.expectedError.Error(), err.Error())
			} else {
				assert.Equal(t, test.expectedError, err)
			}
		})
	}
}
