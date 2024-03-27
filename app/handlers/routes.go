package handlers

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"webApp/app/configs"
	"webApp/app/usecase"
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
	SetMiddleware(api, h.cfg)
	h.userIdentity = userIdentity

	authGroup := api.Group("auth")
	{
		authGroup.Post("/", h.CreateNewUser)
		authGroup.Get("/", h.LogIn)
		authGroup.Put("/refresh", h.RefreshToken)
	}

	userGroup := api.Group("/user")
	{
		userGroup.Get("/", h.GetAllSolutions)
		userGroup.Delete("/", h.DeleteSolution)
		userGroup.Patch("/settings", h.UpdateUser)
		userGroup.Delete("/settings", h.DeleteUser)
	}

	solGroup := api.Group("/solution")
	{
		solSettings := solGroup.Group("/settings")
		{
			solSettings.Post("/", h.CreateSolution)
			solSettings.Put("/", h.UpdateSolution)
			solSettings.Get("/", h.GetSolution)
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

		solGroup.Get("/connect", h.ConnectToSolution)

		solRating := solGroup.Group("/rating")
		{
			solRating.Post("/", h.CreateMatrix)
			solRating.Put("/", h.UpdateMatrix)
			solRating.Get("/", h.GetRatings)
		}

		solExperts := solGroup.Group("/experts")
		{
			solExperts.Get("/experts", h.GetExperts)
			solExperts.Put("/experts", h.SetExpertsWeights)
		}

		solGroup.Get("/final", h.GetFinal)
	}
}

type response struct {
	Message string `json:"message"`
}

func sendErrorResponse(c *fiber.Ctx, status int, err error) error {
	logrus.Error(err)
	if err := c.SendStatus(status); err != nil {
		return err
	}

	data, err := json.Marshal(response{Message: err.Error()})
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Send(data)
}
