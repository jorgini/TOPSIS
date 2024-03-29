package handlers

import (
	"errors"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"webApp/configs"
)

var claimsKey string

func SetMiddleware(api fiber.Router, config *configs.AppConfig) {
	middleware := jwtware.New(jwtware.Config{
		Filter:     SkipMiddleware,
		SigningKey: jwtware.SigningKey{Key: []byte(config.SignKey)},
		ContextKey: "id",
	})

	claimsKey = config.ClaimsKey
	api.Use(middleware)

}

func SkipMiddleware(c *fiber.Ctx) bool {
	url := strings.Clone(c.OriginalURL())
	if url == "/auth" || url == "/auth/refresh" {
		return true
	}
	return false
}

func userIdentity(c *fiber.Ctx) (int64, error) {
	if c.Locals("id") == nil {
		return 0, errors.New("cant read user id")
	}
	token := c.Locals("id").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	if claims[claimsKey] == nil {
		return 0, errors.New("cant read user id")
	}
	intiId := claims[claimsKey].(float64)
	id := int64(intiId)
	//logrus.Infof("successful auth with id %d", id)
	return id, nil
}
