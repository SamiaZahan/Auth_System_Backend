package handler

import (
	"errors"

	"github.com/emamulandalib/airbringr-auth/dto"
	"github.com/emamulandalib/airbringr-auth/response"
	"github.com/emamulandalib/airbringr-auth/service"
	"github.com/gofiber/fiber/v2"
)

func (receiver *Handler) MobileVerificationOTP(c *fiber.Ctx) (err error) {
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
	if err = svc.MobileVerificationOtp(*input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Payload{
			Message: err.Error(),
			Errors:  err,
			Data:    dto.VerificationOutput{Verified: false},
		})
	}

	return c.JSON(response.Payload{
		Message: "Please check your SMS.",
		Data:    dto.VerificationOutput{Verified: true},
	})
}

func (receiver *Handler) VerifyMobile(c *fiber.Ctx) (err error) {
	input := new(dto.VerifyMobileInput)

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
	if err = svc.VerifyAndRegisterMobileNumber(*input); err != nil {
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
