package handler

import (
	"github.com/emamulandalib/airbringr-auth/dto"
	"github.com/emamulandalib/airbringr-auth/response"
	"github.com/emamulandalib/airbringr-auth/service"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) EmailVerification(c *fiber.Ctx) (err error) {
	input := new(dto.EmailVerificationInput)
	if err = c.BodyParser(input); err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(response.Payload{
				Message: response.BodyParseFailedErrorMsg,
				Errors:  err,
			})
	}

	if err = input.Validate(); err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(response.Payload{
				Message: response.ValidationFailedMsg,
				Errors:  err,
			})
	}

	auth := service.Auth{}
	if err = auth.EmailVerification(*input); err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(response.Payload{
				Message: err.Error(),
				Errors:  err,
				Data:    dto.EmailVerificationOutput{Verified: false},
			})
	}

	return c.JSON(response.Payload{
		Message: "Your account activated successfully. Please try to login now.",
		Data:    dto.EmailVerificationOutput{Verified: true},
	})
}
