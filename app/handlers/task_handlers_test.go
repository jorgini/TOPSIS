package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http/httptest"
	"testing"
	"webApp/entity"
	"webApp/lib/eval"
	v "webApp/lib/variables"
	"webApp/usecase"
	mock_service "webApp/usecase/mocks"
)

func TestHandler_CreateSolution(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockTask, di *mock_service.MockDiService, task *entity.TaskModel, svc *usecase.Service)
	tests := []struct {
		name                 string
		userIdentify         func(c *fiber.Ctx) (int64, error)
		inputBody            string
		inputTask            entity.TaskModel
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody: `{"title": "title", "description": "description", "task_type": "individuals", "method": "topsis", "calc_settings": 42}`,
			inputTask: entity.TaskModel{
				Title:        "title",
				Description:  "description",
				MaintainerID: 1,
				TaskType:     "individuals",
				Method:       "topsis",
				CalcSettings: 42,
				LingScale:    *eval.DefaultT1FSScale,
				Status:       entity.Draft,
			},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, task *entity.TaskModel, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().CreateNewTask(context.Background(), task).Return(int64(1), nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"success"}`,
		},
		{
			name: "Partition request",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody: `{"title": "title", "task_type": "individuals", "method": "topsis", "calc_settings": 42}`,
			inputTask: entity.TaskModel{
				Title:        "title",
				MaintainerID: 1,
				TaskType:     "individuals",
				Method:       "topsis",
				CalcSettings: 42,
				LingScale:    *eval.DefaultT1FSScale,
				Status:       entity.Draft,
			},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, task *entity.TaskModel, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().CreateNewTask(context.Background(), task).Return(int64(1), nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"success"}`,
		},
		{
			name: "Wrong Input",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody: `{"title": "title", "task_type": "individuals", "method": "topsisss", "calc_settings": 42}`,
			inputTask: entity.TaskModel{},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, task *entity.TaskModel, svc *usecase.Service) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input arguments for task, check required fields"}`,
		},
		{
			name: "Server error",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 0, errors.New("cant read id")
			},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, task *entity.TaskModel, svc *usecase.Service) {
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"cant read id"}`,
		},
		{
			name: "Service Error",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody: `{"title": "title", "task_type": "individuals", "method": "topsis", "calc_settings": 42}`,
			inputTask: entity.TaskModel{
				Title:        "title",
				MaintainerID: 1,
				TaskType:     "individuals",
				Method:       "topsis",
				CalcSettings: 42,
				LingScale:    *eval.DefaultT1FSScale,
				Status:       entity.Draft,
			},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, task *entity.TaskModel, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().CreateNewTask(context.Background(), task).Return(int64(0), errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockTask(c)
			di := mock_service.NewMockDiService(c)

			svc := usecase.Service{Task: repo}
			test.mockBehavior(repo, di, &test.inputTask, &svc)

			handler := Handler{di: di, userIdentity: test.userIdentify}

			r := fiber.New()
			r.Post("/solution", handler.CreateSolution)

			// Create Request
			req := httptest.NewRequest("POST", "/solution",
				bytes.NewBufferString(test.inputBody))
			req.Header.Set("Content-Type", "application/json")

			// Make Request
			w, err := r.Test(req)
			if err != nil {
				t.Error(err)
			}

			// Read response
			resp, err := ReadResponse(w.Body)
			if err != nil {
				t.Fatal(err)
			}

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.StatusCode)
			assert.Equal(t, test.expectedResponseBody, resp)
		})
	}
}

func TestHandler_UpdateSolution(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockTask, di *mock_service.MockDiService, task *entity.TaskModel, svc *usecase.Service)
	tests := []struct {
		name                 string
		paramsName           string
		userIdentify         func(c *fiber.Ctx) (int64, error)
		inputBody            string
		inputTask            entity.TaskModel
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:       "Ok",
			paramsName: "sid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody: `{"title": "title", "description": "description", "task_type": "individuals", "method": "topsis", "calc_settings": 42}`,
			inputTask: entity.TaskModel{
				Title:        "title",
				Description:  "description",
				MaintainerID: 1,
				TaskType:     "individuals",
				Method:       "topsis",
				CalcSettings: 42,
				LingScale:    *eval.DefaultT1FSScale,
				Status:       entity.Draft,
			},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, task *entity.TaskModel, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().ValidateUser(context.Background(), int64(1), int64(1)).Return(nil)
				r.EXPECT().UpdateTask(context.Background(), int64(1), task).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"success"}`,
		},
		{
			name:       "Partition request",
			paramsName: "sid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody: `{"title": "title", "task_type": "individuals", "method": "topsis", "calc_settings": 42}`,
			inputTask: entity.TaskModel{
				Title:        "title",
				MaintainerID: 1,
				TaskType:     "individuals",
				Method:       "topsis",
				CalcSettings: 42,
				LingScale:    *eval.DefaultT1FSScale,
				Status:       entity.Draft,
			},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, task *entity.TaskModel, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().ValidateUser(context.Background(), int64(1), int64(1)).Return(nil)
				r.EXPECT().UpdateTask(context.Background(), int64(1), task).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"success"}`,
		},
		{
			name:       "Wrong Input",
			paramsName: "sid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody: `{"title": "title", "task_type": "individuals", "method": "topsisss", "calc_settings": 42}`,
			inputTask: entity.TaskModel{},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, task *entity.TaskModel, svc *usecase.Service) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input arguments for task, check required fields"}`,
		},
		{
			name:       "Wrong URL params",
			paramsName: "smt",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody: `{"title": "title", "task_type": "individuals", "method": "topsisss", "calc_settings": 42}`,
			inputTask: entity.TaskModel{},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, task *entity.TaskModel, svc *usecase.Service) {
			},
			expectedStatusCode:   404,
			expectedResponseBody: `{"message":"task doesn't specified"}`,
		},
		{
			name:       "Server error",
			paramsName: "sid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 0, errors.New("cant read id")
			},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, task *entity.TaskModel, svc *usecase.Service) {
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"cant read id"}`,
		},
		{
			name:       "Validate error",
			paramsName: "sid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody: `{"title": "title", "task_type": "individuals", "method": "topsis", "calc_settings": 42}`,
			inputTask: entity.TaskModel{},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, task *entity.TaskModel, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().ValidateUser(context.Background(), int64(1), int64(1)).Return(errors.New("forbidden"))
			},
			expectedStatusCode:   403,
			expectedResponseBody: `{"message":"forbidden"}`,
		},
		{
			name:       "Service Error",
			paramsName: "sid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody: `{"title": "title", "task_type": "individuals", "method": "topsis", "calc_settings": 42}`,
			inputTask: entity.TaskModel{
				Title:        "title",
				MaintainerID: 1,
				TaskType:     "individuals",
				Method:       "topsis",
				CalcSettings: 42,
				LingScale:    *eval.DefaultT1FSScale,
				Status:       entity.Draft,
			},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, task *entity.TaskModel, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().ValidateUser(context.Background(), int64(1), int64(1)).Return(nil)
				r.EXPECT().UpdateTask(context.Background(), int64(1), task).Return(errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockTask(c)
			di := mock_service.NewMockDiService(c)

			svc := usecase.Service{Task: repo}
			test.mockBehavior(repo, di, &test.inputTask, &svc)

			handler := Handler{di: di, userIdentity: test.userIdentify}

			r := fiber.New()
			r.Put("/solution", handler.UpdateSolution)

			// Create Request
			req := httptest.NewRequest("PUT", "/solution?"+test.paramsName+"=1",
				bytes.NewBufferString(test.inputBody))
			req.Header.Set("Content-Type", "application/json")

			// Make Request
			w, err := r.Test(req)
			if err != nil {
				t.Error(err)
			}

			// Read response
			resp, err := ReadResponse(w.Body)
			if err != nil {
				t.Fatal(err)
			}

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.StatusCode)
			assert.Equal(t, test.expectedResponseBody, resp)
		})
	}
}

func TestHandler_GetSolution(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockTask, di *mock_service.MockDiService, task *entity.TaskModel, svc *usecase.Service)

	tests := []struct {
		name                 string
		paramsName           string
		userIdentify         func(c *fiber.Ctx) (int64, error)
		outputTask           entity.TaskModel
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:       "Ok",
			paramsName: "sid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			outputTask: entity.TaskModel{
				Title:        "title",
				Description:  "description",
				TaskType:     "individuals",
				Method:       "topsis",
				CalcSettings: 42,
				LingScale:    *eval.DefaultT1FSScale,
			},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, task *entity.TaskModel, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().CheckAccess(context.Background(), int64(1), int64(1)).Return(nil)
				r.EXPECT().GetTask(context.Background(), int64(1)).Return(task, nil)
			},
			expectedStatusCode: 200,
		},
		{
			name:       "Wrong URL params",
			paramsName: "smt",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			outputTask: entity.TaskModel{},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, task *entity.TaskModel, svc *usecase.Service) {
			},
			expectedStatusCode:   404,
			expectedResponseBody: `{"message":"task doesn't specified"}`,
		},
		{
			name:       "Server error",
			paramsName: "sid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 0, errors.New("cant read id")
			},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, task *entity.TaskModel, svc *usecase.Service) {
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"cant read id"}`,
		},
		{
			name:       "Check access error",
			paramsName: "sid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			outputTask: entity.TaskModel{},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, task *entity.TaskModel, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().CheckAccess(context.Background(), int64(1), int64(1)).Return(errors.New("forbidden"))
			},
			expectedStatusCode:   403,
			expectedResponseBody: `{"message":"hasn't access to solution"}`,
		},
		{
			name:       "Service Error",
			paramsName: "sid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			outputTask: entity.TaskModel{},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, task *entity.TaskModel, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().CheckAccess(context.Background(), int64(1), int64(1)).Return(nil)
				r.EXPECT().GetTask(context.Background(), int64(1)).Return(nil, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockTask(c)
			di := mock_service.NewMockDiService(c)

			svc := usecase.Service{Task: repo}
			test.mockBehavior(repo, di, &test.outputTask, &svc)

			handler := Handler{di: di, userIdentity: test.userIdentify}

			r := fiber.New()
			r.Get("/solution", handler.GetSolution)

			// Create Request
			req := httptest.NewRequest("GET", "/solution?"+test.paramsName+"=1",
				bytes.NewBufferString(""))
			req.Header.Set("Content-Type", "application/json")

			// Make Request
			w, err := r.Test(req)
			if err != nil {
				t.Error(err)
			}

			// Read response
			resp, err := ReadResponse(w.Body)
			if err != nil {
				t.Fatal(err)
			}

			if test.expectedResponseBody == "" {
				data, err := json.Marshal(TaskInput{
					Title:        test.outputTask.Title,
					Description:  test.outputTask.Description,
					TaskType:     test.outputTask.TaskType,
					Method:       test.outputTask.Method,
					CalcSettings: test.outputTask.CalcSettings,
					LingScale:    test.outputTask.LingScale,
				})
				if err != nil {
					t.Fatal(err)
				}
				test.expectedResponseBody = string(data)
			}

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.StatusCode)
			assert.Equal(t, test.expectedResponseBody, resp)
		})
	}
}

func TestHandler_DeleteSolution(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockTask, di *mock_service.MockDiService, svc *usecase.Service)
	tests := []struct {
		name                 string
		userIdentify         func(c *fiber.Ctx) (int64, error)
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody: `{"sid": 1}`,
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().CheckAccess(context.Background(), int64(1), int64(1)).Return(nil)
				r.EXPECT().DeleteTask(context.Background(), int64(1), int64(1)).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"success"}`,
		},
		{
			name: "Wrong Input",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody: `{}`,
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, svc *usecase.Service) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"cant find task identification"}`,
		},
		{
			name: "Server error",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 0, errors.New("cant read id")
			},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, svc *usecase.Service) {
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"cant read id"}`,
		},
		{
			name: "Check access error",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody: `{"sid": 1}`,
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().CheckAccess(context.Background(), int64(1), int64(1)).Return(errors.New("forbidden"))
			},
			expectedStatusCode:   403,
			expectedResponseBody: `{"message":"hasn't access to solution"}`,
		},
		{
			name: "Service Error",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody: `{"sid": 1}`,
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().CheckAccess(context.Background(), int64(1), int64(1)).Return(nil)
				r.EXPECT().DeleteTask(context.Background(), int64(1), int64(1)).Return(errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockTask(c)
			di := mock_service.NewMockDiService(c)

			svc := usecase.Service{Task: repo}
			test.mockBehavior(repo, di, &svc)

			handler := Handler{di: di, userIdentity: test.userIdentify}

			r := fiber.New()
			r.Delete("/solution", handler.DeleteSolution)

			// Create Request
			req := httptest.NewRequest("DELETE", "/solution",
				bytes.NewBufferString(test.inputBody))
			req.Header.Set("Content-Type", "application/json")

			// Make Request
			w, err := r.Test(req)
			if err != nil {
				t.Error(err)
			}

			// Read response
			resp, err := ReadResponse(w.Body)
			if err != nil {
				t.Fatal(err)
			}

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.StatusCode)
			assert.Equal(t, test.expectedResponseBody, resp)
		})
	}
}

func TestHandler_SetPassword(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockTask, di *mock_service.MockDiService, pass string, svc *usecase.Service)
	tests := []struct {
		name                 string
		paramsName           string
		userIdentify         func(c *fiber.Ctx) (int64, error)
		inputBody            string
		inputPass            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:       "Ok",
			paramsName: "sid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody: `{"password": "password"}`,
			inputPass: "password",
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, pass string, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().ValidateUser(context.Background(), int64(1), int64(1)).Return(nil)
				r.EXPECT().SetPassword(context.Background(), int64(1), pass).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"success"}`,
		},
		{
			name:       "Wrong Input",
			paramsName: "sid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody: ``,
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, pass string, svc *usecase.Service) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"password doesn't specified"}`,
		},
		{
			name:       "Wrong URL params",
			paramsName: "smt",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody: `{"passssword": "qwerty"}`,
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, pass string, svc *usecase.Service) {
			},
			expectedStatusCode:   404,
			expectedResponseBody: `{"message":"task doesn't specified"}`,
		},
		{
			name:       "Server error",
			paramsName: "sid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 0, errors.New("cant read id")
			},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, pass string, svc *usecase.Service) {
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"cant read id"}`,
		},
		{
			name:       "Validate error",
			paramsName: "sid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody: `{"password": "password"}`,
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, pass string, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().ValidateUser(context.Background(), int64(1), int64(1)).Return(errors.New("forbidden"))
			},
			expectedStatusCode:   403,
			expectedResponseBody: `{"message":"forbidden"}`,
		},
		{
			name:       "Service Error",
			paramsName: "sid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody: `{"password": "password"}`,
			inputPass: "password",
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, pass string, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().ValidateUser(context.Background(), int64(1), int64(1)).Return(nil)
				r.EXPECT().SetPassword(context.Background(), int64(1), pass).Return(errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockTask(c)
			di := mock_service.NewMockDiService(c)

			svc := usecase.Service{Task: repo}
			test.mockBehavior(repo, di, test.inputPass, &svc)

			handler := Handler{di: di, userIdentity: test.userIdentify}

			r := fiber.New()
			r.Patch("/solution", handler.SetPassword)

			// Create Request
			req := httptest.NewRequest("PATCH", "/solution?"+test.paramsName+"=1",
				bytes.NewBufferString(test.inputBody))
			req.Header.Set("Content-Type", "application/json")

			// Make Request
			w, err := r.Test(req)
			if err != nil {
				t.Error(err)
			}

			// Read response
			resp, err := ReadResponse(w.Body)
			if err != nil {
				t.Fatal(err)
			}

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.StatusCode)
			assert.Equal(t, test.expectedResponseBody, resp)
		})
	}
}

func TestHandler_ConnectToSolution(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockTask, di *mock_service.MockDiService, input ConnectInput, svc *usecase.Service)
	tests := []struct {
		name                 string
		userIdentify         func(c *fiber.Ctx) (int64, error)
		inputBody            string
		inputRequest         ConnectInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody: `{"sid": 1, "password": "qwerty"}`,
			inputRequest: ConnectInput{
				SID:      1,
				Password: "qwerty",
			},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, input ConnectInput, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().ConnectToTask(context.Background(), input.SID, input.Password).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"success"}`,
		},
		{
			name: "Wrong Input",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody:    `{"sid": 1}`,
			inputRequest: ConnectInput{},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, input ConnectInput, svc *usecase.Service) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input to connect to task"}`,
		},
		{
			name: "Server error",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 0, errors.New("cant read id")
			},
			inputBody:    ``,
			inputRequest: ConnectInput{},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, input ConnectInput, svc *usecase.Service) {
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"cant read id"}`,
		},
		{
			name: "Service Error",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody: `{"sid": 1, "password": "qwerty"}`,
			inputRequest: ConnectInput{
				SID:      1,
				Password: "qwerty",
			},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, input ConnectInput, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().ConnectToTask(context.Background(), input.SID, input.Password).Return(errors.New("something went wrong"))
			},
			expectedStatusCode:   403,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockTask(c)
			di := mock_service.NewMockDiService(c)

			svc := usecase.Service{Task: repo}
			test.mockBehavior(repo, di, test.inputRequest, &svc)

			handler := Handler{di: di, userIdentity: test.userIdentify}

			r := fiber.New()
			r.Get("/connect", handler.ConnectToSolution)

			// Create Request
			req := httptest.NewRequest("GET", "/connect",
				bytes.NewBufferString(test.inputBody))
			req.Header.Set("Content-Type", "application/json")

			// Make Request
			w, err := r.Test(req)
			if err != nil {
				t.Error(err)
			}

			// Read response
			resp, err := ReadResponse(w.Body)
			if err != nil {
				t.Fatal(err)
			}

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.StatusCode)
			assert.Equal(t, test.expectedResponseBody, resp)
		})
	}
}

func TestHandler_SetExpertsWeights(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockTask, di *mock_service.MockDiService, w entity.Weights, svc *usecase.Service)

	type InReq struct {
		Weights entity.Weights `json:"weights"`
	}

	tests := []struct {
		name                 string
		paramsName           string
		userIdentify         func(c *fiber.Ctx) (int64, error)
		inputBody            string
		inputRequest         InReq
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:       "Ok",
			paramsName: "sid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody: `{"weights": [0.5, {"start":0.2, "end":0.4}, 0.6]}`,
			inputRequest: InReq{
				Weights: []eval.Rating{{eval.Number(0.5)}, {eval.Interval{Start: 0.2, End: 0.4}},
					{eval.Number(0.6)}},
			},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, w entity.Weights, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().ValidateUser(context.Background(), int64(1), int64(1)).Return(nil)
				r.EXPECT().SetExpertsWeights(context.Background(), int64(1), w).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"success"}`,
		},
		{
			name:       "Wrong Input",
			paramsName: "sid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody:    `{"weights": "test"}`,
			inputRequest: InReq{},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, w entity.Weights, svc *usecase.Service) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"json: cannot unmarshal string into Go struct field .weights of type entity.Weights"}`,
		},
		{
			name:       "Wrong URL",
			paramsName: "sssssid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody:    `{"weights": [{0.5}, {"start":0.2, "end":"0.4"}, {0.6}]}`,
			inputRequest: InReq{},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, w entity.Weights, svc *usecase.Service) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"task doesnt specified"}`,
		},
		{
			name:       "Validate error",
			paramsName: "sid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody:    `{"weights": [0.5, {"start":0.2, "end":0.4}, 0.6]}`,
			inputRequest: InReq{},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, w entity.Weights, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().ValidateUser(context.Background(), int64(1), int64(1)).Return(errors.New("forbidden"))
			},
			expectedStatusCode:   403,
			expectedResponseBody: `{"message":"forbidden"}`,
		},
		{
			name:       "Server error",
			paramsName: "sid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 0, errors.New("cant read id")
			},
			inputBody:    ``,
			inputRequest: InReq{},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, w entity.Weights, svc *usecase.Service) {
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"cant read id"}`,
		},
		{
			name:       "Service Error",
			paramsName: "sid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody: `{"weights": [0.5, {"start":0.2, "end":0.4}, 0.6]}`,
			inputRequest: InReq{
				Weights: []eval.Rating{{eval.Number(0.5)}, {eval.Interval{Start: 0.2, End: 0.4}},
					{eval.Number(0.6)}},
			},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, w entity.Weights, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().ValidateUser(context.Background(), int64(1), int64(1)).Return(nil)
				r.EXPECT().SetExpertsWeights(context.Background(), int64(1), w).Return(errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockTask(c)
			di := mock_service.NewMockDiService(c)

			svc := usecase.Service{Task: repo}
			test.mockBehavior(repo, di, test.inputRequest.Weights, &svc)

			handler := Handler{di: di, userIdentity: test.userIdentify}

			r := fiber.New()
			r.Put("/experts", handler.SetExpertsWeights)

			// Create Request
			req := httptest.NewRequest("PUT", "/experts?"+test.paramsName+"=1",
				bytes.NewBufferString(test.inputBody))
			req.Header.Set("Content-Type", "application/json")

			// Make Request
			w, err := r.Test(req)
			if err != nil {
				t.Error(err)
			}

			// Read response
			resp, err := ReadResponse(w.Body)
			if err != nil {
				t.Fatal(err)
			}

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.StatusCode)
			assert.Equal(t, test.expectedResponseBody, resp)
		})
	}
}

func TestHandler_GetExperts(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockTask, u *mock_service.MockUser, di *mock_service.MockDiService,
		task *entity.TaskModel, svc *usecase.Service)

	type InReq struct {
		Weights entity.Weights `json:"weights"`
	}

	tests := []struct {
		name                 string
		paramsName           string
		userIdentify         func(c *fiber.Ctx) (int64, error)
		outTask              entity.TaskModel
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:       "Ok",
			paramsName: "sid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			outTask: entity.TaskModel{
				TaskType: v.Group,
			},
			mockBehavior: func(r *mock_service.MockTask, u *mock_service.MockUser, di *mock_service.MockDiService,
				task *entity.TaskModel, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().ValidateUser(context.Background(), int64(1), int64(1)).Return(nil)
				r.EXPECT().GetTask(context.Background(), int64(1)).Return(task, nil)
				u.EXPECT().GetUsersRelateToTask(context.Background(), int64(1)).Return([]string{"user1", "user2"}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `["user1","user2"]`,
		},
		{
			name:       "Task type error",
			paramsName: "sid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			outTask: entity.TaskModel{
				TaskType: v.Individuals,
			},
			mockBehavior: func(r *mock_service.MockTask, u *mock_service.MockUser, di *mock_service.MockDiService,
				task *entity.TaskModel, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().ValidateUser(context.Background(), int64(1), int64(1)).Return(nil)
				r.EXPECT().GetTask(context.Background(), int64(1)).Return(task, nil)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"individuals task doesnt provide this function"}`,
		},
		{
			name:       "Wrong URL",
			paramsName: "sssssid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			outTask: entity.TaskModel{},
			mockBehavior: func(r *mock_service.MockTask, u *mock_service.MockUser, di *mock_service.MockDiService,
				task *entity.TaskModel, svc *usecase.Service) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"task doesnt specified"}`,
		},
		{
			name:       "Validate error",
			paramsName: "sid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			outTask: entity.TaskModel{},
			mockBehavior: func(r *mock_service.MockTask, u *mock_service.MockUser, di *mock_service.MockDiService,
				task *entity.TaskModel, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().ValidateUser(context.Background(), int64(1), int64(1)).Return(errors.New("forbidden"))
			},
			expectedStatusCode:   403,
			expectedResponseBody: `{"message":"forbidden"}`,
		},
		{
			name:       "Server error",
			paramsName: "sid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 0, errors.New("cant read id")
			},
			outTask: entity.TaskModel{},
			mockBehavior: func(r *mock_service.MockTask, u *mock_service.MockUser, di *mock_service.MockDiService,
				task *entity.TaskModel, svc *usecase.Service) {
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"cant read id"}`,
		},
		{
			name:       "Service Error",
			paramsName: "sid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			outTask: entity.TaskModel{
				TaskType: v.Group,
			},
			mockBehavior: func(r *mock_service.MockTask, u *mock_service.MockUser, di *mock_service.MockDiService,
				task *entity.TaskModel, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().ValidateUser(context.Background(), int64(1), int64(1)).Return(nil)
				r.EXPECT().GetTask(context.Background(), int64(1)).Return(task, nil)
				u.EXPECT().GetUsersRelateToTask(context.Background(), int64(1)).Return(nil, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
		{
			name:       "Service Error 2",
			paramsName: "sid",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			outTask: entity.TaskModel{},
			mockBehavior: func(r *mock_service.MockTask, u *mock_service.MockUser, di *mock_service.MockDiService,
				task *entity.TaskModel, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().ValidateUser(context.Background(), int64(1), int64(1)).Return(nil)
				r.EXPECT().GetTask(context.Background(), int64(1)).Return(nil, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockTask(c)
			userRepo := mock_service.NewMockUser(c)
			di := mock_service.NewMockDiService(c)

			svc := usecase.Service{Task: repo, User: userRepo}
			test.mockBehavior(repo, userRepo, di, &test.outTask, &svc)

			handler := Handler{di: di, userIdentity: test.userIdentify}

			r := fiber.New()
			r.Get("/experts", handler.GetExperts)

			// Create Request
			req := httptest.NewRequest("GET", "/experts?"+test.paramsName+"=1",
				bytes.NewBufferString(""))
			req.Header.Set("Content-Type", "application/json")

			// Make Request
			w, err := r.Test(req)
			if err != nil {
				t.Error(err)
			}

			// Read response
			resp, err := ReadResponse(w.Body)
			if err != nil {
				t.Fatal(err)
			}

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.StatusCode)
			assert.Equal(t, test.expectedResponseBody, resp)
		})
	}
}
