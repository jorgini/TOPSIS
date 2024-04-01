package controller

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"time"
	"webApp/entity"
)

// CreateNewUser godoc
// @summary Sign-Up
// @description sign-up for new user
// @id sign-up
// @tags auth
// @accept json
// @produce json
// @param input body entity.UserModel true "account info"
// @success 200 {object} response
// @failure 400 {object} response
// @failure 500 {object} response
// @router /auth/sign-up [post]
func (h *Handler) CreateNewUser(c *fiber.Ctx) error {
	var user entity.UserModel
	if err := c.BodyParser(&user); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, err)
	}

	service := h.di.GetInstanceService()
	id, err := service.User.CreateNewUser(c.UserContext(), &user)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	logrus.Infof("successful create new user with id %d", id)
	return c.JSON(response{Message: "success"})
}

type UserInput struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserInput) UnmarshalJSON(data []byte) error {
	result := struct {
		Login    string `json:"login"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.Email == "" && result.Password == "" && result.Login == "" {
		return errors.New("invalid input data for update user")
	} else {
		u.Login = result.Login
		u.Email = result.Email
		u.Password = result.Password
	}
	return nil
}

// LogIn godoc
// @summary Log-In
// @description log-in for user
// @id log-in
// @tags auth
// @accept json
// @produce json
// @param input body UserInput true "account info"
// @success 200 {object} string
// @failure 400 {object} response
// @failure 401 {object} response
// @router /auth/log-in [post]
func (h *Handler) LogIn(c *fiber.Ctx) error {
	var input UserInput
	err := c.BodyParser(&input)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, err)
	}

	if input.Login == "" && input.Email == "" || input.Password == "" {
		return sendErrorResponse(c, fiber.StatusBadRequest, errors.New("empty input"))
	}

	user := entity.UserModel{
		Login:    input.Login,
		Email:    input.Email,
		Password: input.Password,
	}

	service := h.di.GetInstanceService()
	uid, err := service.User.GetUID(c.UserContext(), &user)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, err)
	}

	tokens, err := service.Session.GenerateToken(c.UserContext(), uid, h.cfg)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, err)
	}

	c.Cookie(&fiber.Cookie{
		Name:     h.cfg.CookieName,
		Value:    tokens.Refresh,
		Expires:  time.Now().Add(48 * time.Hour),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{"token": tokens.Access})
}

// UpdateUser godoc
// @summary UpdateUser
// @description update account info
// @security ApiKeyAuth
// @id update-user
// @tags user
// @accept json
// @produce json
// @param input body UserInput true "update account info"
// @success 200 {object} response
// @failure 400 {object} response
// @failure 401 {string} string
// @failure 500 {object} response
// @router /user/settings [patch]
func (h *Handler) UpdateUser(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	var update UserInput
	if err := c.BodyParser(&update); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, err)
	}
	user := entity.UserModel{
		Email:    update.Email,
		Password: update.Password,
	}

	service := h.di.GetInstanceService()
	if err := service.User.UpdateUser(c.UserContext(), uid, &user); err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(response{Message: "success"})
}

// DeleteUser godoc
// @summary DeleteUser
// @description delete user account
// @security ApiKeyAuth
// @id delete-user
// @tags user
// @accept json
// @produce json
// @success 200 {object} response
// @failure 401 {object} response
// @failure 500 {object} response
// @router /user/settings [delete]
func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	service := h.di.GetInstanceService()
	if err := service.User.DeleteUser(c.UserContext(), uid); err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}
	return c.JSON(response{Message: "success"})
}

// RefreshToken godoc
// @summary RefreshToken
// @description refresh tokens for identify user
// @id refresh-token
// @tags auth
// @accept json
// @produce json
// @success 200 {object} response
// @failure 401 {object} response
// @failure 500 {object} response
// @router /auth/refresh [put]
func (h *Handler) RefreshToken(c *fiber.Ctx) error {
	refresh := c.Cookies(h.cfg.CookieName)
	if refresh == "" {
		return sendErrorResponse(c, fiber.StatusUnauthorized, errors.New("cant find refresh token"))
	}

	service := h.di.GetInstanceService()
	tokens, err := service.Session.RefreshToken(c.UserContext(), refresh, h.cfg)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, err)
	}

	c.Cookie(&fiber.Cookie{
		Name:     h.cfg.CookieName,
		Value:    tokens.Refresh,
		Expires:  time.Now().Add(360 * time.Hour),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{"token": tokens.Access})
}

// GetAllSolutions godoc
// @summary GetAllSolutions
// @description get all solutions related to current user
// @security ApiKeyAuth
// @id get-all-solutions
// @tags user
// @accept json
// @produce json
// @success 200 {object} []entity.TaskShortCard
// @failure 401 {object} response
// @failure 500 {object} response
// @router /user [get]
func (h *Handler) GetAllSolutions(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	svc := h.di.GetInstanceService()
	tasks, err := svc.Task.GetAllSolutions(c.UserContext(), uid)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(tasks)
}
