package handler

import (
	"errors"
	"github.com/emamulandalib/airbringr-auth/dto"
	"github.com/emamulandalib/airbringr-auth/response"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Login(c *fiber.Ctx) error {
	loginInput := new(dto.LoginInput)

	if err := c.BodyParser(loginInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Payload{
			Message: response.BodyParseFailedErrorMsg,
			Errors:  errors.New(response.BodyParseFailedErrorMsg),
		})
	}

	err := loginInput.Validate()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Payload{
			Message: response.ValidationFailedMsg,
			Errors:  err,
		})
	}

	return c.JSON(response.Payload{
		Message: "Please check your email for the OTP.",
	})
}
