package handler

import (
	"errors"

	"github.com/emamulandalib/airbringr-auth/response"
	"github.com/gofiber/fiber/v2"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Home(c *fiber.Ctx) error {
	return c.JSON(response.Payload{
		Message: "V1!",
	})
}

func (h *Handler) NotFound(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(response.Payload{
		Message: "404!",
		Errors:  errors.New("path not found"),
	})
}
