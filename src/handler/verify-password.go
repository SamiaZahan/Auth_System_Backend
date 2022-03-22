package handler

import (
	"errors"
	"github.com/emamulandalib/airbringr-auth/dto"
	"github.com/emamulandalib/airbringr-auth/response"
	"github.com/emamulandalib/airbringr-auth/service"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func (h *Handler) VerifyPassword(c *fiber.Ctx) (err error) {
	input := new(dto.VerifyPassword)
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

	//token := c.Get("Authorization", "")
	//ExtractJwtClaims := service.ExtractJwtClaims{}
	//claims, ok := ExtractJwtClaims.ExtractClaims(token)
	//if !ok {
	//	return c.JSON(response.Payload{
	//		Message: "Something  Wrong",
	//	})
	//}
	//email := claims["email"].(string)
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	svc := service.VerifyUserPassword{}
	if err = svc.VerifyPassword(input.Password, email); err != nil {
		return c.JSON(response.Payload{
			Message: err.Error(),
			Errors:  err,
			Data:    dto.VerificationOutput{Verified: false},
		})
	}

	return c.JSON(response.Payload{
		Message: "Password matched",
		Data:    dto.VerificationOutput{Verified: true},
	})
}
