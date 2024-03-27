package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"strconv"
	"webApp/app/entity"
	"webApp/lib/eval"
	v "webApp/lib/variables"
)

type TaskInput struct {
	Title        string               `json:"title"`
	Description  string               `json:"description"`
	TaskType     string               `json:"task_type"`
	Method       string               `json:"method"`
	CalcSettings int64                `json:"calc_settings"`
	LingScale    eval.LinguisticScale `json:"ling_scale"`
}

func (t *TaskInput) UnmarshalJSON(data []byte) error {
	result := struct {
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

func (h *Handler) CreateSolution(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
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
	sid, err := svc.Task.CreateNewTask(c.UserContext(), &task)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	logrus.Infof("user with id %d successfully create new input with sid %d", uid, sid)
	return c.JSON(fiber.Map{"status": "success"})
}

func (h *Handler) UpdateSolution(c *fiber.Ctx) error {
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
	return c.JSON(fiber.Map{"status": "success"})
}

func (h *Handler) GetSolution(c *fiber.Ctx) error {
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

	task, err := svc.Task.GetTask(c.UserContext(), sid)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}
	output := TaskInput{
		Title:        task.Title,
		Description:  task.Description,
		TaskType:     task.TaskType,
		Method:       task.Method,
		CalcSettings: task.CalcSettings,
		LingScale:    task.LingScale,
	}

	return c.JSON(output)
}

func (h *Handler) DeleteSolution(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	request := struct {
		Sid int64 `json:"sid"`
	}{}
	if err := c.BodyParser(&request); err != nil || request.Sid == 0 {
		return sendErrorResponse(c, fiber.StatusBadRequest, errors.New("cant find task identification"))
	}

	svc := h.di.GetInstanceService()
	if err := svc.Task.CheckAccess(c.UserContext(), uid, request.Sid); err != nil {
		return sendErrorResponse(c, fiber.StatusForbidden, errors.New("hasn't access to solution"))
	}

	if err := svc.Task.DeleteTask(c.UserContext(), uid, request.Sid); err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	logrus.Infof("user with id %d successfully delete task with sid %d", uid, request.Sid)
	return c.JSON(fiber.Map{"status": "success"})
}

func (h *Handler) SetPassword(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	sid, err := strconv.ParseInt(c.Query("sid"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusNotFound, errors.New("task doesn't specified"))
	}

	request := struct {
		Pass string `json:"password"`
	}{}
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
	return c.JSON(fiber.Map{"status": "success"})
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

func (h *Handler) ConnectToSolution(c *fiber.Ctx) error {
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

	return c.JSON(fiber.Map{"status": "success"})
}

func (h *Handler) SetExpertsWeights(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	sid, err := strconv.ParseInt(c.Query("sid"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, errors.New("task doesnt specified"))
	}

	request := struct {
		Weights entity.Weights `json:"weights"`
	}{}

	if err := c.BodyParser(&request); err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, err)
	}

	svc := h.di.GetInstanceService()
	if err := svc.Task.ValidateUser(c.UserContext(), uid, sid); err != nil {
		return sendErrorResponse(c, fiber.StatusForbidden, err)
	}

	if err := svc.Task.SetExpertsWeights(c.UserContext(), sid, request.Weights); err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(fiber.Map{"status": "success"})
}

func (h *Handler) GetExperts(c *fiber.Ctx) error {
	uid, err := h.userIdentity(c)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	sid, err := strconv.ParseInt(c.Query("sid"), 10, 64)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, errors.New("task doesnt specified"))
	}

	svc := h.di.GetInstanceService()
	if err := svc.Task.ValidateUser(c.UserContext(), uid, sid); err != nil {
		return sendErrorResponse(c, fiber.StatusForbidden, err)
	}

	task, err := svc.Task.GetTask(c.UserContext(), sid)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	} else if task.TaskType == v.Individuals {
		return sendErrorResponse(c, fiber.StatusBadRequest, errors.New("individuals task doesnt provide this function"))
	}

	users, err := svc.User.GetUsersRelateToTask(c.UserContext(), sid)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(users)
}
