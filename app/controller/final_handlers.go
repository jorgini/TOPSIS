package controller

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type ThresholdInput struct {
	Threshold float64 `json:"threshold"`
}

// GetFinal godoc
// @summary GetFinal
// @description gets final report
// @security ApiKeyAuth
// @id get-final
// @tags final
// @accept json
// @produce json
// @param input body ThresholdInput false "threshold for sensitivity analysis"
// @param sid query int true "task identifier"
// @success 200 {object} entity.FinalModel
// @success 400 {object} response
// @success 403 {object} response
// @success 404 {object} response
// @failure 500 {object} response
// @router /solution/final [post]
func (h *Handler) GetFinal(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	sid, err := strconv.ParseInt(c.Query("sid"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusNotFound, errors.New("task doesn't specified"))
	}

	request := ThresholdInput{}
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
