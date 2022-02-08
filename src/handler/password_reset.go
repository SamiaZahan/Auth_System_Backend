package handler

import (
	"github.com/emamulandalib/airbringr-auth/dto"
	"github.com/emamulandalib/airbringr-auth/response"
	"github.com/emamulandalib/airbringr-auth/service"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

func (receiver *Handler) PasswordReset(c *fiber.Ctx) (err error) {
	input := new(dto.PasswordReset)

	if err = c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Payload{
			Message: response.BodyParseFailedErrorMsg,
			Errors:  errors.New(response.BodyParseFailedErrorMsg),
		})
	}

	if err = input.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Payload{
			Message: response.ValidationFailedMsg,
			Errors:  err,
		})
	}

	svc := service.PassReset{}
	if err = svc.UpdatePassword(*input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Payload{
			Message: err.Error(),
			Errors:  err,
			Data:    dto.VerificationOutput{Verified: false},
		})
	}

	return c.JSON(response.Payload{
		Message: "Password Reset Successfully",
		Data:    dto.VerificationOutput{Verified: true},
	})
}
