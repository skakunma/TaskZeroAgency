package handlers

import (
	"github.com/gofiber/fiber/v2"
	config "github.com/skakunma/TaskZeroAgency/internal/config"
	"github.com/skakunma/TaskZeroAgency/internal/middleware"
)

func LoadHandlers(c *fiber.App, cfg *config.Config) {
	c.Use(middleware.AuthMiddleWare(cfg))
	c.Post("/edit/:id", UpdateNew(cfg))
	c.Post("/edit/", CreateNew(cfg))
	c.Get("/list", GetNews(cfg))
}
