package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"webApp/configs"
	"webApp/entity"
	"webApp/usecase"
	mock_service "webApp/usecase/mocks-service"
)

func TestHandler_CreateNewUser(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockUser, di *mock_service.MockDiService, user entity.UserModel, svc *usecase.Service)
	tests := []struct {
		name                 string
		inputBody            string
		inputUser            entity.UserModel
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"login": "login", "email": "test", "password": "test"}`,
			inputUser: entity.UserModel{
				Login:    "login",
				Email:    "test",
				Password: "test",
			},
			mockBehavior: func(r *mock_service.MockUser, di *mock_service.MockDiService, user entity.UserModel, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().CreateNewUser(context.Background(), &user).Return(int64(1), nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"success"}`,
		},
		{
			name:      "Wrong Input",
			inputBody: `{"email": "test"}`,
			inputUser: entity.UserModel{},
			mockBehavior: func(r *mock_service.MockUser, di *mock_service.MockDiService, user entity.UserModel, svc *usecase.Service) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input data for sign up"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"login": "login", "email": "test", "password": "test"}`,
			inputUser: entity.UserModel{
				Login:    "login",
				Email:    "test",
				Password: "test",
			},
			mockBehavior: func(r *mock_service.MockUser, di *mock_service.MockDiService, user entity.UserModel, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().CreateNewUser(context.Background(), &user).Return(int64(0), errors.New("something went wrong"))
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

			repo := mock_service.NewMockUser(c)
			di := mock_service.NewMockDiService(c)

			svc := usecase.Service{User: repo}
			test.mockBehavior(repo, di, test.inputUser, &svc)

			handler := Handler{di: di}

			r := fiber.New()
			r.Post("/user", handler.CreateNewUser)

			// Create Request
			req := httptest.NewRequest("POST", "/user",
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

func TestHandler_LogIn(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r1 *mock_service.MockUser, r2 *mock_service.MockSession, di *mock_service.MockDiService,
		user *entity.UserModel, svc *usecase.Service, cfg *configs.AppConfig, tokens entity.Tokens)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            entity.UserModel
		tokens               entity.Tokens
		cfg                  *configs.AppConfig
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
		expectedCookie       int
	}{
		{
			name:      "Ok",
			inputBody: `{"email": "test", "password": "test"}`,
			inputUser: entity.UserModel{
				Email:    "test",
				Password: "test",
			},
			tokens: entity.Tokens{
				Access:  "access",
				Refresh: "refresh",
			},
			cfg: &configs.AppConfig{
				CookieName: "cookie",
			},
			mockBehavior: func(r1 *mock_service.MockUser, r2 *mock_service.MockSession, di *mock_service.MockDiService,
				user *entity.UserModel, svc *usecase.Service, cfg *configs.AppConfig, tokens entity.Tokens) {
				di.EXPECT().GetInstanceService().Return(svc)
				r1.EXPECT().GetUID(context.Background(), user).Return(int64(1), nil)
				r2.EXPECT().GenerateToken(context.Background(), int64(1), cfg).Return(tokens, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"token":"access"}`,
			expectedCookie:       1,
		},
		{
			name:      "Ok",
			inputBody: `{"login": "test", "password": "test"}`,
			inputUser: entity.UserModel{
				Login:    "test",
				Password: "test",
			},
			tokens: entity.Tokens{
				Access:  "access",
				Refresh: "refresh",
			},
			cfg: &configs.AppConfig{
				CookieName: "cookie",
			},
			mockBehavior: func(r1 *mock_service.MockUser, r2 *mock_service.MockSession, di *mock_service.MockDiService,
				user *entity.UserModel, svc *usecase.Service, cfg *configs.AppConfig, tokens entity.Tokens) {
				di.EXPECT().GetInstanceService().Return(svc)
				r1.EXPECT().GetUID(context.Background(), user).Return(int64(1), nil)
				r2.EXPECT().GenerateToken(context.Background(), int64(1), cfg).Return(tokens, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"token":"access"}`,
			expectedCookie:       1,
		},
		{
			name:      "Wrong Input",
			inputBody: `{"password": "qwerty"}`,
			inputUser: entity.UserModel{},
			mockBehavior: func(r1 *mock_service.MockUser, r2 *mock_service.MockSession, di *mock_service.MockDiService,
				user *entity.UserModel, svc *usecase.Service, cfg *configs.AppConfig, tokens entity.Tokens) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"empty input"}`,
			expectedCookie:       0,
		},
		{
			name:      "Unauthorized",
			inputBody: `{"email": "test", "password": "test"}`,
			inputUser: entity.UserModel{
				Email:    "test",
				Password: "test",
			},
			mockBehavior: func(r1 *mock_service.MockUser, r2 *mock_service.MockSession, di *mock_service.MockDiService,
				user *entity.UserModel, svc *usecase.Service, cfg *configs.AppConfig, tokens entity.Tokens) {
				di.EXPECT().GetInstanceService().Return(svc)
				r1.EXPECT().GetUID(context.Background(), user).Return(int64(1), nil)
				r2.EXPECT().GenerateToken(context.Background(), int64(1), cfg).Return(tokens, errors.New("something went wrong"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"something went wrong"}`,
			expectedCookie:       0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			userRepo := mock_service.NewMockUser(c)
			sessionRepo := mock_service.NewMockSession(c)
			di := mock_service.NewMockDiService(c)

			svc := usecase.Service{User: userRepo, Session: sessionRepo}
			test.mockBehavior(userRepo, sessionRepo, di, &test.inputUser, &svc, test.cfg, test.tokens)

			handler := Handler{di: di, cfg: test.cfg}

			r := fiber.New()
			r.Get("/user", handler.LogIn)

			// Create Request
			req := httptest.NewRequest("GET", "/user",
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
			cookies := w.Cookies()

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.StatusCode)
			assert.Equal(t, test.expectedResponseBody, resp)
			assert.Equal(t, test.expectedCookie, len(cookies))
			if test.expectedCookie == 1 {
				if len(cookies) == 1 {
					assert.Equal(t, test.cfg.CookieName, cookies[0].Name)
					assert.Equal(t, test.tokens.Refresh, cookies[0].Value)
					assert.True(t, cookies[0].HttpOnly, "cookie is httpOnly")
				}
			}
		})
	}
}

func TestHandler_UpdateUser(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockUser, di *mock_service.MockDiService, user *entity.UserModel, svc *usecase.Service)
	tests := []struct {
		name                 string
		userIdentify         func(c *fiber.Ctx) (int64, error)
		inputBody            string
		inputUser            entity.UserModel
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody: `{"email": "test", "password": "test"}`,
			inputUser: entity.UserModel{
				Email:    "test",
				Password: "test",
			},
			mockBehavior: func(r *mock_service.MockUser, di *mock_service.MockDiService, user *entity.UserModel, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().UpdateUser(context.Background(), int64(1), user).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"success"}`,
		},
		{
			name: "Partition request",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody: `{"password": "test"}`,
			inputUser: entity.UserModel{
				Password: "test",
			},
			mockBehavior: func(r *mock_service.MockUser, di *mock_service.MockDiService, user *entity.UserModel, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().UpdateUser(context.Background(), int64(1), user).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"success"}`,
		},
		{
			name: "Wrong Input",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody: `{}`,
			inputUser: entity.UserModel{},
			mockBehavior: func(r *mock_service.MockUser, di *mock_service.MockDiService, user *entity.UserModel, svc *usecase.Service) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input data for update user"}`,
		},
		{
			name: "Server error",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 0, errors.New("cant read id")
			},
			mockBehavior: func(r *mock_service.MockUser, di *mock_service.MockDiService, user *entity.UserModel, svc *usecase.Service) {
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"cant read id"}`,
		},
		{
			name: "Service Error",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			inputBody: `{"email": "test", "password": "test"}`,
			inputUser: entity.UserModel{
				Email:    "test",
				Password: "test",
			},
			mockBehavior: func(r *mock_service.MockUser, di *mock_service.MockDiService, user *entity.UserModel, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().UpdateUser(context.Background(), int64(1), user).Return(errors.New("something went wrong"))
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

			repo := mock_service.NewMockUser(c)
			di := mock_service.NewMockDiService(c)

			svc := usecase.Service{User: repo}
			test.mockBehavior(repo, di, &test.inputUser, &svc)

			handler := Handler{di: di, userIdentity: test.userIdentify}

			r := fiber.New()
			r.Patch("/user", handler.UpdateUser)

			// Create Request
			req := httptest.NewRequest("PATCH", "/user",
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

func TestHandler_DeleteUser(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockUser, di *mock_service.MockDiService, svc *usecase.Service)
	tests := []struct {
		name                 string
		userIdentify         func(c *fiber.Ctx) (int64, error)
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			mockBehavior: func(r *mock_service.MockUser, di *mock_service.MockDiService, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().DeleteUser(context.Background(), int64(1)).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"success"}`,
		},
		{
			name: "Server error",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 0, errors.New("cant read id")
			},
			mockBehavior: func(r *mock_service.MockUser, di *mock_service.MockDiService, svc *usecase.Service) {
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"cant read id"}`,
		},
		{
			name: "Service Error",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			mockBehavior: func(r *mock_service.MockUser, di *mock_service.MockDiService, svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().DeleteUser(context.Background(), int64(1)).Return(errors.New("something went wrong"))
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

			repo := mock_service.NewMockUser(c)
			di := mock_service.NewMockDiService(c)

			svc := usecase.Service{User: repo}
			test.mockBehavior(repo, di, &svc)

			handler := Handler{di: di, userIdentity: test.userIdentify}

			r := fiber.New()
			r.Patch("/user", handler.DeleteUser)

			// Create Request
			req := httptest.NewRequest("PATCH", "/user",
				bytes.NewBufferString(""))

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

func TestHandler_RefreshToken(t *testing.T) {
	// Init Test Table
	type mockBehavior func(s *mock_service.MockSession, di *mock_service.MockDiService,
		svc *usecase.Service, refresh string, cfg *configs.AppConfig, tokens entity.Tokens)
	tests := []struct {
		name                 string
		cookie               http.Cookie
		cfg                  *configs.AppConfig
		tokens               entity.Tokens
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
		expectedCookie       int
	}{
		{
			name: "Ok",
			cookie: http.Cookie{
				Name:  "cookie",
				Value: "refresh",
			},
			cfg: &configs.AppConfig{
				CookieName: "cookie",
			},
			tokens: entity.Tokens{
				Access:  "access",
				Refresh: "refresh",
			},
			mockBehavior: func(s *mock_service.MockSession, di *mock_service.MockDiService,
				svc *usecase.Service, refresh string, cfg *configs.AppConfig, tokens entity.Tokens) {
				di.EXPECT().GetInstanceService().Return(svc)
				s.EXPECT().RefreshToken(context.Background(), refresh, cfg).Return(tokens, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"access":"access"}`,
			expectedCookie:       1,
		},
		{
			name: "Incorrect cookie",
			cookie: http.Cookie{
				Name: "test",
			},
			cfg: &configs.AppConfig{
				CookieName: "cookie",
			},
			mockBehavior: func(s *mock_service.MockSession, di *mock_service.MockDiService,
				svc *usecase.Service, refresh string, cfg *configs.AppConfig, tokens entity.Tokens) {
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"cant find refresh token"}`,
			expectedCookie:       0,
		},
		{
			name: "Unauthorized",
			cookie: http.Cookie{
				Name:  "cookie",
				Value: "refresh",
			},
			cfg: &configs.AppConfig{
				CookieName: "cookie",
			},
			mockBehavior: func(s *mock_service.MockSession, di *mock_service.MockDiService,
				svc *usecase.Service, refresh string, cfg *configs.AppConfig, tokens entity.Tokens) {
				di.EXPECT().GetInstanceService().Return(svc)
				s.EXPECT().RefreshToken(context.Background(), refresh, cfg).Return(tokens, errors.New("something went wrong"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"something went wrong"}`,
			expectedCookie:       0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockSession(c)
			di := mock_service.NewMockDiService(c)

			svc := usecase.Service{Session: repo}
			test.mockBehavior(repo, di, &svc, test.cookie.Value, test.cfg, test.tokens)

			handler := Handler{di: di, cfg: test.cfg}

			r := fiber.New()
			r.Put("/refresh", handler.RefreshToken)

			// Create Request
			req := httptest.NewRequest("PUT", "/refresh",
				bytes.NewBufferString(""))
			req.AddCookie(&test.cookie)

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
			cookies := w.Cookies()

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.StatusCode)
			assert.Equal(t, test.expectedResponseBody, resp)
			assert.Equal(t, test.expectedCookie, len(cookies))
			if test.expectedCookie == 1 && len(cookies) == 1 {
				assert.Equal(t, test.tokens.Refresh, cookies[0].Value)
				assert.Equal(t, test.cfg.CookieName, cookies[0].Name)
				assert.True(t, cookies[0].HttpOnly, "cookie is httpOnly")
			}
		})
	}
}

func TestHandler_GetAllSolutions(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockTask, di *mock_service.MockDiService, task entity.TaskShortCard,
		svc *usecase.Service)
	tests := []struct {
		name                 string
		userIdentify         func(c *fiber.Ctx) (int64, error)
		outputTask           entity.TaskShortCard
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			outputTask: entity.TaskShortCard{
				Title:       "title",
				Description: "description",
				TaskType:    "group",
				Method:      "topsis",
				LastChange:  time.Now(),
				Status:      entity.Draft,
			},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, task entity.TaskShortCard,
				svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().GetAllSolutions(context.Background(), int64(1)).Return([]entity.TaskShortCard{task}, nil)
			},
			expectedStatusCode: 200,
		},
		{
			name: "Server error",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 0, errors.New("cant read id")
			},
			outputTask: entity.TaskShortCard{},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, task entity.TaskShortCard,
				svc *usecase.Service) {
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"cant read id"}`,
		},
		{
			name: "Service Error",
			userIdentify: func(c *fiber.Ctx) (int64, error) {
				return 1, nil
			},
			outputTask: entity.TaskShortCard{},
			mockBehavior: func(r *mock_service.MockTask, di *mock_service.MockDiService, task entity.TaskShortCard,
				svc *usecase.Service) {
				di.EXPECT().GetInstanceService().Return(svc)
				r.EXPECT().GetAllSolutions(context.Background(), int64(1)).Return(nil, errors.New("something went wrong"))
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
			test.mockBehavior(repo, di, test.outputTask, &svc)

			handler := Handler{di: di, userIdentity: test.userIdentify}

			r := fiber.New()
			r.Get("/user", handler.GetAllSolutions)

			// Create Request
			req := httptest.NewRequest("GET", "/user",
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
				data, err := json.Marshal([]entity.TaskShortCard{test.outputTask})
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
