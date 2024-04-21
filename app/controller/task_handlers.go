package controller

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"strconv"
	"webApp/entity"
	"webApp/lib/eval"
	v "webApp/lib/variables"
)

type TaskInput struct {
	SID          int64                `json:"sid"`
	Title        string               `json:"title"`
	Description  string               `json:"description"`
	TaskType     string               `json:"task_type"`
	Method       string               `json:"method"`
	CalcSettings int64                `json:"calc_settings"`
	LingScale    eval.LinguisticScale `json:"ling_scale"`
}

func (t *TaskInput) UnmarshalJSON(data []byte) error {
	result := struct {
		SID          int64                 `json:"sid"`
		Title        *string               `json:"title"`
		Description  string                `json:"description"`
		TaskType     *string               `json:"task_type"`
		Method       *string               `json:"method"`
		CalcSettings *int64                `json:"calc_settings"`
		LingScale    *eval.LinguisticScale `json:"ling_scale"`
	}{}

	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.Title == nil || result.TaskType == nil || result.Method == nil || result.CalcSettings == nil ||
		(*result.Method != v.TOPSIS && *result.Method != v.SMART) ||
		(*result.TaskType != v.Individuals && *result.TaskType != v.Group) {
		return errors.New("invalid input arguments for task, check required fields")
	} else {
		t.SID = result.SID
		t.Title = *result.Title
		t.Description = result.Description
		t.Method = *result.Method
		t.TaskType = *result.TaskType
		t.CalcSettings = *result.CalcSettings
		if result.LingScale != nil {
			t.LingScale = *result.LingScale
		} else {
			t.LingScale = *eval.DefaultT1FSScale
		}
	}
	return nil
}

type TitleInput struct {
	Title string `json:"title"`
}

// CreateTask godoc
// @summary CreateTask
// @description creates new task with required options
// @security ApiKeyAuth
// @id create-task
// @tags task
// @accept json
// @produce json
// @param input body TitleInput true "task title"
// @success 200 {object} response
// @success 400 {object} response
// @failure 500 {object} response
// @router /solution/settings [post]
func (h *Handler) CreateTask(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	var input TitleInput
	if err := c.BodyParser(&input); err != nil || input.Title == "" {
		return sendErrorResponse(c, fiber.StatusBadRequest, errors.New("invalid empty title"))
	}

	svc := h.di.GetInstanceService()
	sid, err := svc.Task.CreateNewTask(c.UserContext(), input.Title, uid)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	logrus.Infof("user with id %d successfully create new task with sid %d", uid, sid)
	return c.JSON(fiber.Map{"message": "success", "sid": sid})
}

// UpdateTask godoc
// @summary UpdateTask
// @description updates options for current task
// @security ApiKeyAuth
// @id update-task
// @tags task
// @accept json
// @produce json
// @param sid query int true "task identifier"
// @param input body TaskInput true "task options"
// @success 200 {object} response
// @success 400 {object} response
// @success 403 {object} response
// @success 404 {object} response
// @failure 500 {object} response
// @router /solution/settings [put]
func (h *Handler) UpdateTask(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	sid, err := strconv.ParseInt(c.Query("sid", ""), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusNotFound, errors.New("task doesn't specified"))
	}

	var input TaskInput
	if err := c.BodyParser(&input); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, err)
	}
	task := entity.TaskModel{
		Title:        input.Title,
		Description:  input.Description,
		MaintainerID: uid,
		TaskType:     input.TaskType,
		Method:       input.Method,
		CalcSettings: input.CalcSettings,
		LingScale:    input.LingScale,
		Status:       entity.Draft,
	}

	svc := h.di.GetInstanceService()
	if err := svc.Task.ValidateUser(c.UserContext(), uid, sid); err != nil {
		return sendErrorResponse(c, fiber.StatusForbidden, err)
	}

	if err := svc.Task.UpdateTask(c.UserContext(), sid, &task); err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}
	return c.JSON(response{Message: "success"})
}

// GetTask godoc
// @summary GetTask
// @description gets options of current task
// @security ApiKeyAuth
// @id get-task
// @tags task
// @accept json
// @produce json
// @param sid query int true "task identifier"
// @success 200 {object} TaskInput
// @success 404 {object} response
// @failure 500 {object} response
// @router /solution/settings [get]
func (h *Handler) GetTask(c *fiber.Ctx) error {
	_, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	sid, err := strconv.ParseInt(c.Query("sid"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusNotFound, errors.New("task doesn't specified"))
	}

	svc := h.di.GetInstanceService()

	task, err := svc.Task.GetTask(c.UserContext(), sid)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	output := TaskInput{
		SID:          task.SID,
		Title:        task.Title,
		Description:  task.Description,
		TaskType:     task.TaskType,
		Method:       task.Method,
		CalcSettings: task.CalcSettings,
		LingScale:    task.LingScale,
	}

	return c.JSON(output)
}

// DeleteTask godoc
// @summary DeleteTask
// @description deletes current task
// @security ApiKeyAuth
// @id delete-task
// @tags task
// @accept json
// @produce json
// @param sid query int true "task identifier"
// @success 200 {object} response
// @success 403 {object} response
// @success 404 {object} response
// @failure 500 {object} response
// @router /user/tasks [delete]
func (h *Handler) DeleteTask(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	sid, err := strconv.ParseInt(c.Query("sid"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusNotFound, errors.New("task doesn't specified"))
	}

	svc := h.di.GetInstanceService()
	if err := svc.Task.CheckAccess(c.UserContext(), uid, sid); err != nil {
		return sendErrorResponse(c, fiber.StatusForbidden, errors.New("hasn't access to solution"))
	}

	if err := svc.Task.DeleteTask(c.UserContext(), uid, sid); err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	logrus.Infof("user with id %d successfully delete task with sid %d", uid, sid)
	return c.JSON(response{Message: "success"})
}

// GetRole godoc
// @summary GetRole
// @description gets a role of user for current task
// @security ApiKeyAuth
// @id get-task-password
// @tags task
// @accept json
// @produce json
// @param sid query int true "task identifier"
// @success 200 {object} response
// @success 403 {object} response
// @success 404 {object} response
// @failure 500 {object} response
// @router /solution/experts/role [get]
func (h *Handler) GetRole(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	sid, err := strconv.ParseInt(c.Query("sid"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusNotFound, errors.New("task doesn't specified"))
	}

	svc := h.di.GetInstanceService()
	if err := svc.Task.CheckAccess(c.UserContext(), uid, sid); err != nil {
		return sendErrorResponse(c, fiber.StatusForbidden, errors.New("hasn't access to solution"))
	}

	if err := svc.ValidateUser(c.UserContext(), uid, sid); err != nil {
		return c.JSON(response{Message: "expert"})
	} else {
		return c.JSON(response{Message: "maintainer"})
	}
}

type PasswordInput struct {
	Pass string `json:"password"`
}

// SetPassword godoc
// @summary SetTaskPassword
// @description sets password for current task
// @security ApiKeyAuth
// @id set-task-password
// @tags task
// @accept json
// @produce json
// @param sid query int true "task identifier"
// @param input body PasswordInput true "password for task"
// @success 200 {object} response
// @success 400 {object} response
// @success 403 {object} response
// @success 404 {object} response
// @failure 500 {object} response
// @router /solution/settings/pass [patch]
func (h *Handler) SetPassword(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	sid, err := strconv.ParseInt(c.Query("sid"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusNotFound, errors.New("task doesn't specified"))
	}

	request := PasswordInput{}
	if err := c.BodyParser(&request); err != nil || request.Pass == "" {
		return sendErrorResponse(c, fiber.StatusBadRequest, errors.New("password doesn't specified"))
	}

	svc := h.di.GetInstanceService()
	if err := svc.Task.ValidateUser(c.UserContext(), uid, sid); err != nil {
		return sendErrorResponse(c, fiber.StatusForbidden, err)
	}

	if err := svc.Task.SetPassword(c.UserContext(), sid, request.Pass); err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}
	return c.JSON(response{Message: "success"})
}

type ConnectInput struct {
	SID      int64  `json:"sid"`
	Password string `json:"password"`
}

func (c *ConnectInput) UnmarshalJSON(data []byte) error {
	result := struct {
		SID      *int64  `json:"sid"`
		Password *string `json:"password"`
	}{}

	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.Password == nil || result.SID == nil {
		return errors.New("invalid input to connect to task")
	} else {
		c.SID = *result.SID
		c.Password = *result.Password
	}
	return nil
}

// ConnectToTask godoc
// @summary ConnectToTask
// @description connects new user to current task
// @security ApiKeyAuth
// @id connect-to-task
// @tags task
// @accept json
// @produce json
// @param input body ConnectInput true "info for connection"
// @success 200 {object} response
// @success 400 {object} response
// @success 403 {object} response
// @failure 500 {object} response
// @router /solution/connect [post]
func (h *Handler) ConnectToTask(c *fiber.Ctx) error {
	_, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	var request ConnectInput
	if err := c.BodyParser(&request); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, err)
	}

	svc := h.di.GetInstanceService()
	if err := svc.Task.ConnectToTask(c.UserContext(), request.SID, request.Password); err != nil {
		return sendErrorResponse(c, fiber.StatusForbidden, err)
	}

	if alts, err := svc.Task.GetAlts(c.UserContext(), request.SID); err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	} else if len(alts) == 0 {
		return sendErrorResponse(c, fiber.StatusBadRequest, errors.New("maintainer doesn't specified alternatives yet"))
	}

	if criteria, err := svc.Task.GetCriteria(c.UserContext(), request.SID); err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	} else if len(criteria) == 0 {
		return sendErrorResponse(c, fiber.StatusBadRequest, errors.New("maintainer doesn't specified criteria yet"))
	}

	return c.JSON(response{Message: "success"})
}

// SetExpertsWeights godoc
// @summary SetExpertsWeights
// @description sets experts weights for current group task
// @security ApiKeyAuth
// @id set-experts-weight
// @tags task
// @accept json
// @produce json
// @param sid query int true "task identifier"
// @param input body entity.Weights false "weights of experts"
// @success 200 {object} response
// @success 400 {object} response
// @success 403 {object} response
// @success 404 {object} response
// @failure 500 {object} response
// @router /solution/experts [put]
func (h *Handler) SetExpertsWeights(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	sid, err := strconv.ParseInt(c.Query("sid"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusNotFound, errors.New("task doesnt specified"))
	}

	request := entity.Weights{}

	if err := c.BodyParser(&request); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, err)
	}

	svc := h.di.GetInstanceService()
	if err := svc.Task.ValidateUser(c.UserContext(), uid, sid); err != nil {
		return sendErrorResponse(c, fiber.StatusForbidden, err)
	}

	if err := svc.Task.SetExpertsWeights(c.UserContext(), sid, request); err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(response{Message: "success"})
}

// GetExperts godoc
// @summary GetExperts
// @description gets experts of current group task
// @security ApiKeyAuth
// @id get-experts
// @tags task
// @accept json
// @produce json
// @param sid query int true "task identifier"
// @success 200 {object} []entity.Expert
// @success 400 {object} response
// @success 403 {object} response
// @success 404 {object} response
// @failure 500 {object} response
// @router /solution/experts [get]
func (h *Handler) GetExperts(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	sid, err := strconv.ParseInt(c.Query("sid"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusNotFound, errors.New("task doesnt specified"))
	}

	svc := h.di.GetInstanceService()
	if err := svc.Task.CheckAccess(c.UserContext(), uid, sid); err != nil {
		return sendErrorResponse(c, fiber.StatusForbidden, err)
	}

	task, err := svc.Task.GetTask(c.UserContext(), sid)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	} else if task.TaskType == v.Individuals {
		return sendErrorResponse(c, fiber.StatusBadRequest, errors.New("individual task doesnt provide this function"))
	}

	users, err := svc.User.GetUsersRelateToTask(c.UserContext(), sid)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(users)
}

// DeactivateStatuses godoc
// @summary DeactivateStatuses
// @description deactivate statuses of all experts relating to current task if maintainer choose this option
// @security ApiKeyAuth
// @id deactivate-statuses
// @tags task
// @accept json
// @produce json
// @param sid query int true "task identifier"
// @success 200 {object} response
// @success 403 {object} response
// @success 404 {object} response
// @failure 500 {object} response
// @router /solution/settings/experts [patch]
func (h *Handler) DeactivateStatuses(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	sid, err := strconv.ParseInt(c.Query("sid"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusNotFound, errors.New("task doesn't specified"))
	}

	service := h.di.GetInstanceService()
	if err := service.Task.ValidateUser(c.UserContext(), uid, sid); err != nil {
		return sendErrorResponse(c, fiber.StatusForbidden, err)
	}

	if err := service.Matrix.DeactivateStatuses(c.UserContext(), sid); err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(response{Message: "success"})
}

// GetDefaultLingScale godoc
// @summary GetDefaultLingScale
// @description get defaults linguistic scales
// @security ApiKeyAuth
// @id get-default-scales
// @tags task
// @accept json
// @produce json
// @success 200 {object} fiber.Map
// @router /solution/defaults [get]
func (h *Handler) GetDefaultLingScale(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"number": eval.DefaultNumberScale, "interval": eval.DefaultIntervalScale,
		"t1fs": eval.DefaultT1FSScale})
}
