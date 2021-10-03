package handler

import (
	"github.com/emamulandalib/airbringr-auth/dto"
	"github.com/emamulandalib/airbringr-auth/response"
	"github.com/emamulandalib/airbringr-auth/service"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) SendEmailOTP(c *fiber.Ctx) (err error) {
	input := new(dto.EmailOtpInput)
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

	svc := service.EmailOtp{}
	if err = svc.Send(*input); err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(response.Payload{
				Message: err.Error(),
				Errors:  err,
			})
	}

	return c.JSON(response.Payload{
		Message: "Please check your email for the OTP.",
		Data:    dto.VerificationOutput{Verified: true},
	})
}
