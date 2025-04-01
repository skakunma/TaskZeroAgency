package main

import (
	"github.com/gofiber/fiber/v2"
	config "github.com/skakunma/TaskZeroAgency/internal/config"
	"github.com/skakunma/TaskZeroAgency/internal/handlers"
)

func main() {
	app := fiber.New()

	config := config.CreateConfig()
	handlers.LoadHandlers(app, config)

	app.Listen(":8080")
}
