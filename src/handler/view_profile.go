package handler

import (
	"github.com/emamulandalib/airbringr-auth/response"
	"github.com/emamulandalib/airbringr-auth/service"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func (h *Handler) ViewProfile(c *fiber.Ctx) error {

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	svc := service.ViewProfile{}
	err, res, address := svc.ViewUserProfile(email, c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Payload{
			Message: err.Error(),
			Errors:  err,
		})
	}
	return c.JSON(response.Payload{
		Message: "Data Found",
		Data:    fiber.Map{"user": res, "address": address},
	})
}
