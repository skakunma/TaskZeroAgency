package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx"
	"github.com/skakunma/TaskZeroAgency/internal/config"
	"github.com/skakunma/TaskZeroAgency/internal/storage"
	"net/http"
	"strconv"
)

func GetNews(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		news, err := cfg.Store.GetNews(ctx)
		if err != nil {
			cfg.Logger.Log(pgx.LogLevelInfo, err.Error(), nil)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"Success": false, "error": err.Error()})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{"Success": true, "News": news})
	}
}

func CreateNew(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var newNews storage.New

		if err := c.BodyParser(&newNews); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		ctx := c.Context()
		err := cfg.Store.CreateNew(ctx, newNews)
		if err != nil {
			cfg.Logger.Log(pgx.LogLevelInfo, err.Error(), nil)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"Success": false, "error": err.Error()})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"Id":         newNews.Id,
			"Title":      newNews.Title,
			"Content":    newNews.Content,
			"Categories": newNews.Categories,
		})
	}
}

func UpdateNew(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		intID, err := strconv.Atoi(id) // Пробуем преобразовать строку в целое число
		var newInfo storage.New

		if id == "" || err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"Success": false, "error": "Id is not found or id can't be int"})
		}

		ctx := c.Context()
		infoOld, err := cfg.Store.GetNewFromID(ctx, intID)
		if errors.Is(err, storage.ErrNotFound) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"Success": false, "error": "Can't find New with this id"})
		}

		if err = c.BodyParser(&newInfo); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		err = cfg.Store.UpdateNewFromID(ctx, infoOld.Id, newInfo)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"Success": false, "error": "error in storage"})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"Id":         newInfo.Id,
			"Title":      newInfo.Title,
			"Content":    newInfo.Content,
			"Categories": newInfo.Categories,
		})
	}
}
