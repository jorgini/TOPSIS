package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"webApp/app/entity"
)

func (h *Handler) SetAlternatives(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	sid, err := strconv.ParseInt(c.Query("sid"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusNotFound, errors.New("this solution not found"))
	}

	service := h.di.GetInstanceService()
	if err := service.Task.ValidateUser(c.UserContext(), uid, sid); err != nil {
		return sendErrorResponse(c, fiber.StatusForbidden, err)
	}

	var request entity.Alts
	if err := c.BodyParser(&request); err != nil || request == nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, err)
	}

	if err := service.Task.SetAlts(c.UserContext(), sid, request); err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(fiber.Map{"status": "success"})
}

func (h *Handler) ReplaceAlternatives(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	sid, err := strconv.ParseInt(c.Query("sid"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusNotFound, errors.New("this solution not found"))
	}

	service := h.di.GetInstanceService()
	if err := service.Task.ValidateUser(c.UserContext(), uid, sid); err != nil {
		return sendErrorResponse(c, fiber.StatusForbidden, err)
	}

	var request entity.Alts
	if err := c.BodyParser(&request); err != nil || request == nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, err)
	}

	if err := service.Task.UpdateAlts(c.UserContext(), sid, request); err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(fiber.Map{"status": "success"})
}

func (h *Handler) GetAlternatives(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	sid, err := strconv.ParseInt(c.Query("sid"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusNotFound, errors.New("this solution not found"))
	}

	service := h.di.GetInstanceService()
	if err := service.Task.CheckAccess(c.UserContext(), uid, sid); err != nil {
		return sendErrorResponse(c, fiber.StatusForbidden, errors.New("hasn't access to solution"))
	}

	alts, err := service.GetAlts(c.UserContext(), sid)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(alts)
}

func (h *Handler) SetCriteria(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	sid, err := strconv.ParseInt(c.Query("sid"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusNotFound, errors.New("this solution not found"))
	}

	service := h.di.GetInstanceService()
	if err := service.Task.ValidateUser(c.UserContext(), uid, sid); err != nil {
		return sendErrorResponse(c, fiber.StatusForbidden, err)
	}

	var request entity.Criteria
	if err := c.BodyParser(&request); err != nil || request == nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, err)
	}

	if err := service.Task.SetCriteria(c.UserContext(), sid, request); err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(fiber.Map{"status": "success"})
}

func (h *Handler) ReplaceCriteria(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	sid, err := strconv.ParseInt(c.Query("sid"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusNotFound, errors.New("this solution not found"))
	}

	service := h.di.GetInstanceService()
	if err := service.Task.ValidateUser(c.UserContext(), uid, sid); err != nil {
		return sendErrorResponse(c, fiber.StatusForbidden, err)
	}

	var request entity.Criteria
	if err := c.BodyParser(&request); err != nil || request == nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, err)
	}

	if err := service.Task.UpdateCriteria(c.UserContext(), sid, request); err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(fiber.Map{"status": "success"})
}

func (h *Handler) GetCriteria(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	sid, err := strconv.ParseInt(c.Query("sid"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusNotFound, errors.New("this solution not found"))
	}

	service := h.di.GetInstanceService()
	if err := service.Task.CheckAccess(c.UserContext(), uid, sid); err != nil {
		return sendErrorResponse(c, fiber.StatusForbidden, errors.New("hasn't access to solution"))
	}

	criteria, err := service.GetCriteria(c.UserContext(), sid)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(criteria)
}
