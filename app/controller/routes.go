package controller

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/sirupsen/logrus"
	"webApp/configs"
	"webApp/usecase"
)

type Handler struct {
	di           usecase.DiService
	cfg          *configs.AppConfig
	userIdentity func(c *fiber.Ctx) (int64, error)
}

func NewHandler(di usecase.DiService, config *configs.AppConfig) *Handler {
	return &Handler{
		di:  di,
		cfg: config,
	}
}

func (h *Handler) SetAllRoutes(api fiber.Router) {
	// Swagger endpoint
	api.Get("/swagger/*", swagger.HandlerDefault)

	SetMiddleware(api, h.cfg)
	h.userIdentity = userIdentity

	authGroup := api.Group("auth")
	{
		authGroup.Post("/sign-up", h.CreateNewUser)
		authGroup.Post("/log-in", h.LogIn)
		authGroup.Put("/refresh", h.RefreshToken)
	}

	userGroup := api.Group("/user")
	{
		userGroup.Post("/", h.CreateTask)
		userGroup.Get("/", h.GetAllSolutions)
		userGroup.Delete("/tasks", h.DeleteTask)
		userGroup.Patch("/settings", h.UpdateUser)
		userGroup.Delete("/settings", h.DeleteUser)
	}

	solGroup := api.Group("/solution")
	{
		solSettings := solGroup.Group("/settings")
		{
			solSettings.Put("/", h.UpdateTask)
			solSettings.Patch("/experts", h.DeactivateStatuses)
			solSettings.Get("/", h.GetTask)
			solSettings.Patch("/pass", h.SetPassword)
		}

		solAlts := solGroup.Group("/alternatives")
		{
			solAlts.Post("/", h.SetAlternatives)
			solAlts.Put("/", h.ReplaceAlternatives)
			solAlts.Get("/", h.GetAlternatives)
		}

		solCriteria := solGroup.Group("/criteria")
		{
			solCriteria.Post("/", h.SetCriteria)
			solCriteria.Put("/", h.ReplaceCriteria)
			solCriteria.Get("/", h.GetCriteria)
		}

		solGroup.Post("/connect", h.ConnectToTask)

		solRating := solGroup.Group("/rating")
		{
			solRating.Post("/", h.CreateMatrix)
			solRating.Put("/", h.UpdateMatrix)
			solRating.Get("/", h.GetRatings)
			solRating.Patch("/complete", h.CompleteStatus)
		}

		solExperts := solGroup.Group("/experts")
		{
			solExperts.Get("/", h.GetExperts)
			solExperts.Put("/", h.SetExpertsWeights)
			solExperts.Get("/role", h.GetRole)
		}

		solGroup.Post("/final", h.GetFinal)
	}
}

type response struct {
	Message string `json:"message"`
}

func sendErrorResponse(c *fiber.Ctx, status int, err error) error {
	logrus.Error(err)
	data, err := json.Marshal(response{Message: err.Error()})
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return fiber.NewError(status, string(data))
}
