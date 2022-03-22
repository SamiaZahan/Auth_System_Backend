package handler

import (
	"github.com/emamulandalib/airbringr-auth/dto"
	"github.com/emamulandalib/airbringr-auth/response"
	"github.com/emamulandalib/airbringr-auth/service"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func (h *Handler) VerifyAndEditEmail(c *fiber.Ctx) (err error) {
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

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	auth := service.Auth{}
	if err = auth.VerifyAndEditEmail(*input, email); err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(response.Payload{
				Message: err.Error(),
				Errors:  err,
				Data:    dto.VerificationOutput{Verified: false},
			})
	}

	return c.JSON(response.Payload{
		Message: "new email updated successfully",
		Data:    dto.VerificationOutput{Verified: true},
	})
}
