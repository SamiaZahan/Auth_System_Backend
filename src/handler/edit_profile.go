package handler

import (
	"errors"
	"github.com/emamulandalib/airbringr-auth/dto"
	"github.com/emamulandalib/airbringr-auth/response"
	"github.com/emamulandalib/airbringr-auth/service"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) EditProfile(c *fiber.Ctx) (err error) {
	input := new(dto.EditProfileInput)
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
	token := c.Get("Authorization", "")
	ExtractJwtClaims := service.ExtractJwtClaims{}
	claims, ok := ExtractJwtClaims.ExtractClaims(token)
	if !ok {
		return c.JSON(response.Payload{
			Message: "Something  Wrong",
		})
	}
	email := claims["email"].(string)
	svc := service.EditProfile{}
	if err := svc.EditUserProfile(input, email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Payload{
			Message: err.Error(),
			Errors:  err,
		})
	}
	return c.JSON(response.Payload{
		Message: "Profile Updated Successfully",
	})

}
