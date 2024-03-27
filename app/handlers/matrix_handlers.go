package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"strconv"
	"webApp/lib/eval"
)

func (h *Handler) CreateMatrix(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	sid, err := strconv.ParseInt(c.Query("sid"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, errors.New("task doesn't specified"))
	}

	service := h.di.GetInstanceService()
	mid, err := service.Matrix.CreateMatrix(c.UserContext(), uid, sid)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	logrus.Infof("user with id %d successfully create new matrix with mid %d", uid, mid)
	return c.JSON(fiber.Map{"status": "success"})
}

type RatingsInput []eval.Rating

func (h *Handler) UpdateMatrix(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	sid, err := strconv.ParseInt(c.Query("sid"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, errors.New("task doesn't specified"))
	}

	ord, err := strconv.ParseInt(c.Query("ord"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, errors.New("alternative doesn't specified"))
	}
	ord--

	request := struct {
		Ratings RatingsInput `json:"ratings"`
	}{}
	if err := c.BodyParser(&request); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, err)
	}

	service := h.di.GetInstanceService()
	mid, err := service.Matrix.GetMID(c.UserContext(), uid, sid)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	if err := service.Matrix.UpdateMatrix(c.UserContext(), sid, mid, ord, request.Ratings); err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(fiber.Map{"status": "success"})
}

func (h *Handler) GetRatings(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	sid, err := strconv.ParseInt(c.Query("sid"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, errors.New("task doesn't specified"))
	}

	ord, err := strconv.ParseInt(c.Query("ord"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, errors.New("alternative doesn't specified"))
	}
	ord--

	service := h.di.GetInstanceService()
	ratings, err := service.Matrix.GetRatings(c.UserContext(), uid, sid, ord)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(ratings)
}
