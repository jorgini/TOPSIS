package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"time"
	"webApp/entity"
)

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
	return c.JSON(fiber.Map{"status": "success"})
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

	return c.JSON(fiber.Map{"status": "success"})
}

func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	service := h.di.GetInstanceService()
	if err := service.User.DeleteUser(c.UserContext(), uid); err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}
	return c.JSON(fiber.Map{"status": "success"})
}

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
		Expires:  time.Now().Add(48 * time.Hour),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{"token": tokens.Access})
}

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
