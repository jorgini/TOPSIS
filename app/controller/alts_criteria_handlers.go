package controller

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"webApp/entity"
)

// ReplaceAlternatives godoc
// @summary ReplaceAlternatives
// @description updates alternatives to current task
// @security ApiKeyAuth
// @id update-alts
// @tags alternatives
// @accept json
// @produce json
// @param input body entity.Alts true "alternatives info"
// @param sid query int true "task identifier"
// @success 200 {object} response
// @success 400 {object} response
// @success 403 {object} response
// @failure 404 {object} response
// @failure 500 {object} response
// @router /solution/alternatives [put]
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

	return c.JSON(response{Message: "success"})
}

// GetAlternatives godoc
// @summary GetAlternatives
// @description gets alternatives to current task
// @security ApiKeyAuth
// @id get-alts
// @tags alternatives
// @accept json
// @produce json
// @param sid query int true "task identifier"
// @success 200 {object} entity.Alts
// @success 403 {object} response
// @failure 404 {object} response
// @failure 500 {object} response
// @router /solution/alternatives [get]
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

// ReplaceCriteria godoc
// @summary ReplaceCriteria
// @description updates criteria to current task
// @security ApiKeyAuth
// @id update-criteria
// @tags criteria
// @accept json
// @produce json
// @param input body entity.Criteria true "criteria info"
// @param sid query int true "task identifier"
// @success 200 {object} response
// @success 400 {object} response
// @success 403 {object} response
// @failure 404 {object} response
// @failure 500 {object} response
// @router /solution/criteria [put]
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

// GetCriteria godoc
// @summary GetCriteria
// @description gets criteria to current task
// @security ApiKeyAuth
// @id get-criteria
// @tags criteria
// @accept json
// @produce json
// @param sid query int true "task identifier"
// @success 200 {object} entity.Criteria
// @success 403 {object} response
// @failure 404 {object} response
// @failure 500 {object} response
// @router /solution/criteria [get]
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
