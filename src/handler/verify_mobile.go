package handler

import (
	"errors"
	"github.com/emamulandalib/airbringr-auth/dto"
	"github.com/emamulandalib/airbringr-auth/response"
	"github.com/emamulandalib/airbringr-auth/service"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func (h *Handler) MobileVerificationOTP(c *fiber.Ctx) (err error) {
	input := new(dto.SendSmsOtpInput)
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

	svc := service.SmsOtp{}
	if err = svc.MobileVerificationOtp(*input, c); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Payload{
			Message: err.Error(),
			Errors:  err,
			Data:    dto.VerificationOutput{Verified: false},
		})
	}

	return c.JSON(response.Payload{
		Message: "Please check your SMS. The OTP will be valid for 5 minutes.",
	})
}

func (h *Handler) VerifyMobile(c *fiber.Ctx) (err error) {
	input := new(dto.VerifyMobileInput)

	if err = c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Payload{
			Message: response.BodyParseFailedErrorMsg,
			Errors:  errors.New(response.BodyParseFailedErrorMsg),
		})
	}

	if input.State == "edit" {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		email := claims["email"].(string)
		input.Auth = email
	}

	if err = input.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Payload{
			Message: response.ValidationFailedMsg,
			Errors:  err,
		})
	}

	svc := service.SmsOtp{}
	if err = svc.VerifyAndRegisterMobileNumber(*input, c); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Payload{
			Message: err.Error(),
			Errors:  err,
			Data:    dto.VerificationOutput{Verified: false},
		})
	}

	return c.JSON(response.Payload{
		Message: "Mobile verified.",
		Data:    dto.VerificationOutput{Verified: true},
	})
}
