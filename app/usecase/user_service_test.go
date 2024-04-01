package usecase

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"webApp/entity"
	"webApp/repository"
	mock_repository "webApp/repository/mocks-repository"
)

func TestUserService_GetUsersRelateToTask(t *testing.T) {
	// Init test cases
	type mockBehavior func(r *mock_repository.MockUser, m *mock_repository.MockMatrix, experts []entity.ExpertStatus,
		logins []string)

	tests := []struct {
		name           string
		outputExperts  []entity.ExpertStatus
		outputLogins   []string
		mockBehavior   mockBehavior
		expectedError  error
		expectedResult []entity.Expert
	}{
		{
			name: "Ok",
			outputExperts: []entity.ExpertStatus{
				{1, entity.Draft},
				{2, entity.Complete},
				{3, entity.Draft},
			},
			outputLogins: []string{
				"user1", "user2", "user3",
			},
			mockBehavior: func(r *mock_repository.MockUser, m *mock_repository.MockMatrix, experts []entity.ExpertStatus,
				logins []string) {
				m.EXPECT().GetExpertsRelateToTask(context.Background(), int64(1)).Return(experts, nil)
				for i := range experts {
					r.EXPECT().GetUserByUID(context.Background(), experts[i].UID).Return(logins[i], nil)
				}
			},
			expectedError: nil,
			expectedResult: []entity.Expert{
				{"user1", entity.Draft},
				{"user2", entity.Complete},
				{"user3", entity.Draft},
			},
		},
		{
			name: "Wrong find relate experts",
			mockBehavior: func(r *mock_repository.MockUser, m *mock_repository.MockMatrix, experts []entity.ExpertStatus,
				logins []string) {
				m.EXPECT().GetExpertsRelateToTask(context.Background(), int64(1)).Return(nil, errors.New("something went wrong"))
			},
			expectedError:  errors.New("something went wrong"),
			expectedResult: nil,
		},
		{
			name: "Wrong find user",
			outputExperts: []entity.ExpertStatus{
				{1, entity.Draft},
				{2, entity.Complete},
				{3, entity.Draft},
			},
			mockBehavior: func(r *mock_repository.MockUser, m *mock_repository.MockMatrix, experts []entity.ExpertStatus,
				logins []string) {
				m.EXPECT().GetExpertsRelateToTask(context.Background(), int64(1)).Return(experts, nil)
				r.EXPECT().GetUserByUID(context.Background(), experts[0].UID).Return("", errors.New("something went wrong"))
			},
			expectedError:  errors.New("something went wrong"),
			expectedResult: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			userRepo := mock_repository.NewMockUser(c)
			matrixRepo := mock_repository.NewMockMatrix(c)

			repo := repository.Repository{User: userRepo, Matrix: matrixRepo}
			test.mockBehavior(userRepo, matrixRepo, test.outputExperts, test.outputLogins)

			svc := NewService(&repo)

			result, err := svc.GetUsersRelateToTask(context.Background(), 1)

			// Assert
			assert.Equal(t, test.expectedResult, result)
			assert.Equal(t, test.expectedError, err)
		})
	}
}
