package controller

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"strconv"
	"webApp/lib/eval"
)

// CreateMatrix godoc
// @summary CreateMatrix
// @description creates new matrix for evaluating
// @security ApiKeyAuth
// @id create-matrix
// @tags matrix
// @accept json
// @produce json
// @param sid query int true "task identifier"
// @success 200 {object} response
// @success 404 {object} response
// @failure 500 {object} response
// @router /solution/rating [post]
func (h *Handler) CreateMatrix(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	sid, err := strconv.ParseInt(c.Query("sid"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusNotFound, errors.New("task doesn't specified"))
	}

	service := h.di.GetInstanceService()
	mid, err := service.Matrix.CreateMatrix(c.UserContext(), uid, sid)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	logrus.Infof("user with id %d successfully create new matrix with mid %d", uid, mid)
	return c.JSON(response{Message: "success"})
}

type RatingsInput struct {
	Ratings []eval.Rating `json:"ratings"`
}

// UpdateMatrix godoc
// @summary UpdateMatrix
// @description updates ratings in matrix
// @security ApiKeyAuth
// @id update-matrix
// @tags matrix
// @accept json
// @produce json
// @param input body RatingsInput true "ratings to update"
// @param sid query int true "task identifier"
// @param ord query int true "order of alternative"
// @success 200 {object} response
// @success 400 {object} response
// @success 404 {object} response
// @failure 500 {object} response
// @router /solution/rating [put]
func (h *Handler) UpdateMatrix(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	sid, err := strconv.ParseInt(c.Query("sid"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusNotFound, errors.New("task doesn't specified"))
	}
	ord, err := strconv.ParseInt(c.Query("ord"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusNotFound, errors.New("alternative doesn't specified"))
	}
	ord--

	request := RatingsInput{}
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

	return c.JSON(response{Message: "success"})
}

// GetRatings godoc
// @summary GetRatings
// @description gets ratings on current alternative in matrix
// @security ApiKeyAuth
// @id get-ratings
// @tags matrix
// @accept json
// @produce json
// @param sid query int true "task identifier"
// @param ord query int true "order of alternative"
// @success 200 {object} RatingsInput
// @success 404 {object} response
// @failure 500 {object} response
// @router /solution/rating [get]
func (h *Handler) GetRatings(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	sid, err := strconv.ParseInt(c.Query("sid"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusNotFound, errors.New("task doesn't specified"))
	}

	ord, err := strconv.ParseInt(c.Query("ord"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusNotFound, errors.New("alternative doesn't specified"))
	}
	ord--

	service := h.di.GetInstanceService()
	ratings, err := service.Matrix.GetRatings(c.UserContext(), uid, sid, ord)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(RatingsInput{ratings})
}

// CompleteStatus godoc
// @summary CompleteStatus
// @description set status complete for current matrix
// @security ApiKeyAuth
// @id set-complete
// @tags matrix
// @accept json
// @produce json
// @param sid query int true "task identifier"
// @success 200 {object} response
// @success 404 {object} response
// @failure 500 {object} response
// @router /solution/rating/complete [patch]
func (h *Handler) CompleteStatus(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	sid, err := strconv.ParseInt(c.Query("sid"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusNotFound, errors.New("task doesn't specified"))
	}

	service := h.di.GetInstanceService()
	mid, err := service.Matrix.GetMID(c.UserContext(), uid, sid)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	if err := service.Matrix.SetStatusComplete(c.UserContext(), mid); err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(response{Message: "success"})
}
