package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func (h *Handler) GetFinal(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	sid, err := strconv.ParseInt(c.Query("sid"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, errors.New("task doesn't specified"))
	}

	request := struct {
		Threshold float64 `json:"threshold"`
	}{}
	if err := c.BodyParser(&request); err != nil {
		request.Threshold = -1
	}

	svc := h.di.GetInstanceService()
	if err := svc.Task.CheckAccess(c.UserContext(), uid, sid); err != nil {
		return sendErrorResponse(c, fiber.StatusForbidden, err)
	}

	if err := svc.Matrix.IsAllStatusesComplete(c.UserContext(), sid); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, err)
	}

	result, err := svc.Final.PresentFinal(c.UserContext(), sid, request.Threshold)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(result)
}
