package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx"
	"github.com/skakunma/TaskZeroAgency/internal/config"
	"github.com/skakunma/TaskZeroAgency/internal/jwtAuth"
	"math/rand"
	"net/http"
	"time"
)

func RandomUnixTime() int {
	now := time.Now()

	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	randomSeconds := rand.Intn(24 * 60 * 60) // В сутках 24 * 60 * 60 секунд
	randomTime := int(startOfDay.Unix()) + randomSeconds

	return randomTime
}

func AuthMiddleWare(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := RandomUnixTime()
		token, err := jwtAuth.BuildJWTString(cfg, user)
		if err != nil {
			cfg.Logger.Log(pgx.LogLevelInfo, "Ошибка генерации JWT:", nil)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "In service"})
		}

		jwtToken := token

		c.Response().Header.Set("Authorization", "Bearer "+jwtToken)
		return c.Next()

	}
}
